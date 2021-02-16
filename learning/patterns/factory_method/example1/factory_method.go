package example1

import (
	"fmt"
	"os"
)

type logger interface {
	log(message string)
}
type stdoutLogger struct{}

func (s *stdoutLogger) log(message string) {
	fmt.Println(message)
}

type fileLogger struct {
	file os.File
}

func (f *fileLogger) log(message string) {
	_, _ = f.file.Write([]byte(message))
}

const LoggerTypeStdOut = "stdout"
const LoggerTypeFile = "file"

type LoggerFactory struct{}

func (lf LoggerFactory) create(loggerType string) logger {
	switch loggerType {
	case LoggerTypeFile:
		return &fileLogger{}
	case LoggerTypeStdOut:
		return &stdoutLogger{}
	}
	return &stdoutLogger{}
}
