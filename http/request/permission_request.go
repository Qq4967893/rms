package request

type GetPermissionInfoRequest struct {
	Id int64 `uri:"id" validate:"required,numeric"`
}

type UpdatePermissionRequest struct {
	Id      int64  `uri:"id" validate:"required,numeric"`
	Name    string `json:"name" validate:"required,max=20,not_exists=rms^roles^name^Id"`
	MenuIds []int  `json:"menu_ids" validate:"required"`
}
