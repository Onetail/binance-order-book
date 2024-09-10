package service

import (
	"binance-order-book/app/dto"
	"binance-order-book/app/utils"
	"fmt"
	"net/url"
)

type Binance struct {
	url string
}

func (b *Binance) Init() {

	b.url = "https://api.binance.com"
}

func (b *Binance) GetDepth(data dto.GetBinanceDepthDto) (*dto.GetBinanceDepthRO, error) {
	params := url.Values{}
	params.Add("symbol", data.Symbol)
	apiUrl := fmt.Sprintf("%s/api/v3/depth?%s", b.url, params.Encode())

	return utils.CallAPI[dto.GetBinanceDepthRO](apiUrl, "GET", nil)
}

func (b *Binance) GetBookTicker(data dto.GetBinanceBookTickerDto) (*dto.GetBinanceBookTickerRO, error) {
	params := url.Values{}
	params.Add("symbol", data.Symbol)
	apiUrl := fmt.Sprintf("%s/api/v3/ticker/bookTicker?%s", b.url, params.Encode())

	return utils.CallAPI[dto.GetBinanceBookTickerRO](apiUrl, "GET", nil)
}
