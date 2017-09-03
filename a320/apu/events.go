package apu

const (
	// APU controls in overhead panel
	EventBleedSwitch  = "oh/apu/bleed"  // [bool] The bleed switch
	EventMasterSwitch = "oh/apu/master" // [bool] The master switch status
	EventPower        = "oh/apu/power"  // [bool] The APU is electrified
	EventStartButton  = "oh/apu/start"  // [bool] The APU start button is pressed

	// From engine.go
	EventBleed     = "apu/bleed"   // [Float64] The pressure of APU bleed in PSI
	EventEngineN1  = "apu/eng/n1"  // [Float64] The speed of N1 from 0.0 to 100.0
	EventEngineEGT = "apu/eng/egt" // [Float64] The escape gasses temperature in ºC

	// From flap.go
	EventFlap = "apu/flap" // [bool] The status of the flap
)
