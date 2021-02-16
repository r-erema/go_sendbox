package example2

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_addNode(t *testing.T) {
	addNode(root, 7)
	addNode(root, 5)
	addNode(root, 2)
	assert.Equal(t, []int{0, 7, 5, 2}, traverse(root))
}

func Test_revers(t *testing.T) {
	addNode(root, 7)
	addNode(root, 5)
	addNode(root, 2)
	assert.Equal(t, []int{2, 5, 7}, reverse(root))
}

func Test_lookupNode(t *testing.T) {
	addNode(root, 7)
	addNode(root, 5)
	addNode(root, 2)

	assert.True(t, lookupNode(root, 5))
	assert.False(t, lookupNode(root, 54))
}

func Test_size(t *testing.T) {
	addNode(root, 7)
	addNode(root, 5)
	addNode(root, 2)
	assert.Equal(t, 4, size(root))
}
