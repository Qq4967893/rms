package model

type Group struct {
	BaseModel

	Name string `json:"name"`
	Book string `json:"book"`

	CommonTimestampsField
}
