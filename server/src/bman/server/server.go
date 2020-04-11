package server

import (
	"fmt"
	"time"
)

const (
	tick = 30
	millisPerTick = time.Duration(time.Millisecond / tick)
	timeToConnect = time.Duration(1 * time.Second)
)

type Server struct {
	serverId int32
	coordinator *Coordinator
	clients  map[int32]*Client

	readQueue chan []byte
	sendQueue chan Message
	registerQueue chan *Client
	unregisterQueue chan *Client

	entities map[int32]Entity

}

func newServer(serverId int32, coord *Coordinator) *Server {
	return &Server{
		serverId,
		coord,
		make(map[int32]*Client),
		make(chan []byte),
		make(chan Message),
		make(chan *Client),
		make(chan *Client),
		make(map[int32]Entity),
	}
}

func (server *Server) run() {
	ticker := time.NewTicker(millisPerTick)
	for {
		select {
		case messageBytes := <-server.readQueue:
			message := messageFromBytes(messageBytes)
			server.process(message)
			for _, client := range server.clients {
				client.writeQueue <- messageBytes
			}
			fmt.Printf("===== %s\n", string(messageBytes))
		case message := <-server.sendQueue:
			message.timestamp = time.Now().UnixNano()
			message.length = len(message.content)
			for _, client := range server.clients {
				client.writeQueue <- message.encodeMessage()
			}
		case client := <-server.registerQueue:
			fmt.Println("adding client")
			server.clients[client.id] = client
		case client := <-server.unregisterQueue:
			if _, ok := server.clients[client.id]; ok {
				fmt.Println("removing client")
				_ = client.conn.Close()
				delete(server.clients, client.id)
				close(client.writeQueue)
			}
		case _ = <-ticker.C:
			server.sendTick()
		}
	}
}

func (server *Server) startDeleteTimer() {
	timer := time.NewTimer(timeToConnect)
	<-timer.C
	if len(server.clients) == 0 {
		server.coordinator.closeServer(server.serverId)
	}
}

func (server *Server) process(message Message) bool {
	processedAny := false
	switch message.messageType {
	case messageClientHeartbeat:
		return true
	case messageUpdated:
		content := message.content
		numberUpdated := int(content[0])
		for i := 0; i < numberUpdated; i++ {
			processed, bytesRead := server.processUpdate(content)
			processedAny = processedAny || processed
			content = content[bytesRead:]
		}
		return true
	default:
		return false
	}
}

func (server *Server) processUpdate(content []byte) (bool, int) {
	entityType := content[0]
	entityId := readInt32FromBuffer(content[1:])
	switch entityType {
	case PLAYER:
		player := server.entities[entityId]
		bytesRead := player.entityInfo.processUpdate(content)
		return true, bytesRead
	default:
		return false, 0
	}
}

func (server *Server) sendTick() {
	for _, client := range server.clients {
		client.writeQueue <- []byte{}
	}
}

func messageFromBytes(bytes []byte) Message {
	messageType := bytes[0]
	length := readInt32FromBuffer(bytes[1:])
	timestamp := readInt64FromBuffer(bytes[5:])
	message := bytes[13:]
	return Message{messageType, int(length), timestamp, message}
}
