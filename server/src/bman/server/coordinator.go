package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"math/rand"
	"net/http"
	"strconv"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Coordinator struct {
	clients map[int32]*Client
	servers map[int32]*Server
}

func NewCoordinator() *Coordinator {
	return &Coordinator{
		make(map[int32]*Client),
		make(map[int32]*Server),
	}
}

func (coord *Coordinator) WsHandler(ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		fmt.Println("error")
		return
	}

	params := ctx.Request.URL.Query()

	idS, hasId := params["userId"]
	name, hasName := params["name"]
	gameIdS, hasName := params["gameId"]
	id, err := strconv.Atoi(idS[0])
	gameId, errGame := strconv.Atoi(gameIdS[0])
	if !hasId || !hasName || err != nil || errGame != nil {
		_ = conn.Close()
	}
	fmt.Printf("afafaa %d %d %s\n", id, gameId, name)

	coord.addClient(conn, int32(id), name[0], int32(gameId))
}

func (coord *Coordinator) CreateGameHandler(ctx *gin.Context) {
	serverId := coord.createServer()
	ctx.JSON(http.StatusOK, gin.H{
		"gameId": serverId,
	})
}

func (coord *Coordinator) addClient(conn *websocket.Conn, id int32, name string, gameId int32) {
	// todo need to check if server exists
	server := coord.getServer(gameId)
	client := &Client{
		conn,
		server,
		id,
		name,
		make([]byte, 0),
		make(chan []byte),
	}
	coord.clients[id] = client
	client.register()
	go client.readMessages()
	go client.writeMessages()
}

func (coord *Coordinator) createServer() int32 {
	serverId := coord.getUniqServerId()
	server := newServer(serverId, coord)
	coord.servers[serverId] = server
	go server.run()
	return serverId
}

func (coord *Coordinator) getServer(serverId int32) *Server {
	return coord.servers[serverId]
}

func (coord *Coordinator) getUniqServerId() int32 {
	for {
		id := rand.Int31()
		if _, ok := coord.servers[id]; !ok {
			return id
		}
	}
}

func (coord *Coordinator) closeServer(serverId int32) {
	fmt.Printf("closing server %d\n", serverId)
	for clientId, client := range coord.servers[serverId].clients {
		client.close()
		delete(coord.clients, clientId)
	}
	delete(coord.servers, serverId)
}
