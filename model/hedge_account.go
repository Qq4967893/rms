package model

type HedgeAccount struct {
	BaseModel

	HedgeAccount int64  `json:"hedge_account"`
	RoleIds      string `json:"role_ids"`
	Desc         string `json:"desc"`

	CommonTimestampsField
}
