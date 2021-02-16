package example1

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose"
	"github.com/stretchr/testify/assert"
	"testing"
)

const MigrationsDir = "."

func TestMigration(t *testing.T) {

	err := goose.SetDialect("sqlite3")
	if err != nil {
		t.Fatal(err)
	}
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err = db.Close()
		if err != nil {
			t.Fatal(err)
		}
	}()

	err = goose.Up(db, MigrationsDir)
	if err != nil {
		t.Fatal(err)
	}

	result, err := db.Exec("INSERT INTO `migrations_test` (`name`) VALUES ('1')")
	if err != nil {
		t.Fatal(err)
	}
	count, _ := result.RowsAffected()
	assert.Greater(t, count, int64(0))

	err = goose.Down(db, MigrationsDir)
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Query("INSERT INTO `migrations_test` (`name`) VALUES ('1')")
	assert.NotNil(t, err)

}
