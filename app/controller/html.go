package controller

import (
	"binance-order-book/app"
	"binance-order-book/app/application"
	"binance-order-book/app/dto"
	"binance-order-book/app/service"
	"binance-order-book/app/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

type Html struct {
	app    *app.App
	server *app.HTTPServer
	router *gin.RouterGroup

	binanceServer service.Binance

	bidMaxLimit decimal.Decimal
	askMaxLimit decimal.Decimal
}

func (h *Html) Init(server *app.HTTPServer) {

	h.server = server
	h.app = server.App

	h.binanceServer = service.Binance{}
	h.binanceServer.Init()

	h.askMaxLimit = decimal.NewFromInt(150)
	h.bidMaxLimit = decimal.NewFromInt(5)

	h.router = server.GetEngine().Group("")
	h.router.GET("", h.getHtml)
}

func (h *Html) getHtml(c *gin.Context) {

	var data dto.GetHtmlDto
	err := c.ShouldBindQuery(&data)
	if err != nil {
		application.HandleError(c, application.NewError(http.StatusBadRequest, err.Error()))
		return
	}

	depthDto := dto.GetBinanceDepthDto{
		Symbol: data.Symbol,
	}
	depthData, err := h.binanceServer.GetDepth(depthDto)
	if err != nil {
		application.HandleError(c, application.NewError(http.StatusForbidden, err.Error()))
		return
	}
	depthData.Bids = utils.SizeLimit(depthData.Bids, h.bidMaxLimit)
	depthData.Asks = utils.SizeLimit(depthData.Asks, h.askMaxLimit)

	bookTickerDto := dto.GetBinanceBookTickerDto{
		Symbol: data.Symbol,
	}
	bookTickerData, err := h.binanceServer.GetBookTicker(bookTickerDto)
	if err != nil {

		application.HandleError(c, application.NewError(http.StatusForbidden, err.Error()))
		return
	}

	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"title":          data.Symbol,
		"bids":           depthData.Bids,
		"asks":           depthData.Asks,
		"bookTickerData": bookTickerData,
	})

}
