package example1

import (
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

type customError struct {
	err error
}

func (e customError) Error() string {
	return errors.Wrap(e.err, "custom err").Error()
}

func TestErr(t *testing.T) {
	err := getErr()

	switch err.(type) {
	case customError:
		assert.IsType(t, customError{}, err)
	default:
		assert.Fail(t, "type casting in switch failed")
	}

	if _, ok := err.(customError); !ok {
		assert.Fail(t, "type casting in type casting failed")
	}

}

func getErr() error {
	return customError{err: errors.New("init errors")}
}
