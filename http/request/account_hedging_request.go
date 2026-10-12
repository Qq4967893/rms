package request

type GetAccountHedgingPreviewRequest struct {
	Login   int64  `form:"login" validate:"required,numeric"`
	ToGroup string `form:"to_group" validate:"required"`
}

type GetAccountHedgingConfirmRequest struct {
	Login   int64  `form:"login" validate:"required,numeric"`
	ToGroup string `form:"to_group" validate:"required"`
}

type HedgeAccountHedgingRequest struct {
	Login   int64  `json:"login" validate:"required,numeric"`
	Group   string `json:"group" validate:"required"`
	ToGroup string `json:"to_group" validate:"required"`
}
