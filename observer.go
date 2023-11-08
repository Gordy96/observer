package observer

import (
	"sync"
	"sync/atomic"
)

func Observe[T comparable](val T) *Observer[T] {
	return &Observer[T]{val: val}
}

type Transition[T any] struct {
	From, To T
}

type listener[T any] struct {
	c      chan Transition[T]
	m      sync.Mutex
	closed bool
}

func (l *listener[T]) close() {
	l.m.Lock()
	defer l.m.Unlock()
	if !l.closed {
		close(l.c)
		l.closed = true
	}
}

func (l *listener[T]) notify(from, to T) {
	l.m.Lock()
	defer l.m.Unlock()
	if !l.closed {
		defer func() {
			_ = recover()
		}()
		l.c <- Transition[T]{from, to}
	}
}

type Observer[T comparable] struct {
	val       T
	mux       sync.RWMutex
	observers sync.Map
	n         atomic.Int64
}

func (o *Observer[T]) Is(other T) bool {
	o.mux.RLock()
	defer o.mux.RUnlock()
	return o.val == other
}

func (o *Observer[T]) Get() T {
	o.mux.RLock()
	defer o.mux.RUnlock()
	return o.val
}

func (o *Observer[T]) Set(val T) {
	o.mux.Lock()
	defer o.mux.Unlock()
	old := o.val
	o.val = val
	go o.notify(old, o.val)
}

func (o *Observer[T]) notify(from, to T) {
	o.observers.Range(func(key, value any) bool {
		go func(key int64, value *listener[T]) {
			value.notify(from, to)
		}(key.(int64), value.(*listener[T]))
		return true
	})
}

func (o *Observer[T]) unsubscribe(n int64) {
	i, ok := o.observers.LoadAndDelete(n)
	if ok {
		i.(*listener[T]).close()
	}
}

func (o *Observer[T]) Subscribe() (changes <-chan Transition[T], cancel func()) {
	n := o.n.Add(1)
	ch := make(chan Transition[T])
	o.observers.Store(n, &listener[T]{c: ch})
	return ch, func() {
		o.unsubscribe(n)
	}
}
