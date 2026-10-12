package response

type SymbolLots struct {
	Symbol   string `json:"symbol"`
	Lots     string `json:"lots"`
	BuyLots  string `json:"buy_lots"`
	SellLots string `json:"sell_lots"`
	NetLots  string `json:"net_lots"`
}

type SymbolHoldingAccount struct {
	Login            int64  `json:"login"`
	BuyLots          string `json:"buy_lots"`
	SellLots         string `json:"sell_lots"`
	NetLots          string `json:"net_lots"`
	SymbolFloatingPL string `json:"symbol_floating_pl"`
	FloatingPL       string `json:"floating_pl"`
}

type TotalSymbolHoldingAccount struct {
	Total            int    `json:"total"`
	BuyLots          string `json:"buy_lots"`
	SellLots         string `json:"sell_lots"`
	NetLots          string `json:"net_lots"`
	SymbolFloatingPL string `json:"symbol_floating_pl"`
	FloatingPL       string `json:"floating_pl"`
}

type SymbolTrade struct {
	Ticket      int64  `json:"ticket"`
	Login       int64  `json:"login"`
	OpenTime    string `json:"open_time"`
	OpenPrice   string `json:"open_price"`
	SymbolPrice string `json:"symbol_price"`
	Cmd         string `json:"cmd"`
	Lots        string `json:"lots"`
	FloatingPL  string `json:"floating_pl"`
	Profit      string `json:"profit"`
	Swap        string `json:"swaps"`
	Commission  string `json:"commission"`
	SL          string `json:"sl"`
	TP          string `json:"tp"`
}

type SymbolTotalTrade struct {
	Total      int    `json:"total"`
	Lots       string `json:"lots"`
	FloatingPL string `json:"floating_pl"`
}

type SymbolPendingTrade struct {
	Ticket       int64  `json:"ticket"`
	Login        int64  `json:"login"`
	PendingTime  string `json:"pending_time"`
	PendingPrice string `json:"pending_price"`
	SymbolPrice  string `json:"symbol_price"`
	Cmd          string `json:"cmd"`
	Lots         string `json:"lots"`
}

type SymbolTotalPendingTrade struct {
	Total int    `json:"total"`
	Lots  string `json:"lots"`
}
