package main

import (
	"time"

	"log"

	"github.com/apoloval/simavionics"
)

type ConnectionWatcher struct {
	hbListener *simavionics.HeartbeatListener
}

func NewConnectionWatcher(bus simavionics.EventBus) *ConnectionWatcher {
	w := &ConnectionWatcher{
		hbListener: simavionics.NewHeartbeatListener(bus, 1*time.Second),
	}
	go w.watch()
	return w
}

func (w *ConnectionWatcher) watch() {
	var was_alive bool
	for {
		alive := <-w.hbListener.AliveChan()
		if alive != was_alive {
			if alive {
				log.Println("The simulator is now online")
			} else {
				log.Println("The simulator is now offline")
			}
			was_alive = alive
		}
	}
}
