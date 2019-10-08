// Based on https://github.com/mailhog/MailHog-Server/blob/master/websockets/connection.go
package ws

import (
	"time"

	"github.com/gorilla/websocket"
	"github.com/ns3777k/mailcage/storage"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan *storage.Event
}

func NewClient(hub *Hub, conn *websocket.Conn) *Client {
	client := &Client{hub: hub, conn: conn, send: make(chan *storage.Event, 256)}
	hub.register <- client

	return client
}

func (c *Client) ReadPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		if _, _, err := c.conn.NextReader(); err != nil {
			break
		}
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.writeControl(websocket.CloseMessage)
				return
			}

			if err := c.writeJSON(message); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.writeControl(websocket.PingMessage); err != nil {
				return
			}
		}
	}
}

func (c *Client) writeJSON(message interface{}) error {
	c.conn.SetWriteDeadline(time.Now().Add(writeWait))
	return c.conn.WriteJSON(message)
}

func (c *Client) writeControl(messageType int) error {
	c.conn.SetWriteDeadline(time.Now().Add(writeWait))
	return c.conn.WriteMessage(messageType, []byte{})
}
