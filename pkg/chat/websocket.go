package chat

import "github.com/gorilla/websocket"

type client struct {
	hub  *hub
	conn *websocket.Conn
	send chan []byte
}

type hub struct {
	register   chan *client
	unregister chan *client
	clients    map[*client]bool
	broadcast  chan []byte
}

func newHub() *hub {
	return &hub{
		register:   make(chan *client),
		unregister: make(chan *client),
		clients:    make(map[*client]bool),
		broadcast:  make(chan []byte),
	}
}

func newClient(hub *hub, conn *websocket.Conn) *client {
	return &client{
		hub:  hub,
		conn: conn,
		send: make(chan []byte, 1),
	}
}

func (h *hub) broadcastMessage(message []byte) {
	h.broadcast <- message
}

func (c *client) trySend(message []byte) (sent bool, closed bool) {
	defer func() {
		if recover() != nil {
			closed = true
			sent = false
		}
	}()

	select {
	case c.send <- message:
		return true, false
	default:
		return false, false
	}
}

func (h *hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				sent, closed := client.trySend(message)
				if sent {
					continue
				}

				delete(h.clients, client)

				if !closed {
					close(client.send)
				}
			}
		}
	}
}
