package response

type TradeSummary struct {
	TradeCount int    `json:"trade_count"`
	Lots       string `json:"lots"`
	FloatingPL string `json:"floating_pl"`
	ClosedPL   string `json:"closed_pl"`
}
