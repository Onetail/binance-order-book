package app

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type App struct {
	HTTPServer *HTTPServer
}

func (app *App) Init() {

	if os.Getenv("GO_ENV") == "production" {
		viper.SetConfigName("config.prod")
		viper.AddConfigPath("./config")
	} else if os.Getenv("GO_ENV") == "dev" {
		viper.SetConfigName("config.dev")
		viper.AddConfigPath("./config")
	} else {
		viper.SetConfigName("config.local")
		viper.AddConfigPath("./config")
	}

	err := viper.ReadInConfig()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"title": "app read config error",
			"error": err,
		}).Error()
		panic(err)
	}

	app.HTTPServer = &HTTPServer{}
	app.HTTPServer.Init(app)

}

func (app *App) Run() {
	app.HTTPServer.Start()
}
