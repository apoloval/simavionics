package a320

import (
	"testing"

	"time"

	"github.com/apoloval/simavionics"
	"github.com/apoloval/simavionics/a320/apu"
	"github.com/apoloval/simavionics/event/local"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type APUTestSuite struct {
	suite.Suite
	simavionics.TimeAsserts
	bus simavionics.EventBus
	apu *APU
}

func (suite *APUTestSuite) SetupTest() {
	println("Setting up suite")
	suite.TimeAsserts = simavionics.NewTimeAsserts(suite.T())
	suite.bus = local.NewEventBus()
	ctx := simavionics.Context{suite.bus, suite.Dilation}
	suite.apu = NewAPU(ctx)
}

func (suite *APUTestSuite) TestMasterSwitchOn() {
	masterSwChan := suite.bus.Subscribe(apu.EventPower)
	flapOpenChan := suite.bus.Subscribe(apu.EventFlap)
	simavionics.PublishEvent(suite.bus, apu.EventMasterSwitch, true)

	ev := <-masterSwChan
	assert.Equal(suite.T(), true, ev.Bool())

	suite.AssertElapsed(6*time.Second, 1*time.Second, func() {
		ev = <-flapOpenChan
		assert.Equal(suite.T(), true, ev.Bool())
	})
}

func (suite *APUTestSuite) TestStartButtonPressed() {
	eventChanEngineN1 := suite.bus.Subscribe(apu.EventEngineN1)
	simavionics.PublishEvent(suite.bus, apu.EventMasterSwitch, true)
	simavionics.PublishEvent(suite.bus, apu.EventStartButton, true)

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
	suite.Run(t, new(APUTestSuite))
}
