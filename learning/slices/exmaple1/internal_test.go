package exmaple1

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInternals(t *testing.T) {
	slice := []int{4, 7, 11, 18}
	data, err := underlyingArr(slice)
	assert.NoError(t, err)
	data[0] = 0

	for i := range slice {
		assert.True(t, data[i] == slice[i], "%d != %d", data[i], slice[i])
	}

	slice = append(slice, 29)
	assert.True(t, len(data) < len(slice))
}
