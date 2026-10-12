package request

type UpdateHedgeAccountRequest struct {
	Id      int64  `uri:"id" json:"id" validate:"required,numeric"`
	RoleIds string `json:"role_ids" validate:"required"`
	Desc    string `json:"desc"`
}
