package ws

import (
	"binance-order-book/app"
	"binance-order-book/app/dto"
	"binance-order-book/app/service"
	"binance-order-book/app/utils"
	"encoding/json"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type Binance struct {
	app            *app.App
	binanceService service.Binance
	symbolList     []string
}

func (b *Binance) Init() {
	b.binanceService = service.Binance{}
	b.binanceService.Init()
	b.symbolList = []string{"ethbtc"}

}

func (b *Binance) emitWsEvent(symbol string, event dto.GetBinanceDepthRO) {
	bytesEvent, _ := json.Marshal(event)
	utils.WsManager.GroupBroadcast(symbol, bytesEvent)
}

func (b *Binance) WsDepth() {
	defer func() {
		if err := recover(); err != nil {
			logrus.WithFields(logrus.Fields{
				"title": "[ws] binance",
				"error": err,
			}).Error()
			return
		}
	}()

	for _, i := range b.symbolList {
		symbol := i
		wsSymbolData := dto.WsBinanceDepthDto{
			Symbol: symbol,
		}
		b.binanceService.WsDepth(wsSymbolData, func(ws *websocket.Conn, event dto.GetBinanceDepthRO) {

			b.emitWsEvent(symbol, event)
		})
	}

}
