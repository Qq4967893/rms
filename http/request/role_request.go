package request

type GetDisplayMenusRequest struct {
	Id int64 `uri:"id" validate:"required,numeric"`
}

type AddRoleRequest struct {
	Name    string `json:"name" validate:"required,not_exists=rms^roles^name"`
	MenuIds []int  `json:"menu_ids" validate:"required"`
}

type DeleteRoleRequest struct {
	Id int64 `uri:"id" validate:"required,numeric"`
}
