package websockets

import (
	"encoding/json"
	)

type ClientManager struct {
    clients    map[*Client]bool
    broadcast  chan []byte
    register   chan *Client
    unregister chan *Client
}

type Message struct {
    Sender    string `json:"sender,omitempty"`
    Recipient string `json:"recipient,omitempty"`
    Type     string `json:"type,omitempty"`
    Content   interface{} `json:"content,omitempty"`
}

func (manager *ClientManager) start(openmsg interface{}) {
    for {
        select {
        case conn := <-manager.register:
            manager.clients[conn] = true
            jsonMessage, _ := json.Marshal(openmsg)
            openMessage, _ := json.Marshal(&Message{Content: "/A socket has disconnected."})
            manager.sendTo(jsonMessage, conn)
            manager.send(openMessage, conn)
        case conn := <-manager.unregister:
            if _, ok := manager.clients[conn]; ok {
                close(conn.send)
                delete(manager.clients, conn)
                jsonMessage, _ := json.Marshal(&Message{Content: "/A socket has disconnected."})
                manager.send(jsonMessage, conn)
            }
        case message := <-manager.broadcast:
            for conn := range manager.clients {
                select {
                case conn.send <- message:
                default:
                    close(conn.send)
                    delete(manager.clients, conn)
                }
            }
        }
    }
}

func (manager *ClientManager) send(message []byte, ignore *Client) {
    for conn := range manager.clients {
        if conn != ignore {
            conn.send <- message
        }
    }
}
func (manager *ClientManager) sendTo(message []byte, client *Client) {
    client.send <- message
}
