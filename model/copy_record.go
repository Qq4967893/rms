package model

import (
	"time"

	"github.com/shopspring/decimal"
)

const Copy_Record_Default_EndTime = "0000-00-00 00:00:00"

type CopyRecord struct {
	BaseModel

	Mt4From   int64           `json:"mt4_from"`
	Mt4To     int64           `json:"mt4_to"`
	StartTime time.Time       `json:"start_time"`
	EndTime   time.Time       `json:"end_time"`
	Amount    decimal.Decimal `json:"amount"`
	Status    int             `json:"status"`
}

type CopiesRecord struct {
	Mt4To     int64     `json:"mt4_to"`
	UserCode  int64     `json:"user_code"`
	Group     string    `json:"group"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

type CopiesMasters struct {
	Mt4To   int64           `json:"mt4_to"`
	Mt4From int64           `json:"mt4_from"`
	Amount  decimal.Decimal `json:"amount"`
}
