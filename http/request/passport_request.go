package request

type GetAccountsRequest struct {
	Passport int64 `form:"passport" validate:"numeric"`
	Login    int64 `form:"login" validate:"numeric"`
}
