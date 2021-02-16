package example2

import (
	"bufio"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"os"
	"sync"
	"testing"
	"time"
)

func TestConcurrentWrite(t *testing.T) {

	entries := make(map[string]int)

	var mu sync.Mutex
	count := func(f io.Reader) error {
		r := bufio.NewReader(f)
		for {
			if charRune, _, err := r.ReadRune(); err != nil {
				if err == io.EOF {
					return nil
				} else {
					return err
				}
			} else {
				char := string(charRune)
				mu.Lock()
				entries[char]++
				mu.Unlock()
			}
		}
	}

	files := []string{
		"file1",
		"file2",
		"file3",
		"file4",
	}

	start := time.Now()
	for _, file := range files {
		r, err := os.Open(file)
		require.NoError(t, err)

		err = count(r)
		require.NoError(t, err)
	}
	msSync := time.Since(start).Microseconds()

	entries = make(map[string]int)

	var wg sync.WaitGroup
	start = time.Now()
	for _, file := range files {
		wg.Add(1)
		go func(file string) {
			r, err := os.Open(file)
			require.NoError(t, err)

			err = count(r)
			require.NoError(t, err)
			wg.Done()
		}(file)
	}
	msConcurrent := time.Since(start).Microseconds()
	wg.Wait()

	assert.Less(t, msConcurrent, msSync)
}

func TestPipeline(t *testing.T) {
	files := []string{
		"file1",
		"file2",
		"file3",
		"file4",
	}

	count := func(f io.Reader) (map[string]int, error) {
		r := bufio.NewReader(f)
		result := make(map[string]int)
		for {
			if charRune, _, err := r.ReadRune(); err != nil {
				if err == io.EOF {
					return result, nil
				} else {
					return nil, err
				}
			} else {
				char := string(charRune)
				result[char]++
			}
		}
	}

	var channels []chan map[string]int
	for _, file := range files {
		ch := make(chan map[string]int)
		channels = append(channels, ch)
		go func(file string, ch chan map[string]int) {
			r, err := os.Open(file)
			require.NoError(t, err)

			result, err := count(r)
			require.NoError(t, err)
			ch <- result
			close(ch)
		}(file, ch)
	}

	time.Sleep(time.Second * 10)

	result := make(map[string]int)
	for _, ch := range channels {
		for r := range ch {
			for k, v := range r {
				result[k] += v
			}
		}
	}

	assert.NotEmpty(t, result)
}
