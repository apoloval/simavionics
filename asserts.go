package simavionics

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type TimeAsserts struct {
	t        *testing.T
	Delta    time.Duration
	Dilation TimeDilation
}

func NewTimeAsserts(t *testing.T) TimeAsserts {
	return TimeAsserts{
		t:        t,
		Delta:    10 * time.Millisecond,
		Dilation: 10,
	}
}

func (ta TimeAsserts) AssertElapsed(elapsed time.Duration, f func()) {
	t0 := time.Now()
	f()
	actualEnd := time.Now()
	expectedEnd := t0.Add(ta.Dilation.Dilated(elapsed))
	assert.WithinDuration(ta.t, expectedEnd, actualEnd, ta.Delta)
}
