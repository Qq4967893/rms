package report

type GetAccountListRequest struct {
	Params AccountListRequestParams `json:"params" validate:"required"`
}
type AccountListRequestParams struct {
	Logins []int64 `json:"logins" validate:"required,gt=0,dive,required"`
}
