package request

type GetPassportDetailInfoRequest struct {
	Passport int64 `form:"passport" validate:"required,numeric"`
}

type GetAccountListRequest struct {
	Login int64 `form:"login" validate:"required,numeric"`
}

type GetTradeListRequest struct {
	Passport int64  `form:"passport" validate:"required,numeric"`
	Login    int64  `form:"login" validate:"numeric"`
	Ticket   int64  `form:"ticket" validate:"numeric"`
	Symbol   string `form:"symbol"`
	Limit    int    `form:"limit" validate:"required,numeric"`
	Offset   int    `form:"offset" validate:"numeric"`
	Sort     string `form:"sort" validate:"required,oneof=open_close_time open_time close_time open_price"`
	Order    string `form:"order" validate:"required,oneof=asc desc"`
}

type GetSymbolListRequest struct {
	Passport int64 `form:"passport" validate:"required,numeric"`
	Login    int64 `form:"login" validate:"numeric"`
}
