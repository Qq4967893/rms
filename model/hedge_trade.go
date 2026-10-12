package model

import "time"

type HedgeTrade struct {
	BaseModel

	Ticket              int64     `json:"ticket"`
	Login               int64     `json:"login"`
	Symbol              string    `json:"symbol"`
	Cmd                 int8      `json:"cmd"`
	Volume              int64     `json:"volume"`
	OpenTime            time.Time `json:"open_time"`
	HedgeAccount        int64     `json:"hedge_account"`
	ChangeGroupRecordId int64     `json:"change_group_record_id"`
	ManagerId           int64     `json:"manager_id"`

	CommonTimestampsField
}
