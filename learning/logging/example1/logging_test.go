package example1

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log"
	"os"
	"testing"
)

var LOGFILE = "./log.log"

func TestLogging(t *testing.T) {
	f, err := os.OpenFile(LOGFILE, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, f.Close())
	}()

	iLog := log.New(f, "customLogLineNumber", log.LstdFlags|log.Lshortfile)
	iLog.Println("Hello there!")
	iLog.Println("Another log entry!")

	f, err = os.Open(LOGFILE)
	require.NoError(t, err)
	stat, err := f.Stat()
	require.NoError(t, err)
	assert.Equal(t, stat.Size(), int64(150))

	err = os.Remove(LOGFILE)
	require.NoError(t, err)

}
