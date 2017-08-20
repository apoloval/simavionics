package apu

import (
	"testing"

	"time"

	"github.com/apoloval/simavionics/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type FlapTestSuite struct {
	suite.Suite
	core.TimeAsserts

	bus  *core.DefaultEventBus
	flap *Flap
}

func (suite *FlapTestSuite) SetupTest() {
	suite.TimeAsserts = core.NewTimeAsserts(suite.T())
	suite.bus = core.NewDefaultEventBus()
	ctx := core.SimContext{suite.bus, suite.Dilation}
	suite.flap = NewFlap(ctx)
}

func (suite *FlapTestSuite) TestOpenAndClose() {
	c := core.NewSubscription(suite.bus, StatusFlapOpen)
	suite.flap.Open()

	suite.AssertElapsed(6*time.Second, func() {
		event := <-c
		assert.Equal(suite.T(), true, event.Value)
	})

	suite.flap.Close()

	suite.AssertElapsed(10*time.Millisecond, func() {
		event := <-c
		assert.Equal(suite.T(), false, event.Value)
	})
}

func TestFlapTestSuite(t *testing.T) {
	suite.Run(t, new(FlapTestSuite))
}
