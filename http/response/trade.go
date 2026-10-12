package response

type Trade struct {
	Ticket        int64  `json:"ticket"`
	Login         int64  `json:"login"`
	Symbol        string `json:"symbol"`
	CmdDesc       string `json:"cmd_desc"`
	Lots          string `json:"lots"`
	OpenPrice     string `json:"open_price"`
	OpenTime      string `json:"open_time"`
	ClosePrice    string `json:"close_price"`
	CloseTime     string `json:"close_time"`
	SL            string `json:"sl"`
	TP            string `json:"tp"`
	Profit        string `json:"profit"`
	Swaps         string `json:"swaps"`
	Commission    string `json:"commission"`
	PL            string `json:"pl"`
	OpenTypeDesc  string `json:"open_type_desc"`
	CloseTypeDesc string `json:"close_type_desc"`
	OpenCloseTime string `json:"-"`
}

type ReportTrade struct {
	Ticket        int64  `json:"ticket"`
	Login         int64  `json:"login"`
	Passport      int64  `json:"passport"`
	Book          string `json:"book"`
	Symbol        string `json:"symbol"`
	CmdDesc       string `json:"cmd_desc"`
	Lots          string `json:"lots"`
	FloatingPL    string `json:"floating_pl"`
	ClosedPL      string `json:"closed_pl"`
	OpenPrice     string `json:"open_price"`
	OpenTime      string `json:"open_time"`
	ClosePrice    string `json:"close_price"`
	CloseTime     string `json:"close_time"`
	SL            string `json:"sl"`
	TP            string `json:"tp"`
	OpenTypeDesc  string `json:"open_type_desc"`
	CloseTypeDesc string `json:"close_type_desc"`
}
