package response

type ReportAccount struct {
	Login   int64  `json:"login"`
	Balance string `json:"balance"`
	Equity  string `json:"equity"`
}
