package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type CopyPay struct {
	Mt4To     int64           `json:"mt4_to"`
	Mt4From   int64           `json:"mt4_from"`
	PayZone   string          `json:"pay_zone"`
	PayTime   time.Time       `json:"pay_time"`
	PayMaster decimal.Decimal `json:"pay_master"`
}
