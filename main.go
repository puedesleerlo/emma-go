package main

import (
	"fmt"
	"net/http"
	"github.com/gorilla/websocket"
    ws"github.com/puedesleerlo/emma-go/websockets"
    "github.com/puedesleerlo/emma-go/programs"
	)

func main() {
    fmt.Println("Starting application...")
    notebookMsg := ws.Message{Type: "notebook", Content: programs.GetPrograms()}
    ws.ManagerStart(notebookMsg)
    http.HandleFunc("/ws", wsPage)
    http.ListenAndServe(":8886", nil)
}

func wsPage(res http.ResponseWriter, req *http.Request) {
    conn, error := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(res, req, nil)
    if error != nil {
        http.NotFound(res, req)
        return
	}
	ws.ClientStart(conn, programs.SaveEdited)
   
}