package main

import (
	"binance-dashboard/app"
	"binance-dashboard/app/controller"
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
