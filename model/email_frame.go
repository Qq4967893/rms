package model

type EmailFrame struct {
	BaseModel

	Content string `json:"connect"`

	CommonTimestampsField
}
