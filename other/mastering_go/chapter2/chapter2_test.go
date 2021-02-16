package chapter2

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestC(t *testing.T) {
	start := time.Now()
	cResult := sumC(10000)
	execTimeC := time.Since(start)

	start = time.Now()
	goResult := sumGo(10000)
	execTimeGo := time.Since(start)

	assert.Equal(t, cResult, goResult)
	assert.Greater(t, execTimeGo.Nanoseconds(), execTimeC.Nanoseconds())
}
