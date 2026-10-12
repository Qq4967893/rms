package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type PendingTrade struct {
	Ticket    int64           `json:"ticket"`
	Login     int64           `json:"login"`
	Cmd       int8            `json:"cmd"`
	Symbol    string          `json:"symbol"`
	Volume    int64           `json:"volume"`
	OpenPrice decimal.Decimal `json:"open_price"`
	OpenTime  time.Time       `json:"open_time"`
}

// -1所有，0买，1卖
func GetPendingCmd(cmd int8) []int8 {
	if cmd == 0 {
		return []int8{2, 4}
	} else if cmd == 1 {
		return []int8{3, 5}
	} else {
		return []int8{2, 3, 4, 5}
	}
}
