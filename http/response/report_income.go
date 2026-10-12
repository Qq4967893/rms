package response

import "github.com/shopspring/decimal"

type ReportIncomeSymbolResponse struct {
	Date               string `json:"date"`
	Symbol             string `json:"symbol"`
	Profit             string `json:"profit"`
	Swaps              string `json:"swaps"`
	Commission         string `json:"commission"`
	RawSpreadProfit    string `json:"raw_spread_profit"`
	MarkUpSpreadProfit string `json:"mark_up_spread_profit"`
	Lots               string `json:"lots"`
}

type ReportIncomeAccountResponse struct {
	Login              int64  `json:"login"`
	Profit             string `json:"profit"`
	Swaps              string `json:"swaps"`
	Commission         string `json:"commission"`
	RawSpreadProfit    string `json:"raw_spread_profit"`
	MarkUpSpreadProfit string `json:"mark_up_spread_profit"`
	Lots               string `json:"lots"`
}

type ReportIncomeTradeResponse struct {
	Login              int64           `json:"login"`
	Ticket             int64           `json:"ticket"`
	Profit             string          `json:"profit"`
	Swaps              string          `json:"swaps"`
	Commission         string          `json:"commission"`
	RawSpreadProfit    string          `json:"raw_spread_profit"`
	MarkUpSpreadProfit string          `json:"mark_up_spread_profit"`
	Lots               string          `json:"lots"`
	OpenPrice          decimal.Decimal `json:"open_price"`
	ClosePrice         decimal.Decimal `json:"close_price"`
	OpenTime           string          `json:"open_time"`
	CloseTime          string          `json:"close_time"`
}
