package main

import (
	"fmt"
	"net/http"

	"log"

	"github.com/gorilla/websocket"
)

type room struct {
	forward chan []byte
	clients map[*client]bool
	join    chan *client
	leave   chan *client
}

// NewRoom is generator of Room
func NewRoom() *room {
	return &room{
		forward: make(chan []byte),
		clients: make(map[*client]bool),
		join:    make(chan *client),
		leave:   make(chan *client),
	}
}

func (r *room) Run() {
	for {
		select {
		case client := <-r.join:
			r.clients[client] = true
			fmt.Println("join clinet")
		case client := <-r.leave:
			delete(r.clients, client)
			close(client.send)
			fmt.Println("leave client")
		case msg := <-r.forward:
			for client := range r.clients {
				select {
				case client.send <- msg:
					//send message
					fmt.Println("send message for all of room client")
				default:
					// failed to send message
					delete(r.clients, client)
					close(client.send)
					fmt.Println("err!!! close client send channel")
				}
			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  socketBufferSize,
	WriteBufferSize: socketBufferSize,
}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal(err)
		return
	}

	client := &client{
		socket: socket,
		send:   make(chan []byte, messageBufferSize),
		room:   r,
	}

	r.join <- client
	defer func() { r.leave <- client }()
	go client.Write()
	client.Read()
}
