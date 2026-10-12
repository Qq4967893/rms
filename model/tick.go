package model

import "github.com/shopspring/decimal"

type Tick struct {
	Ask decimal.Decimal
	Bid decimal.Decimal
	Ctm int64
}
