package taskB

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestTaskB(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Error(err)
	}
	sqlBytes, err := ioutil.ReadFile("./dump.sql")
	if err != nil {
		t.Error(err)
	}
	_, err = db.Exec(string(sqlBytes))
	if err != nil {
		t.Error(err)
	}

	rows, err := db.Query("SELECT g.id, g.name FROM " +
		"(SELECT g.id, g.name, COUNT(tg.tag_id) tc FROM goods g JOIN tags_goods tg on g.id = tg.goods_id GROUP BY g.id) as g " +
		"WHERE g.tc >= (SELECT COUNT(*) FROM tags)")
	if err != nil {
		t.Error(err)
	}

	var id int
	var name string

	neededIds := []int{2, 4}
	for rows.Next() {
		err = rows.Scan(&id, &name)
		assert.Contains(t, neededIds, id)
		if err != nil {
			t.Error(err)
		}
	}
}
