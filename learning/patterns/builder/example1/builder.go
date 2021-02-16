package example1

import (
	"fmt"
	"strings"
)

type Filters map[string]string
type QueryBuilder interface {
	Select(table string) QueryBuilder
	Where(filters Filters) QueryBuilder
	Limit(limit int, offset int) QueryBuilder
	BuildQuery() string
}
type SQLBuilder struct {
	table         string
	filters       string
	limit, offset int
}

func (S *SQLBuilder) Select(table string) QueryBuilder {
	S.table = table
	return S
}
func (S *SQLBuilder) Where(filters Filters) QueryBuilder {
	var pairs []string
	for field, value := range filters {
		pairs = append(pairs, fmt.Sprintf("`%s` = '%s'", field, value))
	}
	S.filters = strings.Join(pairs, " AND ")
	return S
}
func (S *SQLBuilder) Limit(limit int, offset int) QueryBuilder {
	S.limit = limit
	S.offset = offset
	return S
}
func (S *SQLBuilder) BuildQuery() string {
	return fmt.Sprintf(
		"SELECT * FROM %s WHERE %s LIMIT %d, %d",
		S.table,
		S.filters,
		S.limit,
		S.offset,
	)
}

type MongoBuilder struct {
	table         string
	filters       string
	limit, offset int
}

func (m *MongoBuilder) Select(table string) QueryBuilder {
	m.table = table
	return m
}
func (m *MongoBuilder) Where(filters Filters) QueryBuilder {
	var pairs []string
	for field, value := range filters {
		pairs = append(pairs, fmt.Sprintf(`%s: "%s"`, field, value))
	}
	m.filters = fmt.Sprintf("{%s}", strings.Join(pairs, ", "))
	return m
}
func (m *MongoBuilder) Limit(limit int, offset int) QueryBuilder {
	m.limit = limit
	m.offset = offset
	return m
}
func (m *MongoBuilder) BuildQuery() string {
	return fmt.Sprintf(
		"db.%s.find(%s).limit(%d).skip(%d)",
		m.table,
		m.filters,
		m.limit,
		m.offset,
	)
}
