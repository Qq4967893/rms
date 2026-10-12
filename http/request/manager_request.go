package request

type GetManagerListRequest struct {
	UserName string `form:"username"`
	Offset   int    `form:"offset" validate:"numeric"`
	Limit    int    `form:"limit" validate:"required,numeric,min=10,max=20"`
}

type AddManagerRequest struct {
	UserName string `json:"username" validate:"required,max=20,not_exists=rms^managers^username"`
	Password string `json:"password" validate:"required"`
	Email    string `json:"email" validate:"required,email,not_exists=rms^managers^email"`
	RoleId   int64  `json:"role_id" validate:"required,numeric"`
}

type UpdateManagerRequest struct {
	Id     int64  `uri:"id" validate:"required,numeric"`
	Email  string `json:"email" validate:"required,email,not_exists=rms^managers^email^Id"`
	RoleId int64  `json:"role_id" validate:"required,numeric"`
}

type DeleteManagerRequest struct {
	Id int64 `uri:"id" validate:"required,numeric"`
}

type ResetManagerPasswordRequest struct {
	Id int64 `uri:"id" validate:"required,numeric"`
}

type ModifyManagerPasswordRequest struct {
	OldPassword       string `json:"password" validate:"required"`
	NewPassword       string `json:"new_password" validate:"required"`
	RepeatNewPassword string `json:"repeat_new_password" validate:"required"`
}
