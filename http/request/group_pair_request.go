package request

type AddGroupPairRequest struct {
	GroupA string `json:"group_a" validate:"required,max=20"`
	GroupB string `json:"group_b" validate:"required,max=20"`
}

type UpdateGroupPairRequest struct {
	Id     int64  `uri:"id" json:"id" validate:"required,numeric"`
	GroupA string `json:"group_a" validate:"required,max=20"`
	GroupB string `json:"group_b" validate:"required,max=20"`
}
