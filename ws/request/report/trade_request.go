package report

type GetTradeListRequest struct {
	Params TradeListRequestParams `json:"params" validate:"required"`
}

type TradeListRequestParams struct {
	Tickets []int64 `json:"tickets" validate:"required,gt=0,dive,required"`
}
