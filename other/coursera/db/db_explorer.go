package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

const (
	limitParam  = "limit"
	offsetParam = "offset"

	tablesListPath = "/"
)

func prepareCollectionResp(respBody map[string][]interface{}) []byte {
	resp := make(map[string]map[string][]interface{})
	resp["response"] = respBody
	jsonByte, _ := json.Marshal(resp)
	return jsonByte
}
func prepareResp(respBody map[string]interface{}) []byte {
	resp := make(map[string]map[string]interface{})
	resp["response"] = respBody
	jsonByte, _ := json.Marshal(resp)
	return jsonByte
}
func prepareErr(text string) []byte {
	resp := make(map[string]string)
	resp["error"] = text
	jsonByte, _ := json.Marshal(resp)
	return jsonByte
}
func queryError(err error) error {
	return fmt.Errorf("query error: %v", err)
}

type FieldDB struct {
	Field, Type, Collation, Null, Key, Default, Extra, Privileges, Comment *string
}

type TableDB struct {
	TablesInGo                string
	fields                    []FieldDB
	dbDataPointers            []interface{}
	fieldNamesDataPointersMap []string
}

func (t *TableDB) mergeDataPointersValuesWithFieldsName() map[string]interface{} {
	tmp := make(map[string]*string, len(t.dbDataPointers))
	for i, pointer := range t.dbDataPointers {
		sqlString := pointer.(*sql.NullString)
		if sqlString.Valid {
			tmp[t.fieldNamesDataPointersMap[i]] = &sqlString.String
		} else {
			tmp[t.fieldNamesDataPointersMap[i]] = nil
		}
	}

	result := make(map[string]interface{})
	for fieldName, fieldValue := range tmp {
		if *t.getFieldByName(fieldName).Type == "int" {
			v, _ := strconv.Atoi(*fieldValue)
			result[fieldName] = &v
		} else {
			result[fieldName] = fieldValue
		}
	}

	return result
}
func (t *TableDB) flushDataPointers() {
	t.dbDataPointers = make([]interface{}, len(t.fields))
	t.fieldNamesDataPointersMap = make([]string, len(t.fields))
	for i, field := range t.fields {
		var str sql.NullString
		t.dbDataPointers[i] = &str
		t.fieldNamesDataPointersMap[i] = *field.Field
	}
}
func (t *TableDB) getFieldByName(fieldName string) *FieldDB {
	for _, field := range t.fields {
		if *field.Field == fieldName {
			return &field
		}
	}
	return nil
} /*
func (t *TableDB) isFieldPK(fieldName string) bool {
	for _, field := range t.fields {
		if *field.Field == fieldName {
			return *field.Key == "PRI"
		}
	}
	return false
}*/
func (t *TableDB) validate(data map[string]interface{}) []string {
	var errors []string
	for fieldName, fieldValue := range data {
		f := t.getFieldByName(fieldName)
		if *f.Key == "PRI" {
			errors = append(errors, fmt.Sprintf("field %s have invalid type", *f.Field))
		}

		valueTypeToUpdate := fmt.Sprintf("%T", fieldValue)

		if valueTypeToUpdate == "<nil>" {
			if *f.Null != "YES" {
				errors = append(errors, fmt.Sprintf("field %s have invalid type", *f.Field))
			}
		} else {
			if strings.Contains(*f.Type, "varchar") || *f.Type == "text" {
				if valueTypeToUpdate != "string" {
					errors = append(errors, fmt.Sprintf("field %s have invalid type", *f.Field))
				}
			}

			if *f.Type == "int" {
				if valueTypeToUpdate != "float64" {
					errors = append(errors, fmt.Sprintf("field %s have invalid type", *f.Field))
				}
			}
		}
	}
	return errors
}
func (t *TableDB) clean(data map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for fieldName, fieldValue := range data {
		if t.getFieldByName(fieldName) != nil {
			result[fieldName] = fieldValue
		}
	}

	for _, field := range t.fields {
		if _, ok := result[*field.Field]; !ok {
			if *field.Key != "PRI" {
				if *field.Null == "YES" {
					result[*field.Field] = nil
				} else {
					result[*field.Field] = ""
				}
			}
		}
	}

	return result
}

func (t *TableDB) getPKField() *FieldDB {
	for _, field := range t.fields {
		if *field.Key == "PRI" {
			return &field
		}
	}
	return nil
}

type Explorer struct {
	*sql.DB
	tables []*TableDB
}

