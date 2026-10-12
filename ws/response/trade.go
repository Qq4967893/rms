package response

type ReportTrade struct {
	Ticket        int64  `json:"ticket"`
	FloatingPL    string `json:"floating_pl"`
	ClosedPL      string `json:"closed_pl"`
	ClosePrice    string `json:"close_price"`
	CloseTime     string `json:"close_time"`
	CloseTypeDesc string `json:"close_type_desc"`
}
