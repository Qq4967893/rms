package request

type AddGroupRequest struct {
	Name string `json:"name" validate:"required,max=20,not_exists=rms^groups^name"`
	Book string `json:"book" validate:"required,len=1"`
}

type UpdateGroupRequest struct {
	Id   int64  `uri:"id" json:"id" validate:"required,numeric"`
	Name string `json:"name" validate:"required,max=20,not_exists=rms^groups^name^Id"`
	Book string `json:"book" validate:"required,len=1"`
}
