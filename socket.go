package main

import (
	"log"
	"strings"

	"github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
)

var topics = map[string]bool{"room": true}

func emitError(c *gosocketio.Channel, data interface{}) {
	c.Emit("error", map[string]interface{}{
		"error": data,
	})
}

func createSocketServer() *gosocketio.Server {
	s := gosocketio.NewServer(transport.GetDefaultWebsocketTransport())
	s.On(gosocketio.OnConnection, func(c *gosocketio.Channel) {
		log.Println("New websocket connection")
		c.Emit("connected", "")
		s.On("subscribe", HandleSubscription)
		s.On("unsubscribe", HandleUnsubscribe)
		s.On(gosocketio.OnDisconnection, HandleDisconnect)
	})

	return s
}

func HandleDisconnect(c *gosocketio.Channel) {
	log.Println("Client disconnected")
}

// HandleSubscription confirms that a user can subscribe
// to whatever they want to subscribe to
func HandleSubscription(c *gosocketio.Channel, s map[string]interface{}) {
	r := s["room"].(string)
	var split = strings.Split(r, ":")

	if len(split) < 2 {
		emitError(c, "Invalid room specified")
		return
	}

	topic, room := split[0], split[1]

	if topics[topic] == false {
		emitError(c, "Invalid topic specified")
		return
	}

	if room == "" {
		emitError(c, "Invalid room specified")
		return
	}

	c.Join(topic + ":" + room)
	c.Emit("room:join", map[string]string{"room": room})
}

func HandleUnsubscribe(c *gosocketio.Channel, s map[string]interface{}) {
	room := s["room"].(string)
	c.Leave(room)
}
