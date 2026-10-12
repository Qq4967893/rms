package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type Trade struct {
	Ticket          int64           `json:"ticket"`
	Type            int8            `json:"type"`
	Login           int64           `json:"login"`
	Cmd             int8            `json:"cmd"`
	Symbol          string          `json:"symbol"`
	Volume          int64           `json:"volume"`
	SL              decimal.Decimal `json:"sl"`
	TP              decimal.Decimal `json:"tp"`
	Comment         string          `json:"comment"`
	Commission      decimal.Decimal `json:"commission"`
	Swaps           decimal.Decimal `json:"swaps"`
	Profit          decimal.Decimal `json:"profit"`
	OpenPrice       decimal.Decimal `json:"open_price"`
	OpenTime        time.Time       `json:"open_time"`
	ClosePrice      decimal.Decimal `json:"close_price"`
	CloseTime       time.Time       `json:"close_time"`
	MarginRate      decimal.Decimal `json:"margin_rate"`
	Digits          int8            `json:"digits"`
	UsdRate         decimal.Decimal `json:"usd_rate"`
	GbpRate         decimal.Decimal `json:"gbp_rate"`
	TransOpenPrice  decimal.Decimal `json:"trans_open_price"`
	TransClosePrice decimal.Decimal `json:"trans_close_price"`
	TransOpenGroup  string          `json:"trans_open_group"`
	TransCloseGroup string          `json:"trans_close_group"`
}

type ReportTrade struct {
	Ticket     int64           `json:"ticket"`
	Type       int8            `json:"type"`
	Login      int64           `json:"login"`
	Cmd        int8            `json:"cmd"`
	Symbol     string          `json:"symbol"`
	Volume     int64           `json:"volume"`
	SL         decimal.Decimal `json:"sl"`
	TP         decimal.Decimal `json:"tp"`
	Commission decimal.Decimal `json:"commission"`
	Swaps      decimal.Decimal `json:"swaps"`
	Profit     decimal.Decimal `json:"profit"`
	OpenPrice  decimal.Decimal `json:"open_price"`
	OpenTime   time.Time       `json:"open_time"`
	ClosePrice decimal.Decimal `json:"close_price"`
	CloseTime  time.Time       `json:"close_time"`
	Group      string          `json:"group"`
}

type AccountClosedPL struct {
	Login    int64
	ClosedPL decimal.Decimal
}

type SymbolClosedPL struct {
	Symbol     string
	Profit     decimal.Decimal
	Swaps      decimal.Decimal
	Commission decimal.Decimal
}

func (sc *SymbolClosedPL) ClosedPL() decimal.Decimal {
	return sc.Profit.Add(sc.Swaps).Add(sc.Commission)
}

type OpeningTrade struct {
	Ticket     int64
	Login      int64
	Symbol     string
	Cmd        int8
	Volume     int64
	OpenTime   time.Time
	OpenPrice  decimal.Decimal
	Commission decimal.Decimal
	Swaps      decimal.Decimal
	SL         decimal.Decimal
	TP         decimal.Decimal
	Group      string
	Profit     decimal.Decimal
}

type ClosedTrade struct {
	Ticket          int64           `json:"ticket"`
	Login           int64           `json:"login"`
	Cmd             int8            `json:"cmd"`
	Symbol          string          `json:"symbol"`
	Volume          int64           `json:"volume"`
	Commission      decimal.Decimal `json:"commission"`
	Swaps           decimal.Decimal `json:"swaps"`
	Profit          decimal.Decimal `json:"profit"`
	OpenPrice       decimal.Decimal `json:"open_price"`
	OpenTime        time.Time       `json:"open_time"`
	ClosePrice      decimal.Decimal `json:"close_price"`
	CloseTime       time.Time       `json:"close_time"`
	TransOpenPrice  decimal.Decimal `json:"trans_open_price"`
	TransClosePrice decimal.Decimal `json:"trans_close_price"`
	TransOpenGroup  string          `json:"trans_open_group"`
	TransCloseGroup string          `json:"trans_close_group"`
	Group           string          `json:"group"`
	DepositCurrency string          `json:"deposit_currency"`
}

const DEFAULT_CLOSE_TIME = "0000-00-00 00:00:00"

func (Trade) TableName() string {
	return "trade_trans_result"
}
