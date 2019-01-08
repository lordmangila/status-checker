package rest

import (
	"bytes"
	"log"
	"time"

	"github.com/gorilla/websocket"

	"bitbucket.org/lordmangila/status-checker/pkg/checker"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

// Client maintains a mapping to the server, sites and websocket connection.
type Client struct {
	server *Server

	// Registered sites lists.
	sites map[*checker.Site]bool

	// Websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

// Listen processes incoming client messages.
//
// Listen should always be started as a goroutine for each websocket client
// which handles all incoming client requests.
func (c *Client) Listen() {
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	defer func() {
		delete(c.server.clients, c)
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))

		site := &checker.Site{
			URL: string(message),
		}
		site.Validate()
		if "" != site.Error {
			c.send <- site.Marshal()
		} else {
			c.sites[site] = true
		}
	}
}

// Broadcast processes server messages to the websocket client connection.
//
// Broadcast should always be started as a goroutine for each websocket client
// which handles all the information sending to the client.
func (c *Client) Broadcast() {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.sendMessage(message, ok)
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); nil != err {
				return
			}
		}
	}
}

// CheckSites iterates over the client sites to perform health check.
func (c *Client) CheckSites() {
	for site := range c.sites {
		site.HealthCheck()
		c.send <- site.Marshal()
	}
}

// sendMessage sends the message to the websocket client connection.
func (c *Client) sendMessage(message []byte, ok bool) {
	c.conn.SetWriteDeadline(time.Now().Add(writeWait))
	if !ok {
		// The service closed the channel.
		c.conn.WriteMessage(websocket.CloseMessage, []byte{})
		return
	}

	w, err := c.conn.NextWriter(websocket.TextMessage)
	if err != nil {
		return
	}

	w.Write(message)

	if err := w.Close(); nil != err {
		return
	}
}
