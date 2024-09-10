package controller

import (
	"binance-dashboard/app"
	"binance-dashboard/app/service"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Binance struct {
	app            *app.App
	server         *app.HTTPServer
	router         *gin.RouterGroup
	binanceService service.Binance
}

func (b *Binance) Init(server *app.HTTPServer) {

	b.server = server
	b.app = server.App

	b.binanceService = service.Binance{}
	b.binanceService.Init()

	b.router = server.GetEngine().Group("/api/v1/binance")

	b.router.GET("depth", b.getDepth)

}

func (b *Binance) getDepth(c *gin.Context) {

	result, err := b.binanceService.GetDepth()
	if err != nil {
		fmt.Println("\n\033[31m--- Debug ----")
		fmt.Printf("\033[35merr = %+v\n", err)
		fmt.Println("\033[31m\n---------------\033[0m")
	}
	c.JSON(http.StatusOK, result)
}
