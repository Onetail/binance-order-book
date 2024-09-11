package controller

import (
	"binance-order-book/app"
	"binance-order-book/app/application"
	"binance-order-book/app/dto"
	"binance-order-book/app/service"
	"binance-order-book/app/utils"
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
	b.router.GET("ticker/bookTicker", b.getBookTicker)

	b.router = server.GetEngine().Group("/ws/v1/binance")
	b.router.GET("depth", b.wsBinance)

}

func (b *Binance) getDepth(c *gin.Context) {
	var data dto.GetBinanceDepthDto

	err := c.ShouldBindQuery(&data)
	if err != nil {
		application.HandleError(c, application.NewError(http.StatusBadRequest, err.Error()))
		return
	}

	result, getErr := b.binanceService.GetDepth(data)
	if getErr != nil {
		application.HandleError(c, application.NewError(http.StatusForbidden, getErr.Error()))
		return
	}
	c.JSON(http.StatusOK, result)
}

func (b *Binance) getBookTicker(c *gin.Context) {
	var data dto.GetBinanceBookTickerDto

	err := c.ShouldBindQuery(&data)
	if err != nil {
		application.HandleError(c, application.NewError(http.StatusBadRequest, err.Error()))
		return
	}

	result, getErr := b.binanceService.GetBookTicker(data)
	if getErr != nil {
		application.HandleError(c, application.NewError(http.StatusForbidden, getErr.Error()))
		return
	}
	c.JSON(http.StatusOK, result)
}

func (b *Binance) wsBinance(c *gin.Context) {

	var data dto.WsBinanceDepthDto
	err := c.ShouldBindQuery(&data)
	if err != nil {
		application.HandleError(c, application.NewError(http.StatusBadRequest, err.Error()))
		return
	}

	utils.WsManager.RegisterClient(c, data)
}
