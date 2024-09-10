package app

import (
	"binance-dashboard/app/utils"
	"log"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type HTTPServer struct {
	App    *App
	engine *gin.Engine
	host   string
	port   int
	name   string
}

func (hs *HTTPServer) Init(app *App) {

	hs.App = app
	hs.host = viper.GetString("http_server.host")
	hs.port = viper.GetInt("http_server.port")
	hs.name = viper.GetString("http_server.name")

	hs.engine = gin.New()
	hs.engine.Use(utils.LoggerToFile())
	hs.engine.MaxMultipartMemory = 10 << 20 // 10 MiB

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowMethods = []string{"GET", "POST", "DELETE", "PATCH", "OPTIONS", "PUT"}
	corsConfig.AllowHeaders = []string{"Authorization", "Content-Type", "Upgrade", "Origin",
		"Connection", "Accept-Encoding", "Accept-Language", "Host", "Access-Control-Request-Method", "Access-Control-Request-Headers"}
	hs.engine.Use(gin.Recovery())
	hs.engine.Use(cors.New(corsConfig))
	hs.engine.Use(gin.LoggerWithConfig(gin.LoggerConfig{SkipPaths: []string{"/api/v1/messenger/callback", "/ws/v1/symbol/quote"}}))

	hs.engine.LoadHTMLGlob("app/public/*")

}

func (hs *HTTPServer) Start() {
	log.Printf("Listening on http://%s:%d", hs.host, hs.port)
	hs.engine.Run(hs.host + ":" + strconv.Itoa(hs.port))
}

func (hs *HTTPServer) GetEngine() *gin.Engine {
	return hs.engine
}
