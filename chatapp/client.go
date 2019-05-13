package main

import (
	"fmt"

	"github.com/gorilla/websocket"
)

type client struct {
	socket *websocket.Conn
	send   chan []byte
	room   *room
}

/*
Read is WebSocket Message Receiver.
This method detect WebSocket message and send forward all of clients that is join to chat room.
*/
func (c *client) Read() {
	defer c.socket.Close()
	for {
		_, msg, err := c.socket.ReadMessage()
		fmt.Println("receave message")
		if err != nil {
			fmt.Println("error!!! websocket can not receave message.", err)
			break
		}
		fmt.Println("receave message")
		c.room.forward <- msg
	}
}

/*
Write is WebSocket Message Sender
*/
func (c *client) Write() {
	defer c.socket.Close()
	for msg := range c.send {
		if err := c.socket.WriteMessage(websocket.TextMessage, msg); err != nil {
			break
		}
	}
}
