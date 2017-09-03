package apu

import (
	"testing"

	"time"

	"github.com/apoloval/simavionics"
	"github.com/apoloval/simavionics/event/local"
	"github.com/stretchr/testify/assert"
)

func TestNewBleed(t *testing.T) {
	tests := []struct {
		name          string
		n1            float64
		bleedSwitch   bool
		expectedBleed float64
	}{
		{
			name:          "Bleed PSI is zero if APU is not running",
			n1:            0.0,
			bleedSwitch:   true,
			expectedBleed: 0.0,
		},
		{
			name:          "Bleed PSI is half if APU is half running",
			n1:            50.0,
			bleedSwitch:   true,
			expectedBleed: 17.50,
		},
		{
			name:          "Bleed PSI is max if APU is full running",
			n1:            100.0,
			bleedSwitch:   true,
			expectedBleed: 35.0,
		},
		{
			name:          "Bleed PSI is zero if APU is not running with switch off",
			n1:            0.0,
			bleedSwitch:   false,
			expectedBleed: 0.0,
		},
		{
			name:          "Bleed PSI is zero if APU is half running with switch off",
			n1:            50.0,
			bleedSwitch:   false,
			expectedBleed: 0.0,
		},
		{
			name:          "Bleed PSI is zero if APU is full running with switch off",
			n1:            100.0,
			bleedSwitch:   false,
			expectedBleed: 0.0,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			bus := local.NewEventBus()
			NewBleed(bus)
			eventChanBleed := bus.Subscribe(EventBleed)

			simavionics.PublishEvent(bus, EventEngineN1, test.n1)
			simavionics.PublishEvent(bus, EventBleedSwitch, test.bleedSwitch)

			time.Sleep(20 * time.Millisecond)
			var v simavionics.EventValue
			for i := 0; i <= len(eventChanBleed); i++ {
				v = <-eventChanBleed
			}
			assert.Equal(t, test.expectedBleed, v.Float64())
		})
	}
}
