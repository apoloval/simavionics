package remote

import (
	"github.com/apoloval/simavionics"
	"github.com/go-mangos/mangos"
	"github.com/go-mangos/mangos/protocol/bus"
	"github.com/go-mangos/mangos/transport/tcp"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("bus.remote")

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
	log.Debug("Publishing event", ev.Name, "with value", ev.Value)
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
			log.Error("Failed to receive message from socket: ", err)
		}

		ev, err = simavionics.DecodeEvent(msg.Body)
		if err != nil {
			log.Error("Failed to decode message from socket: ", err)
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
		log.Error("Failed to publish on socket ", bus.socket)
	}
}

func NewMasterEventBus(localAddr string) (simavionics.EventBus, error) {
	return NewEventBus(true, localAddr)
}

func NewSlaveEventBus(remoteAddr string) (simavionics.EventBus, error) {
	return NewEventBus(false, remoteAddr)
}

func NewEventBus(masterNode bool, addr string) (simavionics.EventBus, error) {
	socket, err := createSocket(masterNode, addr)
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

func createSocket(masterNode bool, addr string) (mangos.Socket, error) {
	socket, err := bus.NewSocket()
	if err != nil {
		return nil, err
	}

	socket.AddTransport(tcp.NewTransport())
	socket.SetOption(mangos.OptionReadQLen, 65535)
	socket.SetOption(mangos.OptionWriteQLen, 65535)

	if masterNode {
		log.Notice("Listening on", addr)
		if err = socket.SetOption(mangos.OptionRaw, true); err != nil {
			return nil, err
		}
		if err = socket.Listen(addr); err != nil {
			return nil, err
		}
	} else {
		log.Notice("Dialing to", addr)
		if err = socket.Dial(addr); err != nil {
			return nil, err
		}
	}

	return socket, nil
}
