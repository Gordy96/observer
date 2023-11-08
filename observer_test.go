package observer

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestObserve(t *testing.T) {
	v := Observe(42)
	assert.NotNil(t, v)
	assert.Equal(t, 42, v.Get())
	assert.True(t, v.Is(42))
}
