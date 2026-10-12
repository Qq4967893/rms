package request

type GetSymbolHoldingAccountListRequest struct {
	Book   string `form:"book" validate:"required,oneof=all a b"`
	Symbol string `form:"symbol" validate:"required"`
	Limit  int    `form:"limit" validate:"required,numeric,min=1,max=100"`
	Offset int    `form:"offset" validate:"numeric,min=0"`
}

type GetSymbolTradeListRequest struct {
	Book   string `form:"book" validate:"required,oneof=all a b"`
	Symbol string `form:"symbol" validate:"required"`
	Cmd    int8   `form:"cmd" validate:"oneof=-1 0 1"`
	Limit  int    `form:"limit" validate:"required,numeric,min=1,max=100"`
	Offset int    `form:"offset" validate:"numeric,min=0"`
}

type GetSymbolPendingTradeListRequest struct {
	Book   string `form:"book" validate:"required,oneof=all a b"`
	Symbol string `form:"symbol" validate:"required"`
	Cmd    int8   `form:"cmd" validate:"oneof=-1 0 1"`
	Limit  int    `form:"limit" validate:"required,numeric,min=1,max=100"`
	Offset int    `form:"offset" validate:"numeric,min=0"`
}
