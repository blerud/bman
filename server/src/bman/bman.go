package main

import (
	"github.com/gin-gonic/gin"
	"net/http"

	"bman/server"
)

func homeHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func main() {
	r := gin.Default()
	coord := server.NewCoordinator()

	r.Static("/assets", "./client/dist")
	r.LoadHTMLGlob("client/*.html")
	r.GET("/ws", coord.WsHandler)
	r.GET("/", homeHandler)
	r.POST("/creategame", coord.CreateGameHandler)
	_ = r.Run(":6969")
}
