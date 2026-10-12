package service

import (
	"backend-rms/model"
	"backend-rms/util"
	"sort"
	"time"

	"github.com/golang-module/carbon/v2"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"golang.org/x/sync/errgroup"
)

type ReportAccount struct {
	Login               int64
	Group               string
	Passport            int64
	CashInflow          decimal.Decimal
	CashOutflow         decimal.Decimal
	NetCashFlow         decimal.Decimal
	Balance             decimal.Decimal
	StartTimeBalance    decimal.Decimal
	EndTimeBalance      decimal.Decimal
	Equity              decimal.Decimal
	StartTimeEquity     decimal.Decimal
	EndTimeEquity       decimal.Decimal
	StartTimeFloatingPL decimal.Decimal
	EndTimeFloatingPL   decimal.Decimal
	FloatingPL          decimal.Decimal
	ClosedPL            decimal.Decimal
	GrossPL             decimal.Decimal
	PeriodClosedPL      decimal.Decimal
	NetPeriodPL         decimal.Decimal
}

func (a *ReportAccount) GetNetCashFlow() decimal.Decimal {
	return a.CashInflow.Add(a.CashOutflow)
}

func (a *ReportAccount) GetEquity() decimal.Decimal {
	return a.Balance.Add(a.FloatingPL)
}

func (a *ReportAccount) GetStartTimeEquity() decimal.Decimal {
	return a.StartTimeBalance.Add(a.StartTimeFloatingPL)
}

func (a *ReportAccount) GetEndTimeEquity() decimal.Decimal {
	return a.EndTimeBalance.Add(a.EndTimeFloatingPL)
}

func (a *ReportAccount) GetGrossPL() decimal.Decimal {
	return a.FloatingPL.Add(a.ClosedPL)
}

func (a *ReportAccount) GetNetPeriodPL() decimal.Decimal {
	return a.EndTimeFloatingPL.Sub(a.StartTimeFloatingPL).Add(a.PeriodClosedPL)
}

type reportAccountService struct{}

var ReportAccountService = &reportAccountService{}

