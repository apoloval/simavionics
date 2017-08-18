package a320

import (
	"testing"

	"github.com/apoloval/simavionics/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type APUTestSuite struct {
	suite.Suite
	bus      *core.DefaultEventBus
	consumer core.EventBusConsumer
	apu      *APU
}

func (suite *APUTestSuite) SetupTest() {
	suite.bus = core.NewDefaultEventBus()
	suite.consumer = core.NewEventBusConsumer(suite.bus, 16)
	ctx := core.SimContext{suite.bus, 1000}
	suite.apu = NewAPU(ctx)
}

func (suite *APUTestSuite) TestSwitchOn() {
	suite.consumer.Subscribe(apuStateFlapOpen)
	suite.consumer.Subscribe(apuStateMasterSwOn)
	suite.bus.Publish(core.Event{apuActionMasterSwOn, true})

	ev := suite.consumer.Consume()
	assert.Equal(suite.T(), apuStateMasterSwOn, ev.Name)
	assert.Equal(suite.T(), true, ev.Bool())
	ev = suite.consumer.Consume()
	assert.Equal(suite.T(), apuStateFlapOpen, ev.Name)
	assert.Equal(suite.T(), true, ev.Bool())
}

func TestAPUTestSuite(t *testing.T) {
	suite.Run(t, new(APUTestSuite))
}
