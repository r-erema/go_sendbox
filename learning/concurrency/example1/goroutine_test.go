package example1

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"math/rand"
	"os"
	"sync"
	"testing"
)

func TestGoroutinesCount(t *testing.T) {

	os.Args = []string{"", "-n", "5"}
	n := flag.Int("n", 10, "Number of goroutines")
	flag.Parse()

	count := *n

	wg := &sync.WaitGroup{}
	mu := sync.Mutex{}
	var buffer bytes.Buffer
	wg.Add(count)
	for i := 1; i <= count; i++ {
		go func(n int) {
			mu.Lock()
			defer mu.Unlock()
			buffer.WriteString(fmt.Sprintf("goroutine %d\n", n))
			wg.Done()
		}(i)
	}
	wg.Wait()

	outputBytes, err := ioutil.ReadAll(&buffer)
	require.NoError(t, err)
	output := string(outputBytes)
	for i := 1; i <= count; i++ {
		assert.Contains(t, output, fmt.Sprintf("goroutine %d", i))
	}

}

func TestWriteToChan(t *testing.T) {

	c := make(chan int)
	go func(c chan<- int) {
		c <- 11
		c <- 12
		c <- 13
		c <- 14
		close(c)
	}(c)

	var result int
	for i := range c {
		result += i
	}

	assert.Equal(t, 50, result)

}

func TestPipelines(t *testing.T) {

	CLOSE := false
	DATA := make(map[int]bool)

	random := func(min, max int) int {
		return rand.Intn(max - min)
	}

	first := func(min, max int, out chan<- int) {
		for {
			if CLOSE {
				close(out)
				return
			}
			out <- random(min, max)
		}
	}

	second := func(in <-chan int, out chan<- int) {
		for x := range in {
			_, ok := DATA[x]
			if ok {
				CLOSE = true
			} else {
				DATA[x] = true
				out <- x
			}
		}
		close(out)
	}

	third := func(in <-chan int, out chan<- int) {
		for v := range in {
			out <- v
		}
		close(out)
	}

	out1 := make(chan int)
	go first(1, 3, out1)

	out2 := make(chan int)
	go second(out1, out2)

	out3 := make(chan int)
	go third(out2, out3)

	var result []int
	for v := range out3 {
		result = append(result, v)
	}

	assert.GreaterOrEqual(t, len(result), 1)
}
