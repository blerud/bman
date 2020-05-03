package server

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	tick          = 30
	millisPerTick = time.Duration(time.Second / tick)
	timeToConnect = time.Duration(10 * time.Second)
)

type EntitiesView interface {
	collisions(entity *Entity, newX float32, newY float32) []*Entity
	create(entity *Entity) bool
	update(entity *Entity) bool
	delete(entity *Entity) bool
	genEntityId() int32
}

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

	created []int32
	updated []int32
	deleted []*Entity
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
		make([]int32, 0),
		make([]*Entity, 0),
	}
}

func (server *Server) run() {
	server.init()
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
			for _, client := range server.clients {
				client.writeQueue <- message.encodeMessage()
			}
		case client := <-server.registerQueue:
			fmt.Println("adding client")
			server.clients[client.id] = client
			server.sendState(client.id)

			player := newPlayer(client.id, 2, 2)
			server.entities[client.id] = player
			server.clientIdToEntityId[client.id] = client.id
			server.entityIdToClientId[client.id] = client.id
			server.createdEntity(player)
		case client := <-server.unregisterQueue:
			if _, ok := server.clients[client.id]; ok {
				fmt.Println("removing client")
				_ = client.conn.Close()
				delete(server.clients, client.id)
				close(client.writeQueue)
			}
		case _ = <-ticker.C:
			server.heartbeat()
			server.step()
		}
	}
}

func (server *Server) init() {
	for i := 0; i < 15; i++ {
		topWall := newHardWall(server.genEntityId(), float32(i), 0)
		botWall := newHardWall(server.genEntityId(), float32(i), 12)
		server.create(topWall)
		server.create(botWall)
	}

	for i := 1; i < 12; i++ {
		leftWall := newHardWall(server.genEntityId(), 0, float32(i))
		rightWall := newHardWall(server.genEntityId(), 14, float32(i))
		server.create(leftWall)
		server.create(rightWall)
	}

	for i := 2; i < 12; i += 2 {
		for j := 2; j < 13; j += 2 {
			wall := newHardWall(server.genEntityId(), float32(j), float32(i))
			server.create(wall)
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

func (server *Server) sendToClient(clientId int32, message Message) {
	message.timestamp = time.Now().UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
	message.length = len(message.content) + 8
	server.clients[clientId].writeQueue <- message.encodeMessage()
}

func (server *Server) sendState(clientId int32) {
	createBuf := make([]byte, 1)
	createBuf[0] = byte(len(server.entities))

	for _, entity := range server.entities {
		entityBytes := entity.encode()
		createBuf = append(createBuf, entityBytes...)
	}

	createMessage := Message{messageCreated, 0, 0, createBuf}
	server.sendToClient(clientId, createMessage)
}

func (server *Server) createdEntity(entity *Entity) {
	createBuf := make([]byte, 1)
	createBuf[0] = byte(1)
	entityBytes := server.entities[entity.entityId].encode()
	createBuf = append(createBuf, entityBytes...)
	fmt.Printf("created object %d\n", entity.entityId)

	createMessage := Message{messageCreated, 0, 0, createBuf}
	server.sendQueue <- createMessage
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

func (server *Server) step() {
	for entityId, entity := range server.entities {
		updated := entity.entityInfo.step(server)
		if updated {
			server.updated = append(server.updated, entityId)
		}
	}

	if len(server.created) > 0 {
		createBuf := make([]byte, 1)
		createBuf[0] = byte(len(server.created))
		for _, entityId := range server.created {
			entityBytes := server.entities[entityId].encode()
			createBuf = append(createBuf, entityBytes...)
		}

		server.created = make([]int32, 0)

		createMessage := Message{messageCreated, 0, 0, createBuf}
		server.sendQueue <- createMessage
	}

	if len(server.updated) > 0 {
		updateBuf := make([]byte, 1)
		updateBuf[0] = byte(len(server.updated))
		for _, entityId := range server.updated {
			entityBytes := server.entities[entityId].encode()
			updateBuf = append(updateBuf, entityBytes...)
		}

		server.updated = make([]int32, 0)

		updateMessage := Message{messageUpdated, 0, 0, updateBuf}
		server.sendQueue <- updateMessage
	}

	if len(server.deleted) > 0 {
		deleteBuf := make([]byte, 1)
		deleteBuf[0] = byte(len(server.deleted))
		for _, entity := range server.deleted {
			entityBytes := entity.entityInfo.encode()
			deleteBuf = append(deleteBuf, entityBytes...)
		}

		server.deleted = make([]*Entity, 0)

		deleteMessage := Message{messageDeleted, 0, 0, deleteBuf}
		server.sendQueue <- deleteMessage
	}
}

func (server *Server) collisions(entity *Entity, newX float32, newY float32) []*Entity {
	collisions := make([]*Entity, 0)

	for _, e := range server.entities {
		if entity.entityId == e.entityId {
			continue
		}

		collidingNow := server.collides(entity.x, entity.y, entity.width, entity.height, e.x, e.y, e.width, e.height)
		if server.collides(newX, newY, entity.width, entity.height, e.x, e.y, e.width, e.height) && !collidingNow {
			collisions = append(collisions, e)
		}
	}

	return collisions
}

func (server *Server) create(entity *Entity) bool {
	if _, ok := server.entities[entity.entityId]; ok {
		return false
	}

	server.entities[entity.entityId] = entity
	server.created = append(server.created, entity.entityId)
	return true
}

func (server *Server) update(entity *Entity) bool {
	return true
}

func (server *Server) delete(entity *Entity) bool {
	if _, ok := server.entities[entity.entityId]; !ok {
		return false
	}

	delete(server.entities, entity.entityId)
	server.deleted = append(server.deleted, entity)
	return true
}

func (server *Server) collides(x1 float32, y1 float32, w1 float32, h1 float32,
	x2 float32, y2 float32, w2 float32, h2 float32) bool {
	if x1 > x2+w2 || x1+w1 < x2 || y1 > y2+h2 || y1+h1 < y2 {
		return false
	}
	return true
}

func (server *Server) genEntityId() int32 {
	for {
		id := rand.Int31()
		if _, ok := server.entities[id]; !ok {
			return id
		}
	}
}

func messageFromBytes(bytes []byte) Message {
	messageType := bytes[0]
	length := readInt32FromBuffer(bytes[1:])
	timestamp := readInt64FromBuffer(bytes[5:])
	message := bytes[13:]
	return Message{messageType, int(length), timestamp, message}
}
