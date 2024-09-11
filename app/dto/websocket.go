package dto

type WsEventEnum string

const (
	WsSubscribe   WsEventEnum = "subscribe"
	WsUnSubscribe WsEventEnum = "unSubscribe"
	Exception     WsEventEnum = "exception"
)

type WsEventDto struct {
	Event WsEventEnum        `json:"event"  binding:"required"`
	Data  *WsEventDataDto    `json:"data,omitempty" `
	Items *[]*WsEventDataDto `json:"items,omitempty" `
	Msg   string             `json:"msg,omitempty" `
}

type WsEventDataDto struct {
	Symbol string `json:"symbol,omitempty" `
}
