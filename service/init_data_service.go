package service

import (
	"backend-rms/model"

	"github.com/shopspring/decimal"
)

var InitializedData Data

type Data struct {
	Sources  map[string]*model.Source
	Symbols  map[string]*model.Symbol
	Groups   map[string]*model.Group
	Accounts map[int64]*model.Account
	Trades   map[int64]*model.Trade
}

var SummaryBook = make(map[string]map[string]*InitSummary)

type InitSummary struct {
	Symbol     string          `json:"symbol"`
	BuyLots    decimal.Decimal `json:"buy_lots"`
	SellLots   decimal.Decimal `json:"sell_lots"`
	LpNetLots  decimal.Decimal `json:"lp_net_lots"`
	Et6NetLots decimal.Decimal `json:"et6_net_lots"`
	Orders     int             `json:"orders"`
}

var ExposureBook = make(map[string]map[string]*InitExposure)

type InitExposure struct {
	Symbol      string          `json:"symbol"`
	Currency    string          `json:"currency"`
	Exposure    decimal.Decimal `json:"exposure"`
	ExposureUsd decimal.Decimal `json:"exposure_usd"`
}
