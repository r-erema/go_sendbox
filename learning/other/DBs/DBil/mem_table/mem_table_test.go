package mem_table

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestMemTable_sort(t *testing.T) {

	mt := MemTable{KVs: KVPairs{
		map[string]string{"java": "j"},
		map[string]string{"c": "c"},
		map[string]string{"php": "p"},
		map[string]string{"javascript": "js"},
		map[string]string{"golang": "g"},
	}}
	mt.sort()

	for i, kv := range mt.KVs {
		key := reflect.ValueOf(kv).MapKeys()[0].String()
		value := kv[key]
		switch i {
		case 0:
			assert.Equal(t, "c", key)
			assert.Equal(t, "c", value)
		case 1:
			assert.Equal(t, "golang", key)
			assert.Equal(t, "g", value)
		case 2:
			assert.Equal(t, "java", key)
			assert.Equal(t, "j", value)
		case 3:
			assert.Equal(t, "javascript", key)
			assert.Equal(t, "js", value)
		case 4:
			assert.Equal(t, "php", key)
			assert.Equal(t, "p", value)
		}
	}

}
