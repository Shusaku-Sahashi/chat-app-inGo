# Chat app

this app is sample of WebSocket Chat Application

## Explains

### Models

- Client

  - Chat room client.
  - client is able to `join` and `leave` from room.

- Room

  - Chat room
  - Chat room is able to contain N Clients.

### Flow

`Client_A`
| Send `Message` By WebSocket For Chat Room
v
`socket.ReadMessage()` in `Client_A` // this is used for read client message from Front-Side
|
v
`room.Forward()` // send message for Room Members
|
V
All of room clients call `socket.WriteMessage(websocket.TextMessage, msg)` in `write()`
|
v
Message is presents in all of client's screen.
