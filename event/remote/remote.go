package remote

import (
	"log"

	"github.com/apoloval/simavionics"
	"github.com/go-mangos/mangos"
	"github.com/go-mangos/mangos/protocol/bus"
	"github.com/go-mangos/mangos/transport/ipc"
	"github.com/go-mangos/mangos/transport/tcp"
)

type eventBus struct {
	masterNode    bool
	subs          *simavionics.EventSubscriptions
	socket        mangos.Socket
	localPubChan  chan *simavionics.Event
	remotePubChan chan *simavionics.Event
}

func (bus *eventBus) Subscribe(ev simavionics.EventName) <-chan simavionics.EventValue {
	return bus.subs.New(ev)
}

func (bus *eventBus) Publish(ev simavionics.Event) {
	bus.localPubChan <- &ev
	bus.remotePubChan <- &ev
}

func (bus *eventBus) localPublisher() {
	for {
		select {
		case e := <-bus.localPubChan:
			if e == nil {
				return
			}
			bus.publishLocal(e)
		}
	}
}

func (bus *eventBus) remotePublisher() {
	for {
		select {
		case e := <-bus.remotePubChan:
			if e == nil {
				return
			}
			bus.publishRemote(e)
		}
	}
}

func (bus *eventBus) receiver() {
	for {
		var ev *simavionics.Event
		msg, err := bus.socket.RecvMsg()
		if err != nil {
			log.Printf("[bus.remote] Failed to receive message from socket: %v", err)
		}

		ev, err = simavionics.DecodeEvent(msg.Body)
		if err != nil {
			log.Printf("[bus.remote] Failed to decode message from socket: %v", err)
		}
		bus.publishLocal(ev)

		if bus.masterNode {
			bus.socket.SendMsg(msg)
		}
	}
}

func (bus *eventBus) publishLocal(event *simavionics.Event) {
	bus.subs.Each(event.Name, func(c chan simavionics.EventValue) {
		c <- event.Value
	})
}

func (bus *eventBus) publishRemote(event *simavionics.Event) {
	if err := bus.socket.Send(event.Encode()); err != nil {
		log.Printf("[bus.remote]: failed to publish on socket %v", bus.socket)
	}
}

func NewEventBus(masterNode bool, addrs []string) (simavionics.EventBus, error) {
	socket, err := createSocket(masterNode, addrs)
	if err != nil {
		return nil, err
	}

	b := &eventBus{
		masterNode:    masterNode,
		localPubChan:  make(chan *simavionics.Event),
		remotePubChan: make(chan *simavionics.Event),
		subs:          simavionics.NewEventSubscriptions(),
		socket:        socket,
	}
	go b.localPublisher()
	go b.remotePublisher()
	go b.receiver()
	return b, nil
}

func createSocket(masterNode bool, addrs []string) (mangos.Socket, error) {
	socket, err := bus.NewSocket()
	if err != nil {
		return nil, err
	}

	socket.AddTransport(ipc.NewTransport())
	socket.AddTransport(tcp.NewTransport())
	socket.SetOption(mangos.OptionReadQLen, 65535)
	socket.SetOption(mangos.OptionWriteQLen, 65535)

	for _, addr := range addrs {
		if masterNode {
			log.Printf("[bus.remote] Listing on %v", addr)
			if err = socket.SetOption(mangos.OptionRaw, true); err != nil {
				return nil, err
			}
			if err = socket.Listen(addr); err != nil {
				return nil, err
			}
		} else {
			log.Printf("[bus.remote] Dialing to %v", addr)
			if err = socket.Dial(addr); err != nil {
				return nil, err
			}
		}
	}

	return socket, nil
}
