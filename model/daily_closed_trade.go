package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type DailyClosedTrade struct {
	BaseModel

	Ticket                 int64           `json:"ticket"`
	CloseDate              time.Time       `json:"close_date"`
	SpreadProfit           decimal.Decimal `json:"spread_profit"`
	WithoutRawSpreadProfit decimal.Decimal `json:"without_raw_spread_profit"`

	CommonTimestampsField
}

type SymbolDailyClosedTrade struct {
	Symbol             string          `json:"symbol"`
	Profit             decimal.Decimal `json:"profit"`
	Swaps              decimal.Decimal `json:"swaps"`
	Commission         decimal.Decimal `json:"commission"`
	RawSpreadProfit    decimal.Decimal `json:"raw_spread_profit"`
	MarkUpSpreadProfit decimal.Decimal `json:"mark_up_spread_profit"`
	Lots               decimal.Decimal `json:"lots"`
}

type AccountDailyClosedTrade struct {
	Login              int64           `json:"login"`
	Profit             decimal.Decimal `json:"profit"`
	Swaps              decimal.Decimal `json:"swaps"`
	Commission         decimal.Decimal `json:"commission"`
	RawSpreadProfit    decimal.Decimal `json:"raw_spread_profit"`
	MarkUpSpreadProfit decimal.Decimal `json:"mark_up_spread_profit"`
	Lots               decimal.Decimal `json:"lots"`
}

type DailyClosedTradeInfo struct {
	Login              int64           `json:"login"`
	Symbol             string          `json:"symbol"`
	Ticket             int64           `json:"ticket"`
	Profit             decimal.Decimal `json:"profit"`
	Swaps              decimal.Decimal `json:"swaps"`
	Commission         decimal.Decimal `json:"commission"`
	RawSpreadProfit    decimal.Decimal `json:"raw_spread_profit"`
	MarkUpSpreadProfit decimal.Decimal `json:"mark_up_spread_profit"`
	Lots               decimal.Decimal `json:"lots"`
	OpenPrice          decimal.Decimal `json:"open_price"`
	ClosePrice         decimal.Decimal `json:"close_price"`
	OpenTime           time.Time       `json:"open_time"`
	CloseTime          time.Time       `json:"close_time"`
}

func (DailyClosedTrade) TableName() string {
	return "daily_closed_trades"
}