func (ras *reportAccountService) GetReportAccountList(requestID string, groups []string, login []int64, startTime, endTime time.Time, paging model.Paging) ([]*ReportAccount, error) {
	var accounts []*ReportAccount
	filledFieldMap := make(map[string]struct{})
	if paging.Sort == "reg_time" || len(login) > 0 {
		logins := make([]int64, 0)
		if len(login) > 0 {
			logins = login
		}
		paging.Sort = "regdate"
		accountList, err := AccountService.GetReportAccountList(groups, logins, &paging)
		if err != nil {
			return nil, errors.WithMessage(err, util.ErrWrapMessage())
		}

		for _, account := range accountList {
			accounts = append(accounts, &ReportAccount{
				Login:   account.Login,
				Group:   account.Group,
				Balance: account.Balance,
			})
		}
	} else {
		accountList, err := AccountService.GetAllReportAccounts(groups)
		if err != nil {
			return nil, errors.WithMessage(err, util.ErrWrapMessage())
		}

		var logins []int64
		for _, account := range accountList {
			logins = append(logins, account.Login)
			accounts = append(accounts, &ReportAccount{
				Login:   account.Login,
				Group:   account.Group,
				Balance: account.Balance,
			})
		}

		if paging.Sort == "cash_inflow" || paging.Sort == "net_cash_flow" {
			filledFieldMap["cash_inflow"] = struct{}{}
			if err := ras.fillCashAmount(accounts, logins, startTime, endTime, true); err != nil {
				return nil, errors.WithMessage(err, util.ErrWrapMessage())
			}
		}
		if paging.Sort == "cash_outflow" || paging.Sort == "net_cash_flow" {
			filledFieldMap["cash_outflow"] = struct{}{}
			if err := ras.fillCashAmount(accounts, logins, startTime, endTime, false); err != nil {
				return nil, errors.WithMessage(err, util.ErrWrapMessage())
			}
		}
		if paging.Sort == "start_time_balance" || paging.Sort == "start_time_equity" {
			filledFieldMap["start_time_balance"] = struct{}{}
			if err := ras.fillPointBalance(accounts, logins, startTime, true); err != nil {
				return nil, errors.WithMessage(err, util.ErrWrapMessage())
			}
		}
		if paging.Sort == "end_time_balance" || paging.Sort == "end_time_equity" {
			filledFieldMap["end_time_balance"] = struct{}{}
			if err := ras.fillPointBalance(accounts, logins, endTime, false); err != nil {
				return nil, errors.WithMessage(err, util.ErrWrapMessage())
			}
		}
		if paging.Sort == "start_time_floating_pl" || paging.Sort == "start_time_equity" || paging.Sort == "net_period_pl" {
			filledFieldMap["start_time_floating_pl"] = struct{}{}
			if err := ras.fillPointFloatingPL(requestID, accounts, groups, logins, startTime, true); err != nil {
				return nil, errors.WithMessage(err, util.ErrWrapMessage())
			}
		}
		if paging.Sort == "end_time_floating_pl" || paging.Sort == "end_time_equity" || paging.Sort == "net_period_pl" {
			filledFieldMap["end_time_floating_pl"] = struct{}{}
			if err := ras.fillPointFloatingPL(requestID, accounts, groups, logins, endTime, false); err != nil {
				return nil, errors.WithMessage(err, util.ErrWrapMessage())
			}
		}
		if paging.Sort == "floating_pl" || paging.Sort == "gross_pl" {
			filledFieldMap["floating_pl"] = struct{}{}
			if err := ras.fillFloatingPL(accounts, []int64{}); err != nil {
				return nil, errors.WithMessage(err, util.ErrWrapMessage())
			}
		}
		if paging.Sort == "closed_pl" || paging.Sort == "gross_pl" {
			filledFieldMap["closed_pl"] = struct{}{}
			if err := ras.fillClosedPL(accounts, logins, time.Time{}, time.Now(), false); err != nil {
				return nil, errors.WithMessage(err, util.ErrWrapMessage())
			}
		}
		if paging.Sort == "period_closed_pl" || paging.Sort == "net_period_pl" {
			filledFieldMap["period_closed_pl"] = struct{}{}
			if err := ras.fillClosedPL(accounts, logins, startTime, endTime, true); err != nil {
				return nil, errors.WithMessage(err, util.ErrWrapMessage())
			}
		}

		switch paging.Sort {
		case "net_cash_flow":
			filledFieldMap["net_cash_flow"] = struct{}{}
			ras.fillNetCashFlow(accounts)
		case "start_time_equity":
			filledFieldMap["start_time_equity"] = struct{}{}
			ras.fillStartTimeEquity(accounts)
		case "end_time_equity":
			filledFieldMap["end_time_equity"] = struct{}{}
			ras.fillEndTimeEquity(accounts)
		case "gross_pl":
			filledFieldMap["gross_pl"] = struct{}{}
			ras.fillGrossPL(accounts)
		case "net_period_pl":
			filledFieldMap["net_period_pl"] = struct{}{}
			ras.fillNetPeriodPL(accounts)
		}

		sort.Slice(accounts, func(i, j int) bool {
			switch paging.Sort {
			case "cash_inflow":
				if paging.Order == "asc" {
					return accounts[i].CashInflow.Cmp(accounts[j].CashInflow) < 0
				} else {
					return accounts[i].CashInflow.Cmp(accounts[j].CashInflow) > 0
				}
			case "cash_outflow":
				if paging.Order == "asc" {
					return accounts[i].CashOutflow.Cmp(accounts[j].CashOutflow) < 0
				} else {
					return accounts[i].CashOutflow.Cmp(accounts[j].CashOutflow) > 0
				}
			case "net_cash_flow":
				if paging.Order == "asc" {
					return accounts[i].NetCashFlow.Cmp(accounts[j].NetCashFlow) < 0
				} else {
					return accounts[i].NetCashFlow.Cmp(accounts[j].NetCashFlow) > 0
				}
			case "start_time_balance":
				if paging.Order == "asc" {
					return accounts[i].StartTimeBalance.Cmp(accounts[j].StartTimeBalance) < 0
				} else {
					return accounts[i].StartTimeBalance.Cmp(accounts[j].StartTimeBalance) > 0
				}
			case "end_time_balance":
				if paging.Order == "asc" {
					return accounts[i].EndTimeBalance.Cmp(accounts[j].EndTimeBalance) < 0
				} else {
					return accounts[i].EndTimeBalance.Cmp(accounts[j].EndTimeBalance) > 0
				}
			case "start_time_equity":
				if paging.Order == "asc" {
					return accounts[i].StartTimeEquity.Cmp(accounts[j].StartTimeEquity) < 0
				} else {
					return accounts[i].StartTimeEquity.Cmp(accounts[j].StartTimeEquity) > 0
				}
			case "end_time_equity":
				if paging.Order == "asc" {
					return accounts[i].EndTimeEquity.Cmp(accounts[j].EndTimeEquity) < 0
				} else {
					return accounts[i].EndTimeEquity.Cmp(accounts[j].EndTimeEquity) > 0
				}
			case "start_time_floating_pl":
				if paging.Order == "asc" {
					return accounts[i].StartTimeFloatingPL.Cmp(accounts[j].StartTimeFloatingPL) < 0
				} else {
					return accounts[i].StartTimeFloatingPL.Cmp(accounts[j].StartTimeFloatingPL) > 0
				}
			case "end_time_floating_pl":
				if paging.Order == "asc" {
					return accounts[i].EndTimeFloatingPL.Cmp(accounts[j].EndTimeFloatingPL) < 0
				} else {
					return accounts[i].EndTimeFloatingPL.Cmp(accounts[j].EndTimeFloatingPL) > 0
				}
			case "floating_pl":
				if paging.Order == "asc" {
					return accounts[i].FloatingPL.Cmp(accounts[j].FloatingPL) < 0
				} else {
					return accounts[i].FloatingPL.Cmp(accounts[j].FloatingPL) > 0
				}
			case "closed_pl":
				if paging.Order == "asc" {
					return accounts[i].ClosedPL.Cmp(accounts[j].ClosedPL) < 0
				} else {
					return accounts[i].ClosedPL.Cmp(accounts[j].ClosedPL) > 0
				}
			case "gross_pl":
				if paging.Order == "asc" {
					return accounts[i].GrossPL.Cmp(accounts[j].GrossPL) < 0
				} else {
					return accounts[i].GrossPL.Cmp(accounts[j].GrossPL) > 0
				}
			case "period_closed_pl":
				if paging.Order == "asc" {
					return accounts[i].PeriodClosedPL.Cmp(accounts[j].PeriodClosedPL) < 0
				} else {
					return accounts[i].PeriodClosedPL.Cmp(accounts[j].PeriodClosedPL) > 0
				}
			case "net_period_pl":
				if paging.Order == "asc" {
					return accounts[i].NetPeriodPL.Cmp(accounts[j].NetPeriodPL) < 0
				} else {
					return accounts[i].NetPeriodPL.Cmp(accounts[j].NetPeriodPL) > 0
				}
			default:
				return false
			}
		})

		start := paging.Offset
		if start > len(accounts) {
			start = len(accounts)
		}

		end := paging.Offset + paging.Limit
		if end > len(accounts) {
			end = len(accounts)
		}

		accounts = accounts[start:end]
	}

	var logins []int64
	for _, account := range accounts {
		logins = append(logins, account.Login)
	}

	fields := []string{
		"passport",
		"cash_inflow",
		"cash_outflow",
		"start_time_balance",
		"end_time_balance",
		"start_time_floating_pl",
		"end_time_floating_pl",
		"floating_pl",
		"closed_pl",
		"period_closed_pl",
	}

	g := new(errgroup.Group)
	for _, field := range fields {
		if _, ok := filledFieldMap[field]; ok {
			continue
		}
		switch field {
		case "passport":
			g.Go(func() error {
				return ras.fillPassport(accounts, logins)
			})
		case "cash_inflow":
			g.Go(func() error {
				return ras.fillCashAmount(accounts, logins, startTime, endTime, true)
			})
		case "cash_outflow":
			g.Go(func() error {
				return ras.fillCashAmount(accounts, logins, startTime, endTime, false)
			})
		case "start_time_balance":
			g.Go(func() error {
				return ras.fillPointBalance(accounts, logins, startTime, true)
			})
		case "end_time_balance":
			g.Go(func() error {
				return ras.fillPointBalance(accounts, logins, endTime, false)
			})
		case "start_time_floating_pl":
			g.Go(func() error {
				return ras.fillPointFloatingPL(requestID, accounts, groups, logins, startTime, true)
			})
		case "end_time_floating_pl":
			g.Go(func() error {
				return ras.fillPointFloatingPL(requestID, accounts, groups, logins, endTime, false)
			})
		case "floating_pl":
			g.Go(func() error {
				return ras.fillFloatingPL(accounts, logins)
			})
		case "closed_pl":
			g.Go(func() error {
				return ras.fillClosedPL(accounts, logins, time.Time{}, time.Now(), false)
			})
		case "period_closed_pl":
			g.Go(func() error {
				return ras.fillClosedPL(accounts, logins, startTime, endTime, true)
			})
		}
	}

	if err := g.Wait(); err != nil {
		return nil, errors.Wrap(err, "g.Wait() is failed")
	}

	return accounts, nil
}

