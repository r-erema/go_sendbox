package example1

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_traverse(t *testing.T) {
	tree := create(10)
	result := traverse(tree)
	assert.NotEmpty(t, result)
}
