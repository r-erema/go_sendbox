package example6

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func f1(cc chan chan int, f chan bool) {
	c := make(chan int)
	cc <- c
	defer close(c)

	sum := 0
	select {
	case x := <-c:
		for i := 0; i <= x; i++ {
			sum += i
		}
		c <- sum
	case <-f:
		return
	}
}

func Test(t *testing.T) {
	cc := make(chan chan int)
	result := make(map[int]int)
	for i := 0; i < 10; i++ {
		f := make(chan bool)
		go f1(cc, f)
		ch := <-cc
		ch <- i
		for sum := range ch {
			result[i] = sum
		}
		close(f)
	}
	assert.NotEmpty(t, result)
}
