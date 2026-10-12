package model

type Event struct {
	Ticket int64 `json:"ticket"`
	Action int8  `json:"action"`
	Type   int8  `json:"type"`
}

func (Event) TableName() string {
	return "event"
}
