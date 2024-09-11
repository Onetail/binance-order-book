package service

import (
	"binance-order-book/app/dto"
	"binance-order-book/app/utils"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type Binance struct {
	apiUrl string
	wsUrl  string
}

func (b *Binance) Init() {
	b.apiUrl = "https://api.binance.com"
	b.wsUrl = "wss://stream.binance.com:9443"
}

func (b *Binance) GetDepth(data dto.GetBinanceDepthDto) (*dto.GetBinanceDepthRO, error) {
	params := url.Values{}
	params.Add("symbol", data.Symbol)
	params.Add("limit", strconv.Itoa(20))

	apiUrl := fmt.Sprintf("%s/api/v3/depth?%s", b.apiUrl, params.Encode())
	return utils.CallAPI[dto.GetBinanceDepthRO](apiUrl, "GET", nil)
}

func (b *Binance) GetBookTicker(data dto.GetBinanceBookTickerDto) (*dto.GetBinanceBookTickerRO, error) {
	params := url.Values{}
	params.Add("symbol", data.Symbol)
	apiUrl := fmt.Sprintf("%s/api/v3/ticker/bookTicker?%s", b.apiUrl, params.Encode())

	return utils.CallAPI[dto.GetBinanceBookTickerRO](apiUrl, "GET", nil)
}

func (b *Binance) WsDepth(data dto.WsBinanceDepthDto, handler func(ws *websocket.Conn, event dto.GetBinanceDepthRO)) {
	wsUrl := fmt.Sprintf("%s/ws/%s@depth20@100ms", b.wsUrl, data.Symbol)

	errHandler := func(err error) {
		logrus.WithFields(logrus.Fields{
			"title": "ws depth",
			"error": err,
		}).Error()
	}

	for {
		cfg := utils.NewWsConfig(wsUrl)
		wsHandler := func(ws *websocket.Conn, message []byte) {
			var obj dto.GetBinanceDepthRO
			if err := json.Unmarshal(message, &obj); err != nil {
				errHandler(err)
			}
			handler(ws, obj)
		}

		doneC, _, _, err := utils.WsServe(cfg, wsHandler, errHandler)
		if err != nil {
			time.Sleep(time.Second * 3)
			continue
		}
		<-doneC
		time.Sleep(time.Second * 3)
	}
}
