package exmaple1

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Push(t *testing.T) {
	Push(NodeVal{Value: 11})
	Push(NodeVal{Value: 8})
	assert.Equal(t, []int{8, 11, 0}, traverse(stack))
}

func Test_Pop(t *testing.T) {
	Push(NodeVal{Value: 11})
	Push(NodeVal{Value: 8})
	v, res := Pop(stack)
	assert.True(t, res)
	assert.Equal(t, NodeVal{Value: 8}, v)
}
