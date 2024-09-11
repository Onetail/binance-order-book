package utils

import (
	"binance-order-book/app/config"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type WsHandler func(ws *websocket.Conn, message []byte)
type ErrHandler func(err error)
type WsConfig struct {
	Endpoint string
}

func NewWsConfig(endpoint string) *WsConfig {
	return &WsConfig{
		Endpoint: endpoint,
	}
}

var WsServe = func(cfg *WsConfig, handler WsHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, c *websocket.Conn, err error) {
	Dialer := websocket.Dialer{
		Proxy:             http.ProxyFromEnvironment,
		HandshakeTimeout:  45 * time.Second,
		EnableCompression: false,
	}

	c, _, err = Dialer.Dial(cfg.Endpoint, nil)
	if err != nil {
		return nil, nil, c, err
	}

	c.SetReadLimit(config.WebsocketMaxBytes)
	c.SetReadDeadline(time.Now().Add(config.WebsocketTimeout))

	doneC = make(chan struct{})
	stopC = make(chan struct{})

	go func() {
		defer close(doneC)

		keepAlive(c, config.WebsocketKeepAlive)

		silent := false

		go func() {
			select {
			case <-stopC:
				silent = true
			case <-doneC:
			}
			c.Close()
		}()
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				if !silent {
					errHandler(err)
				}
				return
			}
			handler(c, message)
		}
	}()
	return
}

func keepAlive(c *websocket.Conn, timeout time.Duration) {
	ticker := time.NewTicker(timeout)

	lastResponse := time.Now()
	c.SetPongHandler(func(msg string) error {
		lastResponse = time.Now()
		return nil
	})

	go func() {
		defer ticker.Stop()
		for {
			deadline := time.Now().Add(10 * time.Second)
			err := c.WriteControl(websocket.PingMessage, []byte{}, deadline)
			if err != nil {
				return
			}
			<-ticker.C
			if time.Since(lastResponse) > timeout {
				c.Close()
				return
			}
		}
	}()
}
