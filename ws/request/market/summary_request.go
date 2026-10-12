package market

type GetSummaryListRequest struct {
	Params SummaryListRequestParams `json:"params" validate:"required"`
}

type SummaryListRequestParams struct {
	Books []string `json:"books" validate:"required"`
	Sort  string   `json:"sort" validate:"required,oneof=buy_lots sell_lots et6_net_lots orders symbol"`
	Order string   `json:"order" validate:"required,oneof=asc desc"`
}
