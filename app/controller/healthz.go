package controller

import (
	"binance-dashboard/app"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HeartBeat struct {
	app    *app.App
	server *app.HTTPServer
	router *gin.RouterGroup
}

func (hb *HeartBeat) Init(server *app.HTTPServer) {

	hb.server = server
	hb.app = server.App

	hb.router = server.GetEngine().Group("/api/v1")

	hb.router.GET("/healthz", hb.getHeartBeat)

}

func (hb *HeartBeat) getHeartBeat(c *gin.Context) {
	c.JSON(http.StatusOK, 1)
}
