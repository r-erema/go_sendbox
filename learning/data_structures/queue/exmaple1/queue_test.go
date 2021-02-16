package exmaple1

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Push(t *testing.T) {
	Push(queue, 11)
	Push(queue, 8)
	assert.Equal(t, []float64{8, 11, 0}, traverse(queue))
}

func Test_Pop(t *testing.T) {
	Push(queue, 11)
	Push(queue, 8)
	v, res := Pop(queue)
	assert.True(t, res)
	assert.Equal(t, float64(0), v)
}
