package model

import "time"

type LpSpreadRecord struct {
	Source string    `json:"source"`
	LpType int8      `json:"lp_type"`
	LpAvg  int64     `json:"lp_avg"`
	LpDate time.Time `json:"lp_date"`
}

func (LpSpreadRecord) TableName() string {
	return "lp_spread_record"
}
