package utils

import "github.com/shopspring/decimal"

func SizeLimit(inputList [][]string, maxLimit decimal.Decimal) [][]string {
	bidsList := [][]string{}
	currentBid := decimal.NewFromInt(0)

	for _, bid := range inputList {

		size, _ := decimal.NewFromString(bid[0])
		price, _ := decimal.NewFromString(bid[1])
		currentBid = currentBid.Add(size.Mul(price))
		if currentBid.Cmp(maxLimit) == 1 {
			continue
		}
		bidsList = append(bidsList, bid)
	}

	return bidsList
}
