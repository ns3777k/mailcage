package ws

import (
	"github.com/ns3777k/mailcage/storage"
)

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan *storage.Event
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	hub := &Hub{
		broadcast:  make(chan *storage.Event),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}

	go hub.run()

	return hub
}

func (h *Hub) Broadcast(event *storage.Event) {
	h.broadcast <- event
}

func (h *Hub) run() {
	unregister := func(client *Client) {
		if _, ok := h.clients[client]; ok {
			delete(h.clients, client)
			close(client.send)
		}
	}

	for {
		select {
		case client := <-h.register:
			h.clients[client] = true

		case client := <-h.unregister:
			unregister(client)

		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					unregister(client)
				}
			}
		}
	}
}
