package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type WalletUser struct {
	BaseModel

	Passport     int64           `json:"passport"`
	Balance      decimal.Decimal `json:"balance"`
	Currency     string          `json:"currency"`
	FirstPayTime time.Time       `json:"first_pay_time"`

	CommonTimestampsField
}
