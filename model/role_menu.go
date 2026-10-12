package model

type RoleMenu struct {
	BaseModel

	Role_id int64 `json:"role_id"`
	Menu_id int   `json:"menu_id"`

	CommonTimestampsField
}
