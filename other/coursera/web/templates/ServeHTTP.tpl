func (srv *{{.ApiName}}) ServeHTTP(w http.ResponseWriter, r *http.Request) {
{{ range .HandleFunctions }}
    if r.URL.Path == "{{.Url}}" {
        {{ if .Method }}
            if r.Method != "{{.Method}}" {
                resp, _ := jsonResponse("bad method", nil)
                http.Error(w, string(resp), http.StatusNotAcceptable)
                return
            }
        {{ end }}

        {{ if .Auth }}
        if r.Header.Get("X-Auth") != "100500" {
            resp, _ := jsonResponse("unauthorized", nil)
            http.Error(w, string(resp), http.StatusForbidden)
            return
        }
        {{ end }}

        srv.{{.Name}}Handler(w, r)
        return
    }
{{ end }}
    resp, _ := jsonResponse("unknown method", nil)
    http.Error(w, string(resp), http.StatusNotFound)
}

{{ range .HandleFunctions }}
func (srv *{{$.ApiName}}) {{.Name}}Handler(w http.ResponseWriter, r *http.Request) {
    in := {{.InParamsStructName}}{}
    var query url.Values
    if r.Method == "GET" {
        query = r.URL.Query()
        _ = decoder.Decode(&in, r.URL.Query())
    } else {
        data, err := ioutil.ReadAll(r.Body)
        query, err = url.ParseQuery(string(data))
        if err != nil {
            http.Error(w, "Can't parse body", http.StatusBadRequest)
            return
        }

    }
    _ = decoder.Decode(&in, query)

    {{ range $field, $rules := .ValidationRules }}
        {{ range $rule :=  $rules }}

            {{with $queryKey := $field }}
                {{ if eq $rule.Name "paramname" }}
                    {{ $queryKey = $rule.Value }}
                    in.{{ $field }} = query.Get("{{ $queryKey }}")
                {{ end }}

                {{ if eq $rule.Name "assert_int" }}
                    if _, err := strconv.Atoi(query.Get("{{ $queryKey | ToLower }}")); err != nil {
                        resp, _ := jsonResponse("{{ $field | ToLower }} must be int", nil)
                        http.Error(w, string(resp), http.StatusBadRequest)
                        return
                    }
                {{ end }}
            {{ end }}

            {{ if eq $rule.Name "required" }}
                if in.{{ $field }} == "" {
                    resp, _ := jsonResponse("{{ $field | ToLower }} must me not empty", nil)
                    http.Error(w, string(resp), http.StatusBadRequest)
                    return
                }
            {{ end }}

            {{ if eq $rule.Name "min" }}
                if {{if eq $rule.DataType "string"}}len(in.{{ $field }}){{ else }}in.{{ $field }}{{ end }} < {{ $rule.Value }} {
                    resp, _ := jsonResponse("{{ $field | ToLower }} {{if eq $rule.DataType "string"}}len {{end}}must be >= {{ $rule.Value }}", nil)
                    http.Error(w, string(resp), http.StatusBadRequest)
                    return
                }
            {{ end }}

            {{ if eq $rule.Name "max" }}
                if {{if eq $rule.DataType "string"}}len(in.{{ $field }}){{ else }}in.{{ $field }}{{ end }} > {{ $rule.Value }} {
                    resp, _ := jsonResponse("{{ $field | ToLower }} {{if eq $rule.DataType "string"}}len {{end}}must be <= {{ $rule.Value }}", nil)
                    http.Error(w, string(resp), http.StatusBadRequest)
                    return
                }
            {{ end }}

            {{ if eq $rule.Name "enum" }}
                if funk.IndexOf(strings.Split("{{ $rule.Value }}", "|"), in.{{ $field }}) == -1 && in.{{ $field }} != "" {
                    resp, _ := jsonResponse(fmt.Sprintf("{{ $field | ToLower }} must be one of [%s]", strings.ReplaceAll("{{ $rule.Value }}", "|", ", ")), nil)
                    http.Error(w, string(resp), http.StatusBadRequest)
                    return
                }
            {{ end }}

            {{ if eq $rule.Name "default" }}
                if in.{{$field}} == "" {
                    in.{{$field}} = "{{ $rule.Value }}"
                }
            {{ end }}

        {{ end }}
    {{ end }}

    user, err := srv.{{.Name}}(r.Context(), in)
    if err != nil {
        if apiErr, ok := err.(ApiError); ok {
            resp, _ := jsonResponse(apiErr.Error(), nil)
            http.Error(w, string(resp), apiErr.HTTPStatus)
        } else {
            resp, _ := jsonResponse(err.Error(), nil)
            http.Error(w, string(resp), http.StatusInternalServerError)
        }
        return
    }
    jsonBytes, err := json.Marshal(&Response{"error": "", "response": user})
    _, err = w.Write(jsonBytes)
    if err != nil {
        http.Error(w, "Can't encode json", http.StatusInternalServerError)
    }
}
{{ end }}



