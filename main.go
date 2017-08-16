package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/apoloval/simavionics/a320"
	"github.com/apoloval/simavionics/core"
)

func main() {
	log.Printf("[main] Initializing A320 simulator")
	rtm := core.NewRealTimeManager()
	bus := core.NewDefaultEventBus()
	a320.NewAPU(rtm, bus)
	waitForStopSignal()
}

func waitForStopSignal() {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigchan
	log.Printf("[main] Received stop signal: %s", sig.String())
}
