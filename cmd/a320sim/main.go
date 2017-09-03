package main

import (
	"time"

	"os"
	"os/signal"
	"syscall"

	"github.com/apoloval/simavionics"
	"github.com/apoloval/simavionics/a320"
	"github.com/apoloval/simavionics/event/remote"
)

func main() {
	simavionics.EnableLogging()

	bus, err := remote.NewMasterEventBus("tcp://localhost:7001")
	if err != nil {
		panic(err)
	}

	ctx := simavionics.Context{bus, 1}
	simavionics.NewHeartbeat(bus, 250*time.Millisecond)
	a320.NewAPU(ctx)

	waitForStopSignal()
}

func waitForStopSignal() {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	<-sigchan
}
