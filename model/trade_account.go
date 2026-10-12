package model

import (
	"time"
)

type TradeAccount struct {
	BaseModel

	Passport         int64     `json:"passport"`
	TradeAccountID   int64     `json:"trade_account_id"`
	TradeAccountType int8      `json:"trade_account_type"`
	Currency         string    `json:"currency"`
	Name             string    `json:"name"`
	Envir            int8      `json:"envir"`
	Type             int8      `json:"type"`
	FirstPayTime     time.Time `json:"first_pay_time"`
	LastActiveTime   time.Time `json:"last_active_time"`
	Default          int8      `json:"default"`

	CommonTimestampsField
}
