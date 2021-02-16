package example2

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
	"time"
)

func TestBackgroundNotToDo(t *testing.T) {
	todo := fmt.Sprint(TODO())
	bg := fmt.Sprint(Background())
	assert.NotEqual(t, todo, bg)
}

func TestWithCancel(t *testing.T) {
	ctx, cancel := WithCancel(Background())
	assert.Nil(t, ctx.Err())
	cancel()

	<-ctx.Done()
	assert.Equal(t, ctx.Err(), Canceled)
}

func TestWithCancelConcurrent(t *testing.T) {
	ctx, cancel := WithCancel(Background())

	time.AfterFunc(time.Second*1, cancel)

	assert.Nil(t, ctx.Err())
	cancel()

	<-ctx.Done()
	assert.Equal(t, ctx.Err(), Canceled)
}

func TestWithCancelPropagation(t *testing.T) {
	ctxA, cancelA := WithCancel(Background())
	ctxB, cancelB := WithCancel(ctxA)
	defer cancelB()

	cancelA()

	select {
	case <-ctxB.Done():
	case <-time.After(time.Second):
		t.Errorf("time out")
	}

	assert.Equal(t, ctxB.Err(), Canceled)
}

func TestWithDeadline(t *testing.T) {
	deadline := time.Now().Add(2 * time.Second)
	ctx, cancel := WithDeadline(Background(), deadline)

	d, ok := ctx.Deadline()
	assert.True(t, ok)
	assert.Equal(t, d, deadline)

	then := time.Now()
	<-ctx.Done()
	d2 := time.Since(then)
	assert.Less(t, math.Abs(d2.Seconds())-2.0, 0.1)

	assert.Equal(t, ctx.Err(), DeadlineExceeded)

	cancel()

	assert.Equal(t, ctx.Err(), DeadlineExceeded)
}

func TestWithTimeout(t *testing.T) {
	timeout := 2 * time.Second
	deadline := time.Now().Add(timeout)
	ctx, cancel := WithTimeout(Background(), timeout)

	d, ok := ctx.Deadline()
	d.Sub(deadline)
	assert.True(t, ok)
	assert.True(t, d.Sub(deadline) < time.Millisecond)

	then := time.Now()
	<-ctx.Done()
	d2 := time.Since(then)
	assert.Less(t, math.Abs(d2.Seconds())-2.0, 0.1)

	assert.Equal(t, ctx.Err(), DeadlineExceeded)

	cancel()

	assert.Equal(t, ctx.Err(), DeadlineExceeded)
}

func TestWithValue(t *testing.T) {
	tc := []struct {
		key, val, keyRet, valRet interface{}
		shouldPanic              bool
	}{
		{"a", "b", "a", "b", false},
		{"a", "b", "c", nil, false},
		{42, true, 42, true, false},
		{42, true, int64(42), nil, false},
		{nil, true, nil, nil, true},
		{[]int{1, 2, 3}, true, []int{1, 2, 3}, nil, true},
	}

	for _, tt := range tc {
		var panicked interface{}
		func() {
			defer func() { panicked = recover() }()

			ctx := WithValue(Background(), tt.key, tt.val)
			if val := ctx.Value(tt.keyRet); val != tt.valRet {
				t.Errorf("expected value %v, got %v", tt.valRet, val)
			}
		}()

		if panicked != nil && !tt.shouldPanic {
			t.Errorf("unexpected panic: %v", panicked)
		}
		if panicked == nil && tt.shouldPanic {
			t.Errorf("expected panic, but didn't get it")
		}
	}
}
