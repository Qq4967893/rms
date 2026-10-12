package model

import "time"

type SwitchLpRecord struct {
	Source     string    `json:"source"`
	OldLp      int8      `json:"old_lp"`
	NewLp      int8      `json:"new_lp"`
	SwitchTime time.Time `json:"switch_time"`
}

func (SwitchLpRecord) TableName() string {
	return "switch_lp_record"
}
