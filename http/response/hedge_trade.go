package response

type HedgeTradeList struct {
	Account    int64  `json:"account"`
	HedgeOrder int64  `json:"hedge_order"`
	OpenTime   string `json:"open_time"`
	Market     string `json:"market"`
	Cmd        string `json:"cmd"`
	Lot        string `json:"lot"`
	Operator   string `json:"operator"`
}
