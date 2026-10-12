package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type Account struct {
	Login    int64           `json:"login"`
	Group    string          `json:"group"`
	Currency string          `json:"currency"`
	Balance  decimal.Decimal `json:"balance"`
	Regdate  time.Time       `json:"regdate"`
}

type ReportAccount struct {
	Login   int64           `json:"login"`
	Group   string          `json:"group"`
	Balance decimal.Decimal `json:"balance"`
}

func (Account) TableName() string {
	return "account"
}
