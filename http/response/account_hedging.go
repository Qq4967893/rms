package response

type Preview struct {
	Login   int64             `json:"login"`
	Group   string            `json:"group"`
	Book    string            `json:"book"`
	Balance string            `json:"balance"`
	Equity  string            `json:"equity"`
	GrossPL string            `json:"gross_pl"`
	Net     map[string]string `json:"net"`
}

type ConfirmExposure struct {
	Currency     string `json:"currency"`
	BeforeNetUSD string `json:"before_net_usd"`
	AfterNetUSD  string `json:"after_net_usd"`
}

type ConfirmSummary struct {
	Symbol    string `json:"symbol"`
	BeforeNet string `json:"before_net"`
	AfterNet  string `json:"after_net"`
}

type Confirm struct {
	Summaries map[string][]ConfirmSummary
	Exposures map[string][]ConfirmExposure
}
