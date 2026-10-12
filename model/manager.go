package model

type Manager struct {
	BaseModel

	Username string `json:"username"`
	Password string `json:"-"`
	Email    string `json:"email"`
	RoleId   int64  `json:"role_id"`

	CommonTimestampsField
}

type ManagerGetList struct {
	BaseModel

	Username string `json:"username"`
	Email    string `json:"email"`
	RoleName string `json:"role_name"`
	RoleId   int64  `json:"role_id"`
}
