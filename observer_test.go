package observer

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestObserve(t *testing.T) {
	v := Observe(42)
	assert.NotNil(t, v)
	assert.Equal(t, 42, v.Get())
	assert.True(t, v.Is(42))

	ch, cancelCh := v.Subscribe()
	go func() {
		time.Sleep(10 * time.Millisecond)
		v.Set(55)
	}()
	ex1 := <-ch
	cancelCh()
	assert.Equal(t, Transition[int]{42, 55}, ex1)
	ch, cancelCh = v.Subscribe()
	cancelCh()
	_, ok := <-ch
	assert.False(t, ok)
}
