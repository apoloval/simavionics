package apu

import (
	"testing"

	"time"

	"github.com/apoloval/simavionics"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/apoloval/simavionics/event/local"
)

type FlapTestSuite struct {
	suite.Suite
	simavionics.TimeAsserts

	bus  simavionics.EventBus
	flap *flap
}

func (suite *FlapTestSuite) SetupTest() {
	suite.TimeAsserts = simavionics.NewTimeAsserts(suite.T())
	suite.bus = local.NewEventBus()
	ctx := simavionics.Context{suite.bus, suite.Dilation}
	suite.flap = newFlap(ctx)
}

func (suite *FlapTestSuite) TestOpenAndClose() {
	c := suite.bus.Subscribe(EventFlap)
	suite.flap.open()

	suite.AssertElapsed(6*time.Second, 1*time.Second, func() {
		event := <-c
		assert.Equal(suite.T(), true, event.Bool())
	})

	suite.flap.close()

	suite.AssertElapsed(1*time.Millisecond, 2*time.Millisecond, func() {
		event := <-c
		assert.Equal(suite.T(), false, event.Bool())
	})
}

func TestFlapTestSuite(t *testing.T) {
	suite.Run(t, new(FlapTestSuite))
}
