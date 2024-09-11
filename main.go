package main

import (
	"binance-order-book/app"
	"binance-order-book/app/controller"
	"binance-order-book/app/utils"
	"binance-order-book/app/ws"
)

func main() {

	app := app.App{}
	app.Init()

	binanceController := controller.Binance{}
	binanceController.Init(app.HTTPServer)

	heartBeatController := controller.HeartBeat{}
	heartBeatController.Init(app.HTTPServer)

	htmlController := controller.Html{}
	htmlController.Init(app.HTTPServer)

	binanceWebsocket := ws.Binance{}
	binanceWebsocket.Init()
	go binanceWebsocket.WsDepth()
	go utils.WsManager.Start()

	app.Run()
}
