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
    Content   map[string]interface{} `json:"content,omitempty"`
}

func (manager *ClientManager) start(openmsg func()interface{}) {
    for {
        select {
        case conn := <-manager.register:
            manager.clients[conn] = true
            notebookMsg := Message{Type: "notebook"} //Esto habría que quitarlo luego
            notebookMsg.SetContent(openmsg())
            jsonMessage, _ := json.Marshal(notebookMsg)
            message := Message{}
            message.SetContent("/A socket has connected.")
            openMessage, _ := json.Marshal(&message)
            manager.sendTo(jsonMessage, conn)
            manager.send(openMessage, conn)
        case conn := <-manager.unregister:
            if _, ok := manager.clients[conn]; ok {
                close(conn.send)
                delete(manager.clients, conn)
                message := Message{}
                message.SetContent("/A socket has disconnected.")
                jsonMessage, _ := json.Marshal(&message)
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

func (message *Message) SetContent(content interface{}) {
    m := make(map[string]interface{})
    m["msg"] = content
    message.Content = m
}