func NewDbExplorer(db *sql.DB) (http.Handler, error) {
	t, err := db.Query("SHOW TABLES")
	var tables []*TableDB
	if err != nil {
		return nil, queryError(err)
	}
	for t.Next() {
		var tableDB TableDB
		err = t.Scan(&tableDB.TablesInGo)
		if err != nil {
			return nil, queryError(err)
		}
		tables = append(tables, &tableDB)
	}
	_ = t.Close()

	for _, t := range tables {
		f, err := db.Query(fmt.Sprintf("SHOW FULL COLUMNS FROM `%s`", t.TablesInGo))
		if err != nil {
			return nil, queryError(err)
		}

		for f.Next() {
			var field FieldDB
			err = f.Scan(
				&field.Field, &field.Type, &field.Collation, &field.Null, &field.Key,
				&field.Default, &field.Extra, &field.Privileges, &field.Comment,
			)
			if err != nil {
				return nil, queryError(err)
			}
			t.fields = append(t.fields, field)
		}
		_ = f.Close()
	}

	return &Explorer{db, tables}, nil
}
func (e *Explorer) getTableDB(tableName string) *TableDB {
	for _, table := range e.tables {
		if table.TablesInGo == tableName {
			return table
		}
	}
	return nil
}

type TableRequestParams struct {
	tableName, id, limit, offset string
}

func (p *TableRequestParams) hasId() bool {
	return p.id != ""
}

func parseParams(r *http.Request) *TableRequestParams {
	var onlyTableName = regexp.MustCompile(`^/(.*?)($|/(.*?)(?:/|$))`)
	matched := onlyTableName.FindStringSubmatch(r.URL.Path)

	offset := r.URL.Query().Get(offsetParam)
	limit := r.URL.Query().Get(limitParam)

	if len(matched) == 4 {
		return &TableRequestParams{
			tableName: matched[1],
			id:        matched[3],
			limit:     limit,
			offset:    offset,
		}
	}

	if len(matched) == 2 {
		return &TableRequestParams{
			tableName: matched[1],
			limit:     limit,
			offset:    offset,
		}
	}

	return nil
}

func (e *Explorer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path == tablesListPath && r.Method == http.MethodGet {
		tablesResp := make(map[string][]interface{})
		for _, table := range e.tables {
			tablesResp["tables"] = append(tablesResp["tables"], table.TablesInGo)
		}
		_, _ = w.Write(prepareCollectionResp(tablesResp))
		return
	}

	params := parseParams(r)
	tableDB := e.getTableDB(params.tableName)
	if tableDB == nil {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write(prepareErr("unknown table"))
		return
	}

	if !params.hasId() && r.Method == http.MethodGet {
		wheres := make(map[string]string)
		sqlTpl := prepareSelectSqlTpl(r, wheres)
		rows, err := e.DB.Query(fmt.Sprintf(sqlTpl, params.tableName))
		if err != nil {
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}
		result := make(map[string][]interface{})
		for rows.Next() {
			tableDB.flushDataPointers()
			err := rows.Scan(tableDB.dbDataPointers...)
			if err != nil {
				http.Error(w, "Internal error", http.StatusInternalServerError)
				return
			}
			result["records"] = append(result["records"], tableDB.mergeDataPointersValuesWithFieldsName())
		}
		_, _ = w.Write(prepareCollectionResp(result))
		return
	}

	if params.hasId() && r.Method == http.MethodGet {
		wheres := map[string]string{*tableDB.getPKField().Field: params.id}
		sqlTpl := prepareSelectSqlTpl(r, wheres)
		_ = sqlTpl
		rows, err := e.DB.Query(fmt.Sprintf(sqlTpl, params.tableName), params.id)
		if err != nil {
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}
		result := make(map[string]interface{})
		for rows.Next() {
			tableDB.flushDataPointers()
			err := rows.Scan(tableDB.dbDataPointers...)
			if err != nil {
				http.Error(w, "Internal error", http.StatusInternalServerError)
				return
			}
			result["record"] = tableDB.mergeDataPointersValuesWithFieldsName()
		}

		if len(result) == 0 {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write(prepareErr("record not found"))
			return
		}

		_, _ = w.Write(prepareResp(result))
		return
	}

	if r.Method == http.MethodPut {
		dataToInsert := make(map[string]interface{})
		body, _ := ioutil.ReadAll(r.Body)
		_ = json.Unmarshal(body, &dataToInsert)
		_ = r.Body.Close()
		delete(dataToInsert, *tableDB.getPKField().Field)
		dataToInsert = tableDB.clean(dataToInsert)
		errors := tableDB.validate(dataToInsert)
		if len(errors) > 0 {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write(prepareErr(errors[0]))
			return
		}

		fields, values := fetchFieldsAndValues(dataToInsert)
		sqlStr := prepareInsertSqlTpl(params.tableName, fields)
		result, err := e.DB.Exec(sqlStr, values...)
		if err != nil {
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}
		lastId, _ := result.LastInsertId()
		resp := prepareResp(map[string]interface{}{
			*tableDB.getPKField().Field: lastId,
		})
		_, _ = w.Write(resp)
		return
	}

	if params.hasId() && r.Method == http.MethodPost {
		dataToUpdate := make(map[string]interface{})
		body, _ := ioutil.ReadAll(r.Body)
		_ = json.Unmarshal(body, &dataToUpdate)

		errors := tableDB.validate(dataToUpdate)
		if len(errors) > 0 {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write(prepareErr(errors[0]))
			return
		}

		fields, values := fetchFieldsAndValues(dataToUpdate)
		sqlStr := prepareUpdateSqlTpl(params.tableName, fields, map[string]string{*tableDB.getPKField().Field: params.id})
		bindings := values
		bindings = append(bindings, params.id)
		result, err := e.DB.Exec(sqlStr, bindings...)
		if err != nil {
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}
		updatedCount, _ := result.RowsAffected()
		resp := prepareResp(map[string]interface{}{
			"updated": updatedCount,
		})
		_, _ = w.Write(resp)
		return
	}

	if params.hasId() && r.Method == http.MethodDelete {
		sqlStr := prepareDeleteSqlTpl(params.tableName, map[string]string{*tableDB.getPKField().Field: params.id})
		result, err := e.DB.Exec(sqlStr, params.id)
		if err != nil {
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}
		deletedCount, _ := result.RowsAffected()
		resp := prepareResp(map[string]interface{}{
			"deleted": deletedCount,
		})
		_, _ = w.Write(resp)
		return
	}

}

