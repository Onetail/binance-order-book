package utils

import (
	"binance-order-book/app/dto"
	"encoding/json"
	"fmt"

	"sync"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var WsManager = clientManager{
	clientGroup: &sync.Map{},
	register:    make(chan *wsClient),
	unRegister:  make(chan *wsClient),
	broadcast:   make(chan *boradcastData, 100),
}

type clientManager struct {
	clientGroup *sync.Map
	register    chan *wsClient
	unRegister  chan *wsClient
	broadcast   chan *boradcastData
}

type boradcastData struct {
	Symbol string
	Data   []byte
}

type wsClient struct {
	ID     string
	Group  string
	Socket *websocket.Conn
	Send   chan []byte
	Symbol string
}

func (c *wsClient) PushException(msg string) {
	wsEventReceiveDto := dto.WsEventDto{
		Event: dto.Exception,
		Msg:   string(msg),
	}
	bytes, _ := json.Marshal(wsEventReceiveDto)
	c.Send <- []byte(string(bytes))
}

func (c *wsClient) Read(manager *clientManager) {
	defer func() {
		WsManager.unRegister <- c
		c.Socket.Close()
	}()

	for {
		_, msg, err := c.Socket.ReadMessage()
		if err != nil {
			break
		}
		var data dto.WsEventDto
		json.Unmarshal([]byte(msg), &data)
		bytes, _ := json.Marshal(data)
		c.Send <- []byte(fmt.Sprintf("%s", bytes))

	}
}

func (c *wsClient) Write() {
	defer func() {
		c.Socket.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			c.Socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}

func (manager *clientManager) Start() {

	for {
		select {
		case client := <-manager.register:
			clientG, _ := manager.clientGroup.LoadOrStore(client.Symbol, &sync.Map{})
			clientGroup := clientG.(*sync.Map)
			clientGroup.Store(client.ID, client)

		case client := <-manager.unRegister:
			clientG, _ := manager.clientGroup.LoadOrStore(client.Symbol, &sync.Map{})
			clientGroup := clientG.(*sync.Map)
			clientGroup.LoadAndDelete(client.ID)
			close(client.Send)

		case data := <-manager.broadcast:
			v, ok := manager.clientGroup.Load(data.Symbol)
			if !ok {
				continue
			}
			vv, _ := v.(*sync.Map)
			vv.Range(func(key, value any) bool {
				if value.(*wsClient) != nil {
					value.(*wsClient).Send <- data.Data
				}

				return true
			})

		}
	}
}

func (manager *clientManager) RegisterClient(ctx *gin.Context, data dto.WsBinanceDepthDto) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		Subprotocols: []string{ctx.GetHeader("Sec-WebSocket-Protocol")},
	}

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		return
	}

	client := &wsClient{
		ID:     (uuid.New()).String(),
		Socket: conn,
		Send:   make(chan []byte, 1024),
		Symbol: data.Symbol,
	}

	manager.register <- client
	go client.Read(manager)
	go client.Write()
}

func (manager *clientManager) GroupBroadcast(symbol string, message []byte) {
	data := &boradcastData{
		Symbol: symbol,
		Data:   message,
	}
	manager.broadcast <- data
}
