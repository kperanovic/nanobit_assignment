package main

import (
	"sync"

	"github.com/gorilla/websocket"
)

// clients is used to keep all client connections
type clients struct {
	sync.Mutex
	list []*websocket.Conn
}

// newClients will create a new clients instance
func newClients() *clients {
	return &clients{
		list: make([]*websocket.Conn, 0),
	}
}

// add adds a new client connection to the list
func (c *clients) add(conn *websocket.Conn) {
	c.Lock()
	defer c.Unlock()

	c.list = append(c.list, conn)
}

// remove will remove a client connection from the list
func (c *clients) remove(clientConn *websocket.Conn) {
	c.Lock()
	defer c.Unlock()

	for i, conn := range c.list {
		if conn == clientConn {
			c.list = append(c.list[:i], c.list[i+1:]...)
		}
	}
}

// broadcast sends a message to all connected clients.
// If err happens while sending messages to connected ws clients,
// it will not exit processing and exit, but it will append all errors in
// err slice and return it as a func return value.
// With this we're ensuring that even thought 1 client doesn't receive the message,
// func will attempt to send the message to all of them.
func (c *clients) broadcast(msg []byte) []error {
	c.Lock()
	defer c.Unlock()

	var errs []error
	for _, client := range c.list {
		if client != nil {
			if err := client.WriteMessage(websocket.TextMessage, msg); err != nil {
				errs = append(errs, err)
			}
		}
	}

	return errs
}