func fetchFieldsAndValues(data map[string]interface{}) ([]string, []interface{}) {
	var fields []string
	var values []interface{}
	for fieldName, fieldValue := range data {
		fields = append(fields, fieldName)
		values = append(values, fieldValue)
	}
	return fields, values
}

func prepareSelectSqlTpl(r *http.Request, wheres map[string]string) string {
	sqlTpl := "SELECT * FROM `%s`"

	i := 0
	if len(wheres) > 0 {
		sqlTpl += " WHERE "
		for field := range wheres {
			if i > 0 {
				sqlTpl += " AND "
			}
			sqlTpl += fmt.Sprintf("`%s` = ?", field)
			i++
		}
	}

	limit := r.URL.Query().Get(limitParam)
	if limit != "" {
		_, err := strconv.Atoi(limit)
		if err == nil {
			sqlTpl += fmt.Sprintf(" LIMIT %s", limit)
		}
	}

	offset := r.URL.Query().Get(offsetParam)
	if offset != "" {
		_, err := strconv.Atoi(offset)
		if err == nil {
			sqlTpl += fmt.Sprintf(" OFFSET %s", offset)
		}
	}
	return sqlTpl
}

func prepareInsertSqlTpl(tableName string, fields []string) string {
	fieldsStr := ""
	placeholdersStr := ""
	lastIndex := len(fields) - 1
	if len(fields) > 0 {
		for i := 0; i < len(fields); i++ {
			fieldsStr += fmt.Sprintf("`%s`", fields[i])
			placeholdersStr += "?"
			if i != lastIndex {
				fieldsStr += ", "
				placeholdersStr += ", "
			}
		}
	}
	return fmt.Sprintf("INSERT INTO `%s`(%s) VALUES (%s)", tableName, fieldsStr, placeholdersStr)
}

func prepareUpdateSqlTpl(tableName string, fields []string, wheres map[string]string) string {
	fieldsStr := ""
	lastIndex := len(fields) - 1
	if len(fields) > 0 {
		for i := 0; i < len(fields); i++ {
			fieldsStr += fmt.Sprintf("`%s` = ?", fields[i])
			if i != lastIndex {
				fieldsStr += ", "
			}
		}
	}

	whereStr := ""
	i := 0
	if len(wheres) > 0 {
		whereStr += " WHERE "
		for field := range wheres {
			if i > 0 {
				whereStr += " AND "
			}
			whereStr += fmt.Sprintf("`%s` = ?", field)
			i++
		}
	}

	return fmt.Sprintf("UPDATE `%s` SET %s %s", tableName, fieldsStr, whereStr)
}

func prepareDeleteSqlTpl(tableName string, wheres map[string]string) string {
	whereStr := ""
	i := 0
	if len(wheres) > 0 {
		whereStr += " WHERE "
		for field := range wheres {
			if i > 0 {
				whereStr += " AND "
			}
			whereStr += fmt.Sprintf("`%s` = ?", field)
			i++
		}
	}

	return fmt.Sprintf("DELETE FROM `%s` %s", tableName, whereStr)
}