func (ras *reportAccountService) fillPassport(accounts []*ReportAccount, logins []int64) error {
	loginPassportMap, err := TradeAccountService.GetLoginPassportMap(logins)
	if err != nil {
		return errors.WithMessage(err, util.ErrWrapMessage())
	}
	for _, account := range accounts {
		account.Passport = loginPassportMap[account.Login]
	}
	return nil
}

func (ras *reportAccountService) fillCashAmount(accounts []*ReportAccount, logins []int64, startTime, endTime time.Time, isInflow bool) (err error) {
	startTimeString := carbon.Time2Carbon(startTime).ToDateTimeString()
	endTimeString := carbon.Time2Carbon(endTime).ToDateTimeString()
	cashAmountMap := make(map[int64]decimal.Decimal, len(logins))
	chunkSize := 1000
	for i := 0; i < len(logins); i += chunkSize {
		end := i + chunkSize
		if end > len(logins) {
			end = len(logins)
		}
		var results map[int64]decimal.Decimal
		if isInflow {
			results, err = CashFlowService.GetCashInflowMap(logins[i:end], startTimeString, endTimeString)
			if err != nil {
				return errors.WithMessage(err, util.ErrWrapMessage())
			}
		} else {
			results, err = CashFlowService.GetCashOutflowMap(logins[i:end], startTimeString, endTimeString)
			if err != nil {
				return errors.WithMessage(err, util.ErrWrapMessage())
			}
		}
		for login, amount := range results {
			cashAmountMap[login] = amount
		}
	}
	for _, account := range accounts {
		if isInflow {
			account.CashInflow = cashAmountMap[account.Login]
		} else {
			account.CashOutflow = cashAmountMap[account.Login]
		}
	}
	return nil
}

