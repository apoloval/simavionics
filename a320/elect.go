package a320

type GenState struct {
	Current    float64
	MaxCurrent float64
	Voltage    float64
	Freq       float64
}

func (gs GenState) CurrentPercentage() float64 {
	return gs.Current / gs.MaxCurrent
}
