package model

import "github.com/shopspring/decimal"

type ChangeGroupRecord struct {
	BaseModel

	Login       int64           `json:"login"`
	BeforeGroup string          `json:"before_group"`
	AfterGroup  string          `json:"after_group"`
	OpenTicket  string          `json:"open_ticket"`
	FloatingPl  decimal.Decimal `json:"floating_pl"`
	ManagerId   int64           `json:"manager_id"`
	Status      int8            `json:"status"`

	CommonTimestampsField
}

type Operator struct {
	Username  string `json:"operator"`
	ManagerId int64  `json:"manager_id"`
}
