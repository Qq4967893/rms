package model

type User struct {
	BaseModel

	Passport      int64  `json:"passport"`
	Username      string `json:"username"`
	Avatar        string `json:"avatar"`
	Country       string `json:"country"`
	Language      string `json:"language"`
	AreaID        int8   `json:"area_id"`
	Status        int8   `json:"status"`
	AffiliationID int64  `json:"affiliation_id"`
	Type          int8   `json:"type"`

	CommonTimestampsField
}
