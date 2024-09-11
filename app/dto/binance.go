package dto

type GetBinanceDepthDto struct {
	Symbol string `form:"symbol" binding:"required"`
}

type GetBinanceBookTickerDto struct {
	Symbol string `form:"symbol" binding:"required"`
}

type WsBinanceDepthDto struct {
	Symbol string `form:"symbol" binding:"required"`
}

type GetBinanceDepthRO struct {
	LastUpdateId int        `json:"lastUpdateId"`
	Bids         [][]string `json:"bids"`
	Asks         [][]string `json:"asks"`
}

type GetBinanceBookTickerRO struct {
	Symbol   string `json:"symbol"`
	BidPrice string `json:"bidPrice"`
	BidQty   string `json:"bidQty"`
	AskPrice string `json:"askPrice"`
	AskQty   string `json:"askQty"`
}
