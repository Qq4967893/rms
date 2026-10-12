package response

type Account struct {
	Login       int64  `json:"login"`
	Group       string `json:"group"`
	Book        string `json:"book"`
	RegTime     string `json:"reg_time"`
	Balance     string `json:"balance"`
	Equity      string `json:"equity"`
	MarginLevel string `json:"margin_level"`
	FloatingPL  string `json:"floating_pl"`
	ClosedPL    string `json:"closed_pl"`
	GrossPL     string `json:"gross_pl"`
	CashInflow  string `json:"cash_inflow"`
	CashOutflow string `json:"cash_outflow"`
	NetCashFlow string `json:"net_cash_flow"`
	IsMaster    string `json:"is_master"`
	ToGroups    string `json:"to_groups"`
}

type ReportAccount struct {
	Login               int64  `json:"login"`
	Book                string `json:"book"`
	Passport            int64  `json:"passport"`
	CashInflow          string `json:"cash_inflow"`
	CashOutflow         string `json:"cash_outflow"`
	NetCashFlow         string `json:"net_cash_flow"`
	Balance             string `json:"balance"`
	StartTimeBalance    string `json:"start_time_balance"`
	EndTimeBalance      string `json:"end_time_balance"`
	Equity              string `json:"equity"`
	StartTimeEquity     string `json:"start_time_equity"`
	EndTimeEquity       string `json:"end_time_equity"`
	StartTimeFloatingPL string `json:"start_time_floating_pl"`
	EndTimeFloatingPL   string `json:"end_time_floating_pl"`
	FloatingPL          string `json:"floating_pl"`
	ClosedPL            string `json:"closed_pl"`
	GrossPL             string `json:"gross_pl"`
	PeriodClosedPL      string `json:"period_closed_pl"`
	NetPeriodPL         string `json:"net_period_pl"`
}
