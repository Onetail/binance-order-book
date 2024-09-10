package controller

import (
	"binance-dashboard/app"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Html struct {
	app    *app.App
	server *app.HTTPServer
	router *gin.RouterGroup
}

func (h *Html) Init(server *app.HTTPServer) {

	h.server = server
	h.app = server.App
	h.router = server.GetEngine().Group("")
	h.router.GET("", h.getHtml)
}

func (h *Html) getHtml(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"title": "Main website",
	})
}
