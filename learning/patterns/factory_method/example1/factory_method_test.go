package example1

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFactoryMethod(t *testing.T) {
	factory := LoggerFactory{}
	stdout := factory.create(LoggerTypeStdOut)
	file := factory.create(LoggerTypeFile)
	assert.IsType(t, &fileLogger{}, file)
	assert.IsType(t, &stdoutLogger{}, stdout)
}
