package example2

import (
	"database/sql"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"go_sendbox/config"
	"log"
	"os"
	"testing"
)

var (
	db  *sql.DB
	err error
)

func TestMain(m *testing.M) {
	db, err = sql.Open("postgres", config.PostgresDSN())
	defer func() {
		_ = db.Close()
	}()
	if err != nil {
		log.Fatalln(err)
	}

	_, err = db.Exec("DROP TABLE IF EXISTS test")
	if err != nil {
		log.Fatalln(err)
	}
	_, err = db.Exec("CREATE TABLE test (name VARCHAR, val VARCHAR DEFAULT NULL)")
	if err != nil {
		log.Fatalln(err)
	}
	exitCode := m.Run()

	os.Exit(exitCode)
}

func TestSelectNone(t *testing.T) {
	defer clearTable()
	rows, err := db.Query("SELECT * FROM test")
	if err != nil {
		t.Fatal(err)
	}
	assert.False(t, rows.Next())

}

func TestSelect(t *testing.T) {
	defer clearTable()

	result, err := db.Exec("INSERT INTO test (name) VALUES ($1), ($2), ($3)", "1", "2", "3")
	if err != nil {
		t.Fatal(err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, int64(3), rowsAffected)

	rows, err := db.Query("SELECT COUNT(*) FROM test")
	if err != nil {
		t.Fatal(err)
	}
	var count int

	assert.True(t, rows.Next())
	_ = rows.Scan(&count)
	assert.Equal(t, 3, count)
}

func TestTransaction(t *testing.T) {
	defer clearTable()

	txn, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}

	stmt, err := txn.Prepare(pq.CopyIn("test", "name", "val"))
	if err != nil {
		t.Fatal(err)
	}

	for _, s := range []struct{ name, val string }{
		{name: "_n1_", val: "_v1_"},
		{name: "_n2_", val: "_v2_"},
	} {
		_, err = stmt.Exec(s.name, s.val)
		if err != nil {
			t.Fatal(err)
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		t.Fatal(err)
	}
	err = stmt.Close()
	if err != nil {
		t.Fatal(err)
	}

	err = txn.Commit()
	if err != nil {
		t.Fatal(err)
	}

	rows, err := db.Query("SELECT COUNT(*) FROM test")
	if err != nil {
		t.Fatal(err)
	}
	var count int

	assert.True(t, rows.Next())
	_ = rows.Scan(&count)
	assert.Equal(t, 2, count)

}

func TestConnectorWithNotify(t *testing.T) {
	baseConnector, err := pq.NewConnector(config.PostgresDSN())
	if err != nil {
		t.Fatal(err)
	}

	notificationCaught := false

	connector := pq.ConnectorWithNoticeHandler(baseConnector, func(notice *pq.Error) {
		notificationCaught = true
		assert.Equal(t, "test notice", notice.Message)
	})

	db := sql.OpenDB(connector)
	defer func() {
		_ = db.Close()
	}()

	sqlQuery := "DO language plpgsql $$ BEGIN RAISE NOTICE 'test notice'; END $$"
	if _, err := db.Exec(sqlQuery); err != nil {
		t.Fatal(err)
	}

	assert.True(t, notificationCaught)

}

func clearTable() {
	_, _ = db.Exec("TRUNCATE TABLE test")
}
