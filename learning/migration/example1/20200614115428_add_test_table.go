package example1

import (
	"database/sql"
	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(Up20200614115428, Down20200614115428)
}

func Up20200614115428(tx *sql.Tx) error {
	_, err := tx.Exec("" +
		"CREATE TABLE IF NOT EXISTS `migrations_test` (" +
		"id    INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE," +
		"name  TEXT" +
		");",
	)
	if err != nil {
		return err
	}
	return nil
}

func Down20200614115428(tx *sql.Tx) error {
	_, err := tx.Exec("DROP TABLE IF EXISTS `migrations_test`")
	if err != nil {
		return err
	}
	return nil
}
