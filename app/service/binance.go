package service

import (
	"binance-dashboard/app/dto"
	"binance-dashboard/app/utils"
	"fmt"
	"net/url"
)

type Binance struct {
	url string
}

func (b *Binance) Init() {

	b.url = "https://api.binance.com"
}

func (b *Binance) GetDepth() (*dto.GetBinanceDepthRO, error) {
	params := url.Values{}
	params.Add("symbol", "ETHBTC")
	apiUrl := fmt.Sprintf("%s/api/v3/depth?%s", b.url, params.Encode())

	return utils.CallAPI[dto.GetBinanceDepthRO](apiUrl, "GET", nil)
}
