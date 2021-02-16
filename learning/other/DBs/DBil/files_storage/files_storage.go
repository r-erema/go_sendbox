package files_storage

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	dataFilePath    = "./data/"
	dataFilePattern = dataFilePath + "*.bin"
)

type DataFilesStorage struct {
	files []*os.File
}

func NewDataFilesStorage() *DataFilesStorage {
	return &DataFilesStorage{files: prepareDataFiles()}
}

func (s *DataFilesStorage) Files() []*os.File {
	return s.files
}

func (s *DataFilesStorage) LastFile() *os.File {
	return (s.files)[len(s.files)-1]
}

func (s *DataFilesStorage) AppendNewFile() {
	s.files = append(s.files, createNewDataFile())
}

func (s *DataFilesStorage) GetFileByName(fileName string) *os.File {
	for _, file := range s.files {
		if file.Name() == fileName {
			return file
		}
	}
	return nil
}

func (s *DataFilesStorage) Reset() error {
	matches, _ := filepath.Glob(dataFilePattern)
	for _, match := range matches {
		err := os.RemoveAll(match)
		if err != nil {
			return err
		}
	}
	s.files = prepareDataFiles()
	return nil
}

func prepareDataFiles() (dataFiles []*os.File) {
	matches, _ := filepath.Glob(dataFilePattern)
	filesCount := len(matches)
	if filesCount > 0 {
		dataFiles = make([]*os.File, filesCount)
		for i, match := range matches {
			f, _ := os.OpenFile(match, os.O_APPEND|os.O_RDWR, 0644)
			dataFiles[i] = f
		}
	} else {
		dataFiles = make([]*os.File, 1)
		dataFiles[0] = createNewDataFile()
	}
	return
}

func createNewDataFile() *os.File {
	file, _ := os.OpenFile(
		strings.Replace(dataFilePattern, "*", fmt.Sprintf("%d", time.Now().UnixNano()), -1),
		os.O_APPEND|os.O_CREATE|os.O_RDWR,
		0644,
	)
	return file
}
