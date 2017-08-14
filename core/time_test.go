package core

import (
	"testing"

	"time"

	"github.com/stretchr/testify/assert"
)

func TestRealTimeManager_Observe(t *testing.T) {
	c := make(chan TimeEvent)
	rtm := NewRealTimeManager()
	rtm.Observe(c)

	ev := <-c
	assert.InDelta(t, RealTimeSampling, ev.elapsed, float64((10 * time.Millisecond).Nanoseconds()))
}
