package main

import (
	"github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
)

type Kaladont struct {
	socket interface{}
}

func initKaladont() *Kaladont {
	socketServer := gosocketio.NewServer(transport.GetDefaultWebsocketTransport())
	kt := Kaladont{socket: socketServer}
	return &kt
}
