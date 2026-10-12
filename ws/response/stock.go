package response

import "github.com/shopspring/decimal"

type MarketStock struct {
	Symbol      string          `json:"symbol"`
	BuyLots     decimal.Decimal `json:"buy_lots"`
	SellLots    decimal.Decimal `json:"sell_lots"`
	Et6NetLots  decimal.Decimal `json:"et6_net_lots"`
	LPNetLots   decimal.Decimal `json:"lp_net_lots"`
	Orders      int             `json:"orders"`
	BuyPrice    decimal.Decimal `json:"buy_price"`
	SellPrice   decimal.Decimal `json:"sell_price"`
	ExposureUSD decimal.Decimal `json:"exposure_USD"`
}
