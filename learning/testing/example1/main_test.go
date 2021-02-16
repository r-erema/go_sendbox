package example1

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCalculate(t *testing.T) {
	assert.Equal(t, Calculate(2), 4)
}
