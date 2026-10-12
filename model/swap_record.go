package model

import (
	"github.com/shopspring/decimal"
)

type SwapRecord struct {
	BaseModel

	Ticket         int64           `json:"ticket"`
	Login          int64           `json:"login"`
	Swaps          decimal.Decimal `json:"swaps"`
	IsPartialClose int8            `json:"is_partial_close"`

	CommonTimestampsField
}
