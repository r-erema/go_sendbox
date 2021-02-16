package example1

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	N4 int = 4 << iota
	N8
	N16
	N32
)

func Test_getDayNumber(t *testing.T) {
	tests := []struct {
		constant int
		value    int
	}{
		{constant: N4, value: 4},
		{constant: N8, value: 8},
		{constant: N16, value: 16},
		{constant: N32, value: 32},
	}

	for i, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			assert.Equal(t, tt.value, tt.constant)
		})
	}
}
