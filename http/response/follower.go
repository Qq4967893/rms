package response

import (
	"github.com/shopspring/decimal"
)

type Follower struct {
	Mt4Id           int64  `json:"mt4_id"`
	UserCode        int64  `json:"user_code"`
	PeriodPl        string `json:"period_pl"`
	ClosedPl        string `json:"closed_pl"`
	FloatingPl      string `json:"floating_pl"`
	Balance         string `json:"balance"`
	MarginLevel     string `json:"margin_level"`
	FollowMasterNum int    `json:"follow_master_num"`
	CopiedAmount    string `json:"copied_amount"`
	CopiedAllAmount string `json:"copied_all_amount"`
	Group           string `json:"group"`
	IsFollowing     string `json:"is_following"`
	StartTime       string `json:"start_time"`
}

type FollowerRecord struct {
	StartTime  string          `json:"start_time"`
	EndTime    string          `json:"end_time"`
	Commission decimal.Decimal `json:"commission"`
	ClosedPl   decimal.Decimal `json:"closed_pl"`
}
