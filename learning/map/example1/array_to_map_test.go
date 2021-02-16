package example1

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_arr2map(t *testing.T) {
	tests := []struct {
		arr         [4]string
		expectedMap map[int]string
	}{
		{
			arr:         [4]string{"a", "b", "c", "d"},
			expectedMap: map[int]string{0: "a", 1: "b", 2: "c", 3: "d"},
		},
		{
			arr:         [4]string{},
			expectedMap: map[int]string{0: "", 1: "", 2: "", 3: ""},
		},
		{
			arr:         [4]string{"a", "b"},
			expectedMap: map[int]string{0: "a", 1: "b", 2: "", 3: ""},
		},
		{
			arr:         [4]string{"a", "", "c"},
			expectedMap: map[int]string{0: "a", 1: "", 2: "c", 3: ""},
		},
	}

	for i, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("%d", i+1), func(t *testing.T) {
			m := arr2map(tt.arr)
			assert.Equal(t, tt.expectedMap, m)
		})
	}
}
