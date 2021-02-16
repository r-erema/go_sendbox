package example1

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuilder(t *testing.T) {
	sqlQB := SQLBuilder{}
	sqlQuery := sqlQB.Select("users").Where(map[string]string{"age": "17", "sex": "female"}).Limit(5, 0).BuildQuery()
	mongoQB := MongoBuilder{}
	mongoQuery := mongoQB.Select("users").Where(map[string]string{"age": "17", "sex": "female"}).Limit(5, 0).BuildQuery()
	assert.Equal(t, "SELECT * FROM users WHERE `age` = '17' AND `sex` = 'female' LIMIT 5, 0", sqlQuery)
	assert.Equal(t, `db.users.find({age: "17", sex: "female"}).limit(5).skip(0)`, mongoQuery)
}
