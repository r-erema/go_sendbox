package example3

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestMerge(t *testing.T) {

	valuesCount := 9
	ch1, ch2, ch3, mergeCh := make(chan int), make(chan int), make(chan int), make(chan int, valuesCount)
	exitCh := make(chan struct{})

	go func() {
		for {
			select {
			case v := <-ch1:
				mergeCh <- v
			case v := <-ch2:
				mergeCh <- v
			case v := <-ch3:
				mergeCh <- v
			case <-exitCh:
				close(mergeCh)
				return
			}
		}
	}()

	var wg sync.WaitGroup
	wg.Add(valuesCount)
	go func() { ch1 <- 1; wg.Done() }()
	go func() { ch1 <- 2; wg.Done() }()
	go func() { ch1 <- 3; wg.Done() }()
	go func() { ch2 <- 40; wg.Done() }()
	go func() { ch2 <- 50; wg.Done() }()
	go func() { ch2 <- 60; wg.Done() }()
	go func() { ch3 <- 700; wg.Done() }()
	go func() { ch3 <- 800; wg.Done() }()
	go func() { ch3 <- 900; wg.Done() }()
	wg.Wait()

	exitCh <- struct{}{}

	var result []int
	for v := range mergeCh {
		result = append(result, v)
	}

	assert.ElementsMatch(t, result, []int{1, 2, 3, 40, 50, 60, 700, 800, 900})

}
