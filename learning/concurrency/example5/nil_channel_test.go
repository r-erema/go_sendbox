package example5

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

func Test(t *testing.T) {
	send := func(c chan<- int) {
		for {
			c <- rand.Intn(10)
		}
	}

	add := func(c <-chan int, result chan<- int) {
		sum := 0
		tc := time.NewTimer(time.Millisecond)

		for {
			select {
			case input := <-c:
				sum += input
			case <-tc.C:
				c = nil
				result <- sum
			}
		}
	}

	c, result := make(chan int), make(chan int)
	go add(c, result)
	go send(c)

	assert.Greater(t, <-result, 1000)
}
