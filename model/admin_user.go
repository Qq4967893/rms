package model

type AdminUser struct {
	BaseModel

	Username  string `json:"username"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	City      string `json:"city"`
	Address   string `json:"address"`
	Note      string `json:"note"`
	LeaderID  int64  `json:"leader_id"`
	RoleID    int64  `json:"role_id"`
	Status    int8   `json:"status"`
	LevelID   int64  `json:"level_id"`
	TransSet  int64  `json:"trans_set"`
	SalesTeam string `json:"sales_team"`

	CommonTimestampsField
}
