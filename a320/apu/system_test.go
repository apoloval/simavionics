package apu

import (
	"testing"

	"time"

	"github.com/apoloval/simavionics"
	"github.com/apoloval/simavionics/event/local"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type SystemTestSuite struct {
	suite.Suite
	simavionics.TimeAsserts
	bus simavionics.EventBus
	apu *System
}

func (suite *SystemTestSuite) SetupTest() {
	println("Setting up suite")
	suite.TimeAsserts = simavionics.NewTimeAsserts(suite.T())
	suite.bus = local.NewEventBus()
	ctx := simavionics.Context{suite.bus, suite.Dilation}
	suite.apu = NewSystem(ctx)
}

func (suite *SystemTestSuite) TestMasterSwitchOn() {
	masterSwChan := suite.bus.Subscribe(EventPower)
	flapOpenChan := suite.bus.Subscribe(EventFlap)
	simavionics.PublishEvent(suite.bus, EventMasterSwitch, true)

	ev := <-masterSwChan
	assert.Equal(suite.T(), true, ev.Bool())

	suite.AssertElapsed(6*time.Second, 1*time.Second, func() {
		ev = <-flapOpenChan
		assert.Equal(suite.T(), true, ev.Bool())
	})
}

func (suite *SystemTestSuite) TestStartButtonPressed() {
	eventChanEngineN1 := suite.bus.Subscribe(EventEngineN1)
	simavionics.PublishEvent(suite.bus, EventMasterSwitch, true)
	simavionics.PublishEvent(suite.bus, EventStartButton, true)

	suite.AssertElapsed(60*time.Second, 10*time.Second, func() {
		for {
			ev := <-eventChanEngineN1
			if ev.Float64() >= 100.0 {
				return
			}
		}
	})
}

func TestAPUTestSuite(t *testing.T) {
	suite.Run(t, new(SystemTestSuite))
}
