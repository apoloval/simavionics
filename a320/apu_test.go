package a320

import (
	"testing"

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
	suite.TimeAsserts = simavionics.NewTimeAsserts(suite.T())
	suite.bus = local.NewEventBus()
	ctx := simavionics.Context{suite.bus, suite.Dilation}
	suite.apu = NewAPU(ctx)
}

func (suite *APUTestSuite) TestSwitchOn() {
	masterSwChan := suite.bus.Subscribe(ApuStateMasterSwOn)
	flapOpenChan := suite.bus.Subscribe(apu.StatusFlapOpen)
	suite.bus.Publish(ApuActionMasterSwOn, true)

	ev := <-masterSwChan
	assert.Equal(suite.T(), true, ev.Bool())

	suite.AssertElapsed(apuFlapOpenTime, func() {
		ev = <-flapOpenChan
		assert.Equal(suite.T(), true, ev.Bool())
	})
}

func TestAPUTestSuite(t *testing.T) {
	suite.Run(t, new(APUTestSuite))
}
