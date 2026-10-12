package response

import "github.com/shopspring/decimal"

type MarketSummary struct {
	Symbol     string          `json:"symbol"`
	BuyLots    decimal.Decimal `json:"buy_lots"`
	SellLots   decimal.Decimal `json:"sell_lots"`
	Et6NetLots decimal.Decimal `json:"et6_net_lots"`
	BuyPrice   decimal.Decimal `json:"buy_price"`
	SellPrice  decimal.Decimal `json:"sell_price"`
	Orders     int             `json:"orders"`
}
