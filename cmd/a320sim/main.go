package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"time"

	"github.com/apoloval/simavionics"
	"github.com/apoloval/simavionics/a320"
	"github.com/apoloval/simavionics/event/remote"
)

func main() {
	log.Printf("[main] Initializing A320 simulator")

	bus, err := remote.NewMasterEventBus("tcp://localhost:7001")
	if err != nil {
		panic(err)
	}

	ctx := simavionics.Context{bus, 1}
	simavionics.NewHeartbeat(bus, 250*time.Millisecond)
	apusys := a320.NewAPU(ctx)
	apusys.Start()
	waitForStopSignal()
}

func waitForStopSignal() {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigchan
	log.Printf("[main] Received stop signal: %s", sig.String())
}
