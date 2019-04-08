package websockets

import (
    "log"
    "time"
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/satori/go.uuid"
	)

type Client struct {
	id     string
	socket *websocket.Conn
	send   chan []byte
}

func (c *Client) read(callback func(map[string]interface{})) {
    defer func() {
        manager.unregister <- c
        c.socket.Close()
    }()

    for {
        _, message, err := c.socket.ReadMessage()
        if err != nil {
            log.Printf("la conexión ha sido cerrada %v", err)
            manager.unregister <- c
            c.socket.Close()
            break
        }
        log.Printf("mensaje", string(message))
        compiledMessage := ProcessMessage(message)
        jsonMessage, _ := json.Marshal(&Message{Sender: c.id,
             Type: compiledMessage.Type, 
             Content: compiledMessage.Content,
            })
        if compiledMessage.Type == "edited" {
            callback(compiledMessage.Content)
        }
        //Aquí faltan los efectos colaterales, que en este caso sería la edición de los archivos, habría que hacer
        //un unmarshal, ver el tipo de evento que se mandó y actuar en consecuencia
        
        if compiledMessage.Type != "ping" {
            manager.broadcast <- jsonMessage
        }
        
    }
}

func (c *Client) write() {
    defer func() {
        // log.Fatalf("la conexión ha sido cerrada")
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

func ClientStart(conn *websocket.Conn, readCallback func(map[string]interface{})) {
    id, _ := uuid.NewV4()
    conn.SetReadDeadline(time.Time{})
    conn.SetWriteDeadline(time.Time{})
	client := &Client{id: id.String(), socket: conn, send: make(chan []byte)}

    manager.register <- client

    go client.read(readCallback)
    go client.write()
}

func ProcessMessage(data []byte) Message{
    message := Message{}
    err := json.Unmarshal(data, &message)
    if err != nil {
        log.Fatalf("Este no es el mensaje!! %s", string(data))
    }
    return message
}