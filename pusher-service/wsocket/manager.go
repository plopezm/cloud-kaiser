package wsocket

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/plopezm/cloud-kaiser/core/logger"
	"github.com/satori/go.uuid"
	"net/http"
)

type ClientManager struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

var manager = ClientManager{
	broadcast:  make(chan []byte),
	register:   make(chan *Client),
	unregister: make(chan *Client),
	clients:    make(map[*Client]bool),
}

func Start() {
	var log = logger.GetLogger()
	for {
		select {
		case conn := <-manager.register:
			log.Debug(fmt.Sprintf("New socket connection: %s", conn.id))
			manager.clients[conn] = true
		case conn := <-manager.unregister:
			log.Debug(fmt.Sprintf("Socket disconnected: %s", conn.id))
			if _, ok := manager.clients[conn]; ok {
				close(conn.send)
				delete(manager.clients, conn)
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

func Broadcast(message interface{}, ignore *Client) {
	data, _ := json.Marshal(message)
	for conn := range manager.clients {
		if conn != ignore {
			conn.send <- data
		}
	}
}

func WsPage(res http.ResponseWriter, req *http.Request) {
	conn, error := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(res, req, nil)
	if error != nil {
		http.NotFound(res, req)
		return
	}
	client := &Client{id: uuid.NewV4().String(), socket: conn, send: make(chan []byte)}

	manager.register <- client

	go client.read()
	go client.write()
}
