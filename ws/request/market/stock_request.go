package market

type GetStockListRequest struct {
	Params StockListRequestParams `json:"params" validate:"required"`
}

type StockListRequestParams struct {
	Books []string `json:"books" validate:"required"`
	Sort  string   `json:"sort" validate:"required,oneof=buy_lots sell_lots et6_net_lots orders symbol"`
	Order string   `json:"order" validate:"required,oneof=asc desc"`
}
