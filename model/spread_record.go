package model

import "time"

type SpreadRecord struct {
	Symbol       string    `json:"symbol"`
	BeforeSpread int64     `json:"before_spread"`
	AfterSpread  int64     `json:"after_spread"`
	ChangeTime   time.Time `json:"change_time"`
}

func (SpreadRecord) TableName() string {
	return "spread"
}