func (ras *reportAccountService) fillPointBalance(accounts []*ReportAccount, logins []int64, dateTime time.Time, isStartTime bool) error {
	timePoint := carbon.Time2Carbon(dateTime)
	yesterday := timePoint.Yesterday().Carbon2Time()
	reportDailyTime, err := ReportDailyService.GetReportDailyTime(yesterday)
	if err != nil {
		return errors.WithMessage(err, util.ErrWrapMessage())
	}

	if reportDailyTime == nil {
		return errors.Errorf("time %s does not exist in report daily table", yesterday)
	}

	today := carbon.Time2Carbon(*reportDailyTime).AddSeconds(1).ToDateTimeString()

	var dailyBalanceMap, netCashFlowMap, closedPLMap map[int64]decimal.Decimal

	g := new(errgroup.Group)
	g.Go(func() error {
		dailyBalanceMap, err = ReportDailyService.GetLoginBalanceMapByTime(yesterday, logins)
		return errors.WithMessage(err, util.ErrWrapMessage())
	})

	g.Go(func() error {
		netCashFlowMap, err = CashFlowService.GetNetCashFlowMap(logins, today, timePoint.ToDateTimeString())
		return errors.WithMessage(err, util.ErrWrapMessage())
	})

	g.Go(func() error {
		closedPLMap, err = TradeService.GetClosedPLMap(logins, today, timePoint.ToDateTimeString())
		return errors.WithMessage(err, util.ErrWrapMessage())
	})

	if err = g.Wait(); err != nil {
		return errors.Wrap(err, "g.Wait is failed")
	}

	loginBalanceMap := make(map[int64]decimal.Decimal, len(logins))
	for _, login := range logins {
		loginBalanceMap[login] = dailyBalanceMap[login].Add(netCashFlowMap[login]).Add(closedPLMap[login])
	}

	for _, account := range accounts {
		if isStartTime {
			account.StartTimeBalance = loginBalanceMap[account.Login]
		} else {
			account.EndTimeBalance = loginBalanceMap[account.Login]
		}
	}

	return nil
}

