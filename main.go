package main

import (
	"binance-order-book/app"
	"binance-order-book/app/controller"
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

	app.Run()
}
