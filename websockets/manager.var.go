package websockets

var manager = ClientManager{
    broadcast:  make(chan []byte),
    register:   make(chan *Client),
    unregister: make(chan *Client),
    clients:    make(map[*Client]bool),
}

func ManagerStart(openmsg func()interface{}) {
	go manager.start(openmsg)
}