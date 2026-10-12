package request

type GetChangeGroupRecordListRequest struct {
	Login     int64 `form:"login" validate:"numeric"`
	ManagerId int64 `form:"manager_id" validate:"numeric"`
	Limit     int   `form:"limit" validate:"required,numeric,min=1,max=100"`
	Offset    int   `form:"offset" validate:"numeric,min=0"`
}
