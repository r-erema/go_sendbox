package config

import (
	"fmt"
	"os"
)

func MysqlDSN() string {
	if host, ok := os.LookupEnv("MYSQL_HOST"); ok {
		return fmt.Sprintf("root:123@tcp(%s:3306)/go?charset=utf8", host)
	}
	return "root:123@tcp(localhost:3306)/go?charset=utf8"
}

func PostgresDSN() string {
	if host, ok := os.LookupEnv("POSTGRES_HOST"); ok {
		return fmt.Sprintf("postgres://go:123@%s:5432/go?sslmode=disable", host)
	}
	return "postgres://go:123@localhost:5432/go?sslmode=disable"
}

func Neo4jDSN() string {
	if host, ok := os.LookupEnv("NEO4J_HOST"); ok {
		return fmt.Sprintf("bolt://%s:7687", host)
	}
	return "bolt://localhost:7687"
}
