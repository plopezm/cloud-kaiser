package pusher

import "github.com/gorilla/websocket"

type Client struct {
	id       uint64
	socket   *websocket.Conn
	outbound chan []byte
}

func newClient(socket *websocket.Conn) *Client {
	return &Client{
		socket:   socket,
		outbound: make(chan []byte),
	}
}

func (client *Client) listen() {
	for {
		select {
		case data, ok := <-client.outbound:
			if !ok {
				client.socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			client.socket.WriteMessage(websocket.TextMessage, data)
		}
	}
}

func (client Client) close() {
	client.socket.Close()
	close(client.outbound)
}
