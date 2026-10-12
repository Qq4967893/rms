package request

type AddMenuRequest struct {
	Name    string `json:"name" validate:"required,not_exists=rms^menus^name"`
	Pid     int64  `json:"pid"`
	Path    string `json:"path" validate:"required"`
	Display int8   `json:"display"`
	Sort    int8   `json:"sort"`
}

type UpdateMenuRequest struct {
	Id      int64  `uri:"id" validate:"required,numeric"`
	Name    string `json:"name" validate:"required,not_exists=rms^menus^name^Id"`
	Pid     int64  `json:"pid" validate:"numeric"`
	Path    string `json:"path" validate:"required"`
	Display int8   `json:"display"`
	Sort    int8   `json:"sort"`
}

type DeleteMenuRequest struct {
	Id int64 `uri:"id" validate:"required,numeric"`
}
