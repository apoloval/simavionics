package apu

import (
	"testing"

	"time"

	"github.com/apoloval/simavionics"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type FlapTestSuite struct {
	suite.Suite
	simavionics.TimeAsserts

	bus  *simavionics.LocalEventBus
	flap *Flap
}

func (suite *FlapTestSuite) SetupTest() {
	suite.TimeAsserts = simavionics.NewTimeAsserts(suite.T())
	suite.bus = simavionics.NewDefaultEventBus()
	ctx := simavionics.SimContext{suite.bus, suite.Dilation}
	suite.flap = NewFlap(ctx)
}

func (suite *FlapTestSuite) TestOpenAndClose() {
	c := suite.bus.Subscribe(StatusFlapOpen)
	suite.flap.Open()

	suite.AssertElapsed(6*time.Second, func() {
		event := <-c
		assert.Equal(suite.T(), true, event)
	})

	suite.flap.Close()

	suite.AssertElapsed(10*time.Millisecond, func() {
		event := <-c
		assert.Equal(suite.T(), false, event)
	})
}

func TestFlapTestSuite(t *testing.T) {
	suite.Run(t, new(FlapTestSuite))
}
