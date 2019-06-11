package main

import (
	"time"

	"github.com/gorilla/websocket"
)

type client struct {
	socket    *websocket.Conn
	send      chan *message // channel on which messages are sent
	room      *room
	avatarURL string
	userData  map[string]interface{}
}

func (c *client) read() {
	defer c.socket.Close()
	for {
		var msg message
		if err := c.socket.ReadJSON(&msg); err != nil {
			return
		}
		msg.When = time.Now()
		msg.Name = c.userData["name"].(string)
		msg.AvatarURL = c.avatarURL
		c.room.forward <- &msg
	}
}

func (c *client) write() {
	defer c.socket.Close()
	for msg := range c.send {
		err := c.socket.WriteJSON(msg)
		if err != nil {
			return
		}
	}
}
