package apu

import (
	"testing"

	"time"

	"github.com/apoloval/simavionics"
	"github.com/apoloval/simavionics/event/local"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type FlapTestSuite struct {
	suite.Suite
	simavionics.TimeAsserts

	bus  simavionics.EventBus
	flap *Flap
}

func (suite *FlapTestSuite) SetupTest() {
	suite.TimeAsserts = simavionics.NewTimeAsserts(suite.T())
	suite.bus = local.NewEventBus()
	ctx := simavionics.Context{suite.bus, suite.Dilation}
	suite.flap = NewFlap(ctx)
}

func (suite *FlapTestSuite) TestOpenAndClose() {
	c := suite.bus.Subscribe(StatusFlapOpen)
	suite.flap.Open()

	suite.AssertElapsed(6*time.Second, func() {
		event := <-c
		assert.Equal(suite.T(), true, event.Bool())
	})

	suite.flap.Close()

	suite.AssertElapsed(10*time.Millisecond, func() {
		event := <-c
		assert.Equal(suite.T(), false, event.Bool())
	})
}

func TestFlapTestSuite(t *testing.T) {
	suite.Run(t, new(FlapTestSuite))
}
