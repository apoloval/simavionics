package core

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRealTimeManager_Observe(t *testing.T) {
	c := make(chan time.Time)
	before := time.Now()
	rtm := NewRealTimeManager()
	rtm.Observe(c)

	ts := <-c
	assert.Condition(t, func() bool {
		return before.Before(ts)
	})
}
