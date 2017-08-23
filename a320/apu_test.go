package a320

import (
	"testing"

	"github.com/apoloval/simavionics/a320/apu"
	"github.com/apoloval/simavionics/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type APUTestSuite struct {
	suite.Suite
	core.TimeAsserts
	bus *core.LocalEventBus
	apu *APU
}

func (suite *APUTestSuite) SetupTest() {
	suite.TimeAsserts = core.NewTimeAsserts(suite.T())
	suite.bus = core.NewDefaultEventBus()
	ctx := core.SimContext{suite.bus, suite.Dilation}
	suite.apu = NewAPU(ctx)
}

func (suite *APUTestSuite) TestSwitchOn() {
	masterSwChan := suite.bus.Subscribe(apuStateMasterSwOn)
	flapOpenChan := suite.bus.Subscribe(apu.StatusFlapOpen)
	suite.bus.Publish(apuActionMasterSwOn, true)

	ev := <-masterSwChan
	assert.Equal(suite.T(), true, ev.(bool))

	suite.AssertElapsed(apuFlapOpenTime, func() {
		ev = <-flapOpenChan
		assert.Equal(suite.T(), true, ev.(bool))
	})
}

func TestAPUTestSuite(t *testing.T) {
	suite.Run(t, new(APUTestSuite))
}
