package simavionics

type SimContext struct {
	Bus              EventBus
	RealTimeDilation TimeDilation
}
