package password

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_generatePassword(t *testing.T) {
	p := generatePassword(100)
	assert.Equal(t, 100, len(p))
}
