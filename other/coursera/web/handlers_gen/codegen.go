package main

import (
	"bytes"
	"encoding/json"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"strings"
	"text/template"
)

var (
	tpl = template.Must(template.New("").Funcs(template.FuncMap{
		"ToLower": strings.ToLower,
	}).ParseFiles(
		"./templates/Main.tpl",
		"./templates/ServeHTTP.tpl",
	))
)

type Rule struct {
	DataType, Name, Value string
}

type apiGenFunction struct {
	Name, InParamsStructName string
	Url                      string `json:"url"`
	Method                   string `json:"method"`
	Auth                     bool   `json:"auth"`
	ValidationRules          map[string][]Rule
}

func getValidationRulesByStructName(structName string, file *ast.File) map[string][]Rule {
	rules := make(map[string][]Rule)
	for _, decl := range file.Decls {
		decl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		nodeSpec, ok := decl.Specs[0].(*ast.TypeSpec)
		if !ok {
			continue
		}
		nodeName := nodeSpec.Name.Name
		_, isStruct := nodeSpec.Type.(*ast.StructType)
		if nodeName == structName && isStruct {
			for _, field := range nodeSpec.Type.(*ast.StructType).Fields.List {
				name := field.Names[0].Name
				tag := strings.Trim(field.Tag.Value, "`")
				if strings.HasPrefix(tag, "apivalidator:") {
					rulesStr := strings.Replace(tag, "apivalidator:", "", -1)
					rulesStr = strings.Trim(rulesStr, `"`)
					rulesTmp := strings.Split(rulesStr, ",")
					for _, r := range rulesTmp {
						ro := Rule{
							DataType: field.Type.(*ast.Ident).Name,
							Value:    "",
						}
						splitR := strings.Split(r, "=")
						ro.Name = splitR[0]
						if len(splitR) == 2 {
							ro.Value = splitR[1]
						}
						rules[name] = append(rules[name], ro)
					}
				}
				if field.Type.(*ast.Ident).Name == "int" {
					rules[name] = append(rules[name], Rule{Name: "assert_int"})
				}
			}
		}
	}
	return rules
}

func getApiGenFunctions(file *ast.File) map[string][]apiGenFunction {
	apisFunctions := make(map[string][]apiGenFunction)
	for _, decl := range file.Decls {
		function, ok := decl.(*ast.FuncDecl)
		if !ok || function.Doc == nil {
			continue
		}
		for _, comment := range function.Doc.List {
			if strings.HasPrefix(comment.Text, "// apigen:api") {
				apiObj, ok := function.Recv.List[0].Type.(*ast.StarExpr)
				if !ok {
					continue
				}
				apiName := apiObj.X.(*ast.Ident).Name

				paramsStructName := function.Type.Params.List[1].Type.(*ast.Ident).Name
				f := apiGenFunction{
					Name:               function.Name.Name,
					InParamsStructName: paramsStructName,
					ValidationRules:    getValidationRulesByStructName(paramsStructName, file),
				}
				err := json.Unmarshal([]byte(fetchJson(comment.Text)), &f)
				if err != nil {
					log.Fatal(err)
				}
				apisFunctions[apiName] = append(apisFunctions[apiName], f)
			}
		}
	}
	return apisFunctions
}

func main() {
	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, os.Args[1], nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	out, err := os.Create(os.Args[2])
	_ = out
	if err != nil {
		log.Fatal(err)
	}

	apisFunctions := getApiGenFunctions(file)
	serveHTTPBuffer := &bytes.Buffer{}
	for apiName, functions := range apisFunctions {
		err = tpl.ExecuteTemplate(serveHTTPBuffer, "ServeHTTP.tpl", struct {
			ApiName         string
			HandleFunctions []apiGenFunction
		}{apiName, functions})
		if err != nil {
			log.Fatal(err)
		}
	}

	_ = tpl.ExecuteTemplate(out, "Main.tpl", struct {
		ApisHTTPServe string
	}{serveHTTPBuffer.String()})

}

func fetchJson(str string) string {
	startIndex := strings.Index(str, "{")
	finishIndex := strings.Index(str, "}")
	return str[startIndex : finishIndex+1]
}
