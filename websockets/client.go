package websockets

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/satori/go.uuid"
	)

type Client struct {
	id     string
	socket *websocket.Conn
	send   chan []byte
}

func (c *Client) read() {
    defer func() {
        manager.unregister <- c
        c.socket.Close()
    }()

    for {
        _, message, err := c.socket.ReadMessage()
        if err != nil {
            manager.unregister <- c
            c.socket.Close()
            break
        }
        jsonMessage, _ := json.Marshal(&Message{Sender: c.id, Content: string(message)})
        manager.broadcast <- jsonMessage
    }
}

func (c *Client) write() {
    defer func() {
        c.socket.Close()
    }()

    for {
        select {
        case message, ok := <-c.send:
            if !ok {
                c.socket.WriteMessage(websocket.CloseMessage, []byte{})
                return
            }

            c.socket.WriteMessage(websocket.TextMessage, message)
        }
    }
}

func ClientStart(conn *websocket.Conn) {
	id, _ := uuid.NewV4()
	client := &Client{id: id.String(), socket: conn, send: make(chan []byte)}

    manager.register <- client

    go client.read()
    go client.write()
}