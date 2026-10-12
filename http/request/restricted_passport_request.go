package request

type GetRestrictedPassportListRequest struct {
	Limit  int    `form:"limit" validate:"required,numeric,min=1,max=100"`
	Offset int    `form:"offset" validate:"numeric,min=0"`
	Sort   string `form:"sort" validate:"required,oneof=trigger_time"`
	Order  string `form:"order" validate:"required,oneof=asc desc"`
}

type GetRestrictedPassportRecordListRequest struct {
	Passport int64  `form:"passport" `
	Limit    int    `form:"limit" validate:"required,numeric,min=1,max=100"`
	Offset   int    `form:"offset" validate:"numeric,min=0"`
	Sort     string `form:"sort" validate:"required,oneof=remove_time"`
	Order    string `form:"order" validate:"required,oneof=asc desc"`
}

type RemoveRestrictedPassportRequest struct {
	Id int64 `uri:"id" json:"id" validate:"required,numeric"`
}
