package model

type GroupPair struct {
	BaseModel

	GroupA string `json:"group_a"`
	GroupB string `json:"group_b"`

	CommonTimestampsField
}
