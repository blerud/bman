package server

import (
	"fmt"
	"time"
)

const (
	tick          = 30
	millisPerTick = time.Duration(time.Millisecond / tick)
	timeToConnect = time.Duration(1 * time.Second)
)

type Server struct {
	serverId    int32
	coordinator *Coordinator
	clients     map[int32]*Client

	readQueue       chan []byte
	sendQueue       chan Message
	registerQueue   chan *Client
	unregisterQueue chan *Client

	entities         map[int32]Entity
	clientIdToEntity map[int32]int32
	entityIdToClient map[int32]int32
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
		make(map[int32]int32),
		make(map[int32]int32),
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
			code := message.messageType
			length := message.length
			timestamp := message.timestamp
			fmt.Printf("===== code: %d, length: %d, timestamp: %d\n", code, length, timestamp)
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
	switch message.messageType {
	case messageClientHeartbeat:
		return true
	case messageClientAction:
		server.processClientAction(message)
		return true
	default:
		return false
	}
}

func (server *Server) processClientAction(message Message) bool {
	clientId := readInt32FromBuffer(message.content)
	entityId := server.clientIdToEntity[clientId]
	playerEntity := server.entities[entityId]
	if player, ok := playerEntity.entityInfo.(*Player); ok {
		player.processPlayerAction(message)
		return true
	}
	return false
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
