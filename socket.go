package main

import (
	"encoding/json"
	"log"

	"github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
)

type Command struct {
	Command string            `json:"command"`
	Payload map[string]string `json:"payload"`
}

type Message struct {
	Payload interface{}
}

func emitError(c *gosocketio.Channel, data interface{}) {
	c.Emit("error", map[string]interface{}{
		"error": data,
	})
}

func createSocketServer() *gosocketio.Server {
	s := gosocketio.NewServer(transport.GetDefaultWebsocketTransport())
	s.On("connection", func(c *gosocketio.Channel) {
		log.Println("New client connected")
	})

	s.On("subscribe", onSubscribe)
	s.On("nam", func(c *gosocketio.Channel) {
		c.Emit("NaM", map[string]string{})
	})

	return s
}

func onSubscribe(c *gosocketio.Channel, cmd string) {
	var command Command
	err := json.Unmarshal([]byte(cmd), &command)

	if err != nil {
		emitError(c, err.Error())
		return
	}

	room := command.Payload["room"]

	if room == "" {
		emitError(c, "Invalid property: room")
		return
	}

	c.Join(room)
	c.Emit("room:join", map[string]string{"room": room})
}
