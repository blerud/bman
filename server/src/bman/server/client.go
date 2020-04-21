package server

import (
	"encoding/binary"
	"github.com/gorilla/websocket"
	"log"
)

type Client struct {
	conn               *websocket.Conn
	server             *Server
	id                 int32
	name               string
	incompleteMessages []byte

	writeQueue chan []byte
}

func (c *Client) readMessages() {
	defer func() {
		c.server.unregisterQueue <- c
		_ = c.conn.Close()
	}()
	for {
		_, message, err := c.conn.ReadMessage()
		// todo handle ping
		if err != nil {
			log.Printf("websocket closing, %v", err)
			return
		}
		c.incompleteMessages = append(c.incompleteMessages, message...)
		for _, msg := range c.getCompleteMessages() {
			c.server.readQueue <- msg
		}
	}
}

func (c *Client) writeMessages() {
	defer func() {
		c.server.unregisterQueue <- c
		_ = c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.writeQueue:
			if !ok {
				// server closed channel
				return
			}
			// todo we have a simple write here, need to make it more robust
			err := c.conn.WriteMessage(websocket.BinaryMessage, message)
			if err != nil {
				return
			}
		}
	}
}

func (c *Client) getCompleteMessages() [][]byte {
	if len(c.incompleteMessages) < 5 {
		return make([][]byte, 0)
	}
	//fmt.Println(c.incompleteMessages)
	completeMessages := make([][]byte, 0)
	i := 0
	for i+5 < len(c.incompleteMessages) {
		length := int(binary.BigEndian.Uint32(c.incompleteMessages[i+1 : i+5]))
		if len(c.incompleteMessages)-5 >= length {
			completeMessages = append(completeMessages, c.incompleteMessages[i:i+5+length])
			c.incompleteMessages = c.incompleteMessages[i+5+length:]
			i += 5 + length
		} else {
			break
		}
	}
	return completeMessages
}

func (c *Client) register() {
	c.server.registerQueue <- c
}

func (c *Client) close() {
	c.server.unregisterQueue <- c
}
