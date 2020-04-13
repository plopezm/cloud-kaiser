package wsocket

import (
	"encoding/json"

	"github.com/gorilla/websocket"
	"github.com/plopezm/cloud-kaiser/core/logger"
)

//Client Represents a Web socket client
type Client struct {
	id     string
	socket *websocket.Conn
	send   chan []byte
}

//Message a web socket received message
type Message struct {
	Sender    string `json:"sender,omitempty"`
	Recipient string `json:"recipient,omitempty"`
	Content   string `json:"content,omitempty"`
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
	log := logger.GetLogger()

	defer func() {
		log.Debug("Closing socket connection: " + c.id)
		c.socket.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				log.Debug("Error sending message to socket connection: " + c.id)
				c.socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			log.Debug("Sending message to socket connection: " + c.id)
			c.socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}
