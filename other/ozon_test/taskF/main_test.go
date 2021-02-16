package taskF

import (
	"bufio"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

type Chunk struct {
	start, offset int64
}

func splitFileToChunks(fileName string, startOffset, chunkSize int64) (result []Chunk, err error) {

	file, err := os.Open(fileName)
	defer func() { _ = file.Close() }()
	if err != nil {
		return nil, err
	}

	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}
	size := stat.Size()
	if chunkSize > size {
		chunkSize = size
	}

	var currentOffset int64
	for i := int64(1); ; i++ {

		startOffset, _ = shiftTillSymbol(file.Name(), startOffset)
		if err == io.EOF {
			break
		}

		currentOffset = startOffset + chunkSize
		currentOffset, err = shiftTillSpace(file.Name(), currentOffset)
		result = append(
			result,
			Chunk{startOffset, currentOffset},
		)

		if err == io.EOF || currentOffset >= size {
			break
		}
		startOffset = currentOffset
	}
	return
}

func shiftTillSymbol(fileName string, offset int64) (int64, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return -1, err
	}
	defer func() { _ = file.Close() }()
	_, _ = file.Seek(offset, 0)
	buf := bufio.NewReader(file)
	for {
		b, err := buf.ReadByte()
		if err == io.EOF {
			break
		}
		if err != nil {
			return -1, err
		}
		if b != 32 {
			break
		}
		offset++
	}
	return offset, nil
}

func shiftTillSpace(fileName string, offset int64) (int64, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return -1, err
	}
	defer func() { _ = file.Close() }()
	_, _ = file.Seek(offset, 0)
	buf := bufio.NewReader(file)
	for {
		b, err := buf.ReadByte()
		if err == io.EOF {
			break
		}
		if err != nil {
			return -1, err
		}
		if b == 32 {
			break
		}
		offset++
	}
	return offset, nil
}

func parseFile(fileName string, chunks []Chunk, target int) (result []int, maxNumber int) {
	numbersMu, resultMu, wg := &sync.Mutex{}, &sync.Mutex{}, &sync.WaitGroup{}
	for _, chunk := range chunks {
		wg.Add(1)
		go func(chunk Chunk, numbersMu, resultMu *sync.Mutex, wg *sync.WaitGroup) {

			file, _ := os.Open(fileName)
			defer func() { _ = file.Close() }()
			_, _ = file.Seek(chunk.start, 0)
			buf := bufio.NewReader(file)
			var readOffset int64 = 0
			for {
				numberStr, err := buf.ReadString(' ')
				readOffset += int64(len(numberStr))
				if err == io.EOF && numberStr == "" {
					break
				}
				numberStr = strings.Trim(numberStr, " ")
				numberStr = strings.Trim(numberStr, "\n")

				number, err := strconv.Atoi(numberStr)
				if err == nil && target >= number {
					resultMu.Lock()
					result = append(result, number)
					if number > maxNumber {
						maxNumber = number
					}
					resultMu.Unlock()
				}

				if chunk.start+readOffset >= chunk.offset {
					break
				}

			}

			wg.Done()
		}(chunk, numbersMu, resultMu, wg)
	}
	wg.Wait()
	return
}

func getTarget(fileName string) (target int, currentOffset int64) {
	file, _ := os.Open(fileName)
	defer func() { _ = file.Close() }()
	var offset int64
	buf := bufio.NewReader(file)
	targetString, _ := buf.ReadString('\n')
	offset += int64(len(targetString))
	targetString = strings.Trim(targetString, "\n")
	target, _ = strconv.Atoi(targetString)
	return target, offset
}

func TaskF() error {

	fileName := "input.txt"
	target, offset := getTarget(fileName)

	chunks, _ := splitFileToChunks(fileName, offset, 1500000)

	updNumbers, maxNumber := parseFile(fileName, chunks, target)

	cutStart := time.Now()
	sort.Ints(updNumbers)
	stopIndex := -1
	for i, number := range updNumbers {
		if number+maxNumber < target {
			stopIndex = i
		} else {
			break
		}
	}
	if stopIndex != -1 {
		updNumbers = updNumbers[stopIndex+1:]
	}
	fmt.Println("Cut:", time.Since(cutStart))

	targetAttainedChan, targetNotAttainedChan := make(chan bool), make(chan bool)
	go func(targetAttainedChan, targetNotAttainedChan chan bool) {
		wg, mu := &sync.WaitGroup{}, &sync.Mutex{}
		handledNumbers := make(map[int]bool)
		for i, number := range updNumbers {

			mu.Lock()
			if _, ok := handledNumbers[number]; ok {
				mu.Unlock()
				continue
			}
			handledNumbers[number] = true
			mu.Unlock()

			wg.Add(1)
			go func(wg *sync.WaitGroup, number int, i int) {
				delta := target - number
				for _, nextNumber := range updNumbers[i+1:] {
					if delta == nextNumber {
						targetAttainedChan <- true
						return
					}
				}
				wg.Done()
			}(wg, number, i)
		}

		wg.Wait()
		targetNotAttainedChan <- true
	}(targetAttainedChan, targetNotAttainedChan)

	resultStart := time.Now()
	targetAttained := func(targetAttainedChan, targetNotAttainedChan chan bool) bool {
		select {
		case <-targetAttainedChan:
			return true
		case <-targetNotAttainedChan:
			return false
		}
	}(targetAttainedChan, targetNotAttainedChan)
	fmt.Println("Main work:", time.Since(resultStart))

	_ = os.Remove("output.txt")
	outputFile, _ := os.OpenFile("output.txt", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)

	if targetAttained {
		_, _ = outputFile.WriteString("1")
	} else {
		_, _ = outputFile.WriteString("0")
	}

	return nil
}

func TestTaskF(t *testing.T) {
	err := TaskF()
	if err != nil {
		t.Error(err)
	}

	result, err := ioutil.ReadFile("output.txt")
	if err != nil {
		panic(err)
	}

	assert.Equal(t, "1", string(result))
}
