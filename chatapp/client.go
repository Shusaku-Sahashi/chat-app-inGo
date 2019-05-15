package main

import (
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

type client struct {
	socket   *websocket.Conn
	send     chan *message
	room     *room
	userData map[string]interface{}
}

/*
Read is WebSocket Message Receiver.
This method detect WebSocket message and send forward all of clients that is join to chat room.
*/
func (c *client) Read() {
	defer c.socket.Close()
	for {
		var msg = new(message)
		err := c.socket.ReadJSON(msg)
		if err != nil {
			fmt.Println("error!!! websocket can not receave message.", err)
			break
		}
		msg.When = time.Now()
		msg.Name = c.userData["name"].(string)
		c.room.forward <- msg
	}
}

/*
Write is WebSocket Message Sender
*/
func (c *client) Write() {
	defer c.socket.Close()
	for msg := range c.send {
		if err := c.socket.WriteJSON(msg); err != nil {
			break
		}
	}
}
