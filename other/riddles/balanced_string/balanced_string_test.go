package balanced_string

import (
	"github.com/stretchr/testify/assert"

	"testing"
)

func Test_isBalanced(t *testing.T) {
	assert.True(t, isBalanced(""))
	assert.False(t, isBalanced("{]}"))
	assert.False(t, isBalanced("{]]]}"))
	assert.True(t, isBalanced("{}"))
	assert.True(t, isBalanced("{[]}"))
	assert.False(t, isBalanced("{[{]}"))
	assert.True(t, isBalanced("{[{}]}"))
	assert.True(t, isBalanced("([{{{}}}]{{}[]})"))
	assert.True(t, isBalanced("{(({})()){}[(())]}[][]"))
}
