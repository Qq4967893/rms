package model

type Role struct {
	BaseModel

	Name string `json:"name"`

	CommonTimestampsField
}
