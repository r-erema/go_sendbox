package example2

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGotcha(t *testing.T) {
	A := make([]int, 0, 2)
	siblingA := append(A, 1)
	A = append(A, 7)
	A = append(A, 5)
	siblingA2 := append(siblingA, 3)
	assert.Equal(t, A, siblingA2)
	assert.NotEqual(t, A, siblingA)
	assert.NotEqual(t, siblingA2, siblingA)

	siblingA = append(siblingA, 11)

	assert.Equal(t, A, siblingA)
	assert.Equal(t, A, siblingA2)
	assert.Equal(t, siblingA, siblingA2)

	A = append(A, 35)
	A[0] = -3
	A[1] = -2
	A[2] = -1
	newSiblingA := append(A, 0)
	A = append(A, 0)
	newSiblingA[3] = 22
	assert.NotEqual(t, A, siblingA)
	assert.NotEqual(t, A, siblingA2)
	assert.NotEqual(t, newSiblingA, siblingA)
	assert.NotEqual(t, newSiblingA, siblingA2)
	assert.Equal(t, A, newSiblingA)

	siblingA[0] = 909
	siblingA2[1] = 340
	assert.Equal(t, siblingA, siblingA2)

}
