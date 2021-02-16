package example1

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRecover(t *testing.T) {
	defer func() {
		err := recover()
		assert.NotNil(t, err)
	}()
	panic("some error")
}
