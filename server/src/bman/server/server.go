package server

import (
	"fmt"
	"time"
)

const (
	tick          = 30
	millisPerTick = time.Duration(time.Second / tick)
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

	entities           map[int32]*Entity
	clientIdToEntityId map[int32]int32
	entityIdToClientId map[int32]int32

	updated []int32
}

func newServer(serverId int32, coord *Coordinator) *Server {
	return &Server{
		serverId,
		coord,
		make(map[int32]*Client),
		make(chan []byte, 100),
		make(chan Message, 100),
		make(chan *Client),
		make(chan *Client),
		make(map[int32]*Entity),
		make(map[int32]int32),
		make(map[int32]int32),
		make([]int32, 0),
	}
}

func (server *Server) run() {
	ticker := time.NewTicker(millisPerTick)
	for {
		select {
		case messageBytes := <-server.readQueue:
			message := messageFromBytes(messageBytes)
			server.process(message)
			//code := message.messageType
			//length := message.length
			//timestamp := message.timestamp
			//fmt.Printf("===== code: %d, length: %d, timestamp: %d\n", code, length, timestamp)
		case message := <-server.sendQueue:
			message.timestamp = time.Now().UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
			message.length = len(message.content) + 8
			//fmt.Println("sent message ", message.encodeMessage())
			for _, client := range server.clients {
				client.writeQueue <- message.encodeMessage()
			}
		case client := <-server.registerQueue:
			fmt.Println("adding client")
			server.clients[client.id] = client
			player := newPlayer(client.id, 0, 0)
			server.entities[client.id] = player
			server.clientIdToEntityId[client.id] = client.id
			server.entityIdToClientId[client.id] = client.id
		case client := <-server.unregisterQueue:
			if _, ok := server.clients[client.id]; ok {
				fmt.Println("removing client")
				_ = client.conn.Close()
				delete(server.clients, client.id)
				close(client.writeQueue)
			}
		case _ = <-ticker.C:
			server.heartbeat()
			server.tick()
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
	entityId := server.clientIdToEntityId[clientId]
	playerEntity := server.entities[entityId]
	if player, ok := playerEntity.entityInfo.(*Player); ok {
		player.processPlayerAction(message)
		return true
	}
	return false
}

func (server *Server) heartbeat() {
	server.sendQueue <- Message{messageServerHeartbeat, 0, 0, []byte{}}
}

func (server *Server) tick() {
	for _, entityId := range server.clientIdToEntityId {
		entity := server.entities[entityId]
		if player, ok := entity.entityInfo.(*Player); ok {
			xSpeed := playerXSpeed / tick
			ySpeed := playerYSpeed / tick

			xMove := float32(0)
			yMove := float32(0)
			action := player.action
			if action.up {
				yMove -= ySpeed
			}
			if action.down {
				yMove += ySpeed
			}
			if action.left {
				xMove -= xSpeed
			}
			if action.right {
				xMove += xSpeed
			}

			player.entity.x += xMove
			player.entity.y += yMove

			if xMove != 0 || yMove != 0 {
				server.updated = append(server.updated, entityId)
			}
		}
	}

	if len(server.updated) > 0 {
		updateBuf := make([]byte, 1)
		updateBuf[0] = byte(len(server.updated))
		fmt.Println("updated: ", server.updated)
		for _, entityId := range server.updated {
			entityBytes := server.entities[entityId].encode()
			updateBuf = append(updateBuf, entityBytes...)
		}

		server.updated = make([]int32, 0)

		updateMessage := Message{messageUpdated, 0, 0, updateBuf}
		server.sendQueue <- updateMessage
	}
}

func messageFromBytes(bytes []byte) Message {
	messageType := bytes[0]
	length := readInt32FromBuffer(bytes[1:])
	timestamp := readInt64FromBuffer(bytes[5:])
	message := bytes[13:]
	return Message{messageType, int(length), timestamp, message}
}