func (ras *reportAccountService) fillClosedPL(accounts []*ReportAccount, logins []int64, startTime, endTime time.Time, isPeriod bool) (err error) {
	closedPLMap := make(map[int64]decimal.Decimal, len(logins))
	chunkSize := 1000
	for i := 0; i < len(logins); i += chunkSize {
		end := i + chunkSize
		if end > len(logins) {
			end = len(logins)
		}
		var results map[int64]decimal.Decimal
		results, err = TradeService.GetAccountClosedPL(logins[i:end], []string{}, startTime, endTime, "")
		if err != nil {
			return errors.WithMessage(err, util.ErrWrapMessage())
		}
		for login, closedPL := range results {
			closedPLMap[login] = closedPL
		}
	}
	for _, account := range accounts {
		if isPeriod {
			account.PeriodClosedPL = closedPLMap[account.Login]
		} else {
			account.ClosedPL = closedPLMap[account.Login]
		}
	}
	return nil
}

func (ras *reportAccountService) fillPointFloatingPL(requestID string, accounts []*ReportAccount, groups []string, logins []int64, dateTime time.Time, isStartTime bool) error {
	pointAccountPLMap, err := FloatingPLService.GetPointAccountPLMap(groups, logins, dateTime, requestID)
	if err != nil {
		return errors.WithMessage(err, util.ErrWrapMessage())
	}
	for _, account := range accounts {
		if isStartTime {
			account.StartTimeFloatingPL = pointAccountPLMap[account.Login]
		} else {
			account.EndTimeFloatingPL = pointAccountPLMap[account.Login]
		}
	}
	return nil
}

func (ras *reportAccountService) fillFloatingPL(accounts []*ReportAccount, logins []int64) (err error) {
	floatingPLMap := make(map[int64]decimal.Decimal)
	if len(logins) > 0 {
		floatingPLMap, err = CacheService.GetAccountFloatingPLMap(logins)
	} else {
		floatingPLMap, err = CacheService.GetAllAccountFloatingPLMap()
	}
	if err != nil {
		return errors.WithMessage(err, util.ErrWrapMessage())
	}
	for _, account := range accounts {
		account.FloatingPL = floatingPLMap[account.Login]
	}
	return nil
}

func (ras *reportAccountService) fillNetCashFlow(accounts []*ReportAccount) {
	for _, account := range accounts {
		account.NetCashFlow = account.GetNetCashFlow()
	}
}

func (ras *reportAccountService) fillStartTimeEquity(accounts []*ReportAccount) {
	for _, account := range accounts {
		account.StartTimeEquity = account.GetStartTimeEquity()
	}
}

func (ras *reportAccountService) fillEndTimeEquity(accounts []*ReportAccount) {
	for _, account := range accounts {
		account.EndTimeEquity = account.GetEndTimeEquity()
	}
}

func (ras *reportAccountService) fillGrossPL(accounts []*ReportAccount) {
	for _, account := range accounts {
		account.GrossPL = account.GetGrossPL()
	}
}

func (ras *reportAccountService) fillNetPeriodPL(accounts []*ReportAccount) {
	for _, account := range accounts {
		account.NetPeriodPL = account.GetNetPeriodPL()
	}
}
