package main

import (
	"time"

	"github.com/apoloval/simavionics"
	"github.com/apoloval/simavionics/a320"
	"github.com/apoloval/simavionics/event/remote"
)

func main() {
	bus, err := remote.NewMasterEventBus("tcp://localhost:7001")
	if err != nil {
		panic(err)
	}

	cli, err := NewCLI(bus)
	if err != nil {
		panic(err)
	}

	ctx := simavionics.Context{bus, 1}
	simavionics.NewHeartbeat(bus, 250*time.Millisecond)
	a320.NewAPU(ctx)
	cli.Run()
}
