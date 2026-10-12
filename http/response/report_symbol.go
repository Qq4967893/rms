package response

import (
	"github.com/shopspring/decimal"
)

type ReportSymbolSummaryResult struct {
	Date               string `json:"date"`
	Symbol             string `json:"symbol"`
	CashInflow         string `json:"cash_inflow"`
	CashOutflow        string `json:"cash_outflow"`
	NetCashFlow        string `json:"net_cash_flow"`
	EndTimeBalance     string `json:"end_time_balance"`
	EndTimeEquity      string `json:"end_time_equity"`
	ClosedPL           string `json:"closed_pl"`
	ClosedProfit       string `json:"closed_profit"`
	ClosedSwap         string `json:"closed_swap"`
	ClosedCommission   string `json:"closed_commission"`
	FloatingPL         string `json:"floating_pl"`
	FloatingProfit     string `json:"floating_profit"`
	FloatingSwap       string `json:"floating_swap"`
	FloatingCommission string `json:"floating_commission"`
	PeriodPL           string `json:"period_pl"`
}

type ReportSymbolSummary struct {
	Date               string          `json:"date"`
	Symbol             string          `json:"symbol"`
	CashInflow         decimal.Decimal `json:"cash_inflow"`
	CashOutflow        decimal.Decimal `json:"cash_outflow"`
	NetCashFlow        decimal.Decimal `json:"net_cash_flow"`
	EndTimeBalance     decimal.Decimal `json:"end_time_balance"`
	EndTimeEquity      decimal.Decimal `json:"end_time_equity"`
	ClosedPL           decimal.Decimal `json:"closed_pl"`
	ClosedProfit       decimal.Decimal `json:"closed_profit"`
	ClosedSwap         decimal.Decimal `json:"closed_swap"`
	ClosedCommission   decimal.Decimal `json:"closed_commission"`
	FloatingPL         decimal.Decimal `json:"floating_pl"`
	FloatingProfit     decimal.Decimal `json:"floating_profit"`
	FloatingSwap       decimal.Decimal `json:"floating_swap"`
	FloatingCommission decimal.Decimal `json:"floating_commission"`
	PeriodPL           decimal.Decimal `json:"period_pl"`
}

type ReportSymbolGetListResult struct {
	Date               string `json:"date"`
	Symbol             string `json:"symbol"`
	ClosedPL           string `json:"closed_pl"`
	ClosedProfit       string `json:"closed_profit"`
	ClosedSwap         string `json:"closed_swap"`
	ClosedCommission   string `json:"closed_commission"`
	FloatingPL         string `json:"floating_pl"`
	FloatingProfit     string `json:"floating_profit"`
	FloatingSwap       string `json:"floating_swap"`
	FloatingCommission string `json:"floating_commission"`
	PeriodPL           string `json:"period_pl"`
}

type ReportSymbolGetList struct {
	Date               string          `json:"date"`
	Symbol             string          `json:"symbol"`
	ClosedPL           decimal.Decimal `json:"closed_pl"`
	ClosedProfit       decimal.Decimal `json:"closed_profit"`
	ClosedSwap         decimal.Decimal `json:"closed_swap"`
	ClosedCommission   decimal.Decimal `json:"closed_commission"`
	FloatingPL         decimal.Decimal `json:"floating_pl"`
	FloatingProfit     decimal.Decimal `json:"floating_profit"`
	FloatingSwap       decimal.Decimal `json:"floating_swap"`
	FloatingCommission decimal.Decimal `json:"floating_commission"`
	PeriodPL           decimal.Decimal `json:"period_pl"`
}

type ReportSymbolGetSymbolAccountListResult struct {
	Account             int64  `json:"account"`
	PeriodPL            string `json:"period_pl"`
	StartTimeFloatingPL string `json:"start_time_floating_pl"`
	EndTimeFloatingPL   string `json:"end_time_floating_pl"`
	ClosedPL            string `json:"closed_pl"`
	PeriodFloatingPL    string `json:"period_floating_pl"`
}

type ReportSymbolGetSymbolAccountList struct {
	Account             int64           `json:"account"`
	PeriodPL            decimal.Decimal `json:"period_pl"`
	StartTimeFloatingPL decimal.Decimal `json:"start_time_floating_pl"`
	EndTimeFloatingPL   decimal.Decimal `json:"end_time_floating_pl"`
	ClosedPL            decimal.Decimal `json:"closed_pl"`
	PeriodFloatingPL    decimal.Decimal `json:"period_floating_pl"`
}
