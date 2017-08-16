package a320

import (
	"testing"

	"github.com/apoloval/simavionics/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type APUTestSuite struct {
	suite.Suite
	tm  *core.FakeTimeManager
	bus *core.DefaultEventBus
	ec  chan core.Event
	apu *APU
}

func (suite *APUTestSuite) SetupTest() {
	suite.tm = core.NewFakeTimeManager()
	suite.bus = core.NewDefaultEventBus()
	suite.ec = make(chan core.Event, 16)
	suite.apu = NewAPU(suite.tm, suite.bus)
}

func (suite *APUTestSuite) TestSwitchOn() {
	suite.bus.Subscribe(apuStateFlapOpen, suite.ec)
	suite.bus.Subscribe(apuStateMasterSwOn, suite.ec)
	suite.bus.Publish(core.Event{apuActionMasterSwOn, true})

	ev := <-suite.ec
	assert.Equal(suite.T(), apuStateFlapOpen, ev.Name)
	assert.Equal(suite.T(), true, ev.Bool())
	ev = <-suite.ec
	assert.Equal(suite.T(), apuStateMasterSwOn, ev.Name)
	assert.Equal(suite.T(), true, ev.Bool())
}

func TestAPUTestSuite(t *testing.T) {
	suite.Run(t, new(APUTestSuite))
}
