package example2

import (
	"errors"
	"reflect"
	"sync"
	"time"
)

type Context interface {
	Deadline() (deadline time.Time, ok bool)
	Done() <-chan struct{}
	Err() error
	Value(key interface{}) interface{}
}

type emptyCtx int

func (emptyCtx) Deadline() (deadline time.Time, ok bool) {
	return
}
func (emptyCtx) Done() <-chan struct{} {
	return nil
}
func (emptyCtx) Err() error {
	return nil
}
func (emptyCtx) Value(key interface{}) interface{} {
	return nil
}

var (
	background = new(emptyCtx)
	todo       = new(emptyCtx)
)

func TODO() Context {
	return todo
}
func Background() Context {
	return background
}

type cancelCtx struct {
	Context
	done chan struct{}
	err  error
	mu   sync.Mutex
}

func (c *cancelCtx) Done() <-chan struct{} {
	return c.done
}
func (c *cancelCtx) Err() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.err
}
func (c *cancelCtx) cancel(err error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.err != nil {
		return
	}
	c.err = err
	close(c.done)
}

var Canceled = errors.New("ctx cancelled")

type CancelFunc func()

func WithCancel(parentCtx Context) (Context, CancelFunc) {
	ctx := &cancelCtx{
		Context: parentCtx,
		done:    make(chan struct{}),
	}
	cancel := func() { ctx.cancel(Canceled) }

	go func() {
		select {
		case <-parentCtx.Done():
			ctx.cancel(parentCtx.Err())
		case <-ctx.Done():
		}
	}()

	return ctx, cancel
}

type deadlineCtx struct {
	*cancelCtx
	deadline time.Time
}

func (ctx *deadlineCtx) Deadline() (deadline time.Time, ok bool) {
	return ctx.deadline, true
}

var DeadlineExceeded = errors.New("deadline exceeded")

func WithDeadline(parentCtx Context, deadline time.Time) (Context, CancelFunc) {
	cCtx, cancel := WithCancel(parentCtx)
	ctx := &deadlineCtx{
		cancelCtx: cCtx.(*cancelCtx),
		deadline:  deadline,
	}

	t := time.AfterFunc(time.Until(deadline), func() {
		ctx.cancel(DeadlineExceeded)
	})

	stop := func() {
		t.Stop()
		cancel()
	}

	return ctx, stop
}
func WithTimeout(parentCtx Context, timeout time.Duration) (Context, CancelFunc) {
	return WithDeadline(parentCtx, time.Now().Add(timeout))
}

type valueCtx struct {
	Context
	key, value interface{}
}

func (ctx *valueCtx) Value(key interface{}) interface{} {
	if key == ctx.key {
		return ctx.value
	}
	return ctx.Context.Value(key)
}
func WithValue(parentCtx Context, key, value interface{}) Context {
	if key == nil {
		panic("key is nil")
	}
	if !reflect.TypeOf(key).Comparable() {
		panic("key is not comparable")
	}
	return &valueCtx{parentCtx, key, value}
}
