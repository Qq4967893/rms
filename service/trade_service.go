package service

import (
	"backend-rms/model"
	"backend-rms/repo"
	"backend-rms/util"
	"time"

	carbon "github.com/golang-module/carbon/v2"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type tradeService struct{}

var TradeService = &tradeService{}

func (s *tradeService) GetTradeListByLogins(logins []int64, ticket int64, symbol string) (trades []*model.Trade, err error) {
	trades, err = repo.TradeRepo.GetTradeListByLogins(logins, ticket, symbol)
	return trades, errors.WithMessage(err, util.ErrWrapMessage())
}

func (s *tradeService) GetOpeningTrades(groups []string, login int64) (trades []model.OpeningTrade, err error) {
	trades, err = repo.TradeRepo.GetOpeningTrades(groups, login, "", -1)
	return trades, errors.WithMessage(err, util.ErrWrapMessage())
}

func (s *tradeService) GetPointOpeningTrades(groups []string, dateTime time.Time, symbol string) (trades []model.OpeningTrade, err error) {
	trades, err = s.GetPointOpeningTradeList(groups, []int64{}, dateTime, symbol)
	return trades, errors.WithMessage(err, util.ErrWrapMessage())
}

func (s *tradeService) GetPointOpeningTradeList(groups []string, logins []int64, dateTime time.Time, symbol string) (trades []model.OpeningTrade, err error) {
	startOfToday := carbon.Now().StartOfDay()
	if carbon.Time2Carbon(dateTime).Gte(startOfToday) {
		trades, err = repo.TradeRepo.GetPointOpeningTradeWithAccount(groups, logins, dateTime, symbol)
		if err != nil {
			return nil, errors.WithMessage(err, util.ErrWrapMessage())
		}
	} else {
		trades, err = repo.TradeRepo.GetPointOpeningTradeWithDaily(groups, logins, dateTime, symbol)
		if err != nil {
			return nil, errors.WithMessage(err, util.ErrWrapMessage())
		}
	}
	return
}

func (s *tradeService) GetReportTradeList(tradeStatus, symbol string, groups []string, startTime, endTime time.Time, login []int64, ticket int64, paging model.Paging) (trades []*model.ReportTrade, err error) {
	startTime, endTime = s.prepareDefaultTime(startTime, endTime)
	trades, err = repo.TradeRepo.GetReportTradeList(tradeStatus, symbol, groups, startTime, endTime, login, ticket, paging)
	return trades, errors.WithMessage(err, util.ErrWrapMessage())
}

func (s *tradeService) GetReportTradeCount(tradeStatus, symbol string, groups []string, startTime, endTime time.Time, login []int64, ticket int64) (totalCount int64, err error) {
	startTime, endTime = s.prepareDefaultTime(startTime, endTime)
	totalCount, err = repo.TradeRepo.GetReportTradeCount(tradeStatus, symbol, groups, startTime, endTime, login, ticket)
	return totalCount, errors.WithMessage(err, util.ErrWrapMessage())
}

func (s *tradeService) GetReportTradeListByTickets(tickets []int64, groups []string) (trades []*model.ReportTrade, err error) {
	trades, err = repo.TradeRepo.GetReportTradeListByTickets(tickets, groups)
	return trades, errors.WithMessage(err, util.ErrWrapMessage())
}

func (s *tradeService) GetClosedTrades(groups []string, startTime, endTime time.Time) (trades []model.ClosedTrade, err error) {
	trades, err = repo.TradeRepo.GetClosedTrades(groups, startTime, endTime)
	return trades, errors.WithMessage(err, util.ErrWrapMessage())
}

func (s *tradeService) GetClosedPL(groups []string, startTime, endTime time.Time) (closedPL decimal.Decimal, err error) {
	todayZero := carbon.Now().StartOfDay().ToDateTimeString()
	if carbon.Time2Carbon(endTime).Gte(carbon.Parse(todayZero)) {
		closedPL, err = repo.TradeRepo.GetClosedPLWithAccount(groups, startTime, endTime)
		if err != nil {
			return decimal.Zero, errors.WithMessage(err, util.ErrWrapMessage())
		}
	} else {
		closedPL, err = repo.TradeRepo.GetClosedPLWithDaily(groups, startTime, endTime)
		if err != nil {
			return decimal.Zero, errors.WithMessage(err, util.ErrWrapMessage())
		}
	}
	return
}

func (s *tradeService) GetAccountClosedPL(logins []int64, groups []string, startTime, endTime time.Time, symbol string) (map[int64]decimal.Decimal, error) {
	todayZero := carbon.Now().StartOfDay().ToDateTimeString()
	var accountClosedPL []model.AccountClosedPL
	var err error
	if carbon.Time2Carbon(endTime).Gte(carbon.Parse(todayZero)) {
		accountClosedPL, err = repo.TradeRepo.GetAccountClosedPLWithAccount(logins, groups, startTime, endTime, symbol)
		if err != nil {
			return nil, errors.WithMessage(err, util.ErrWrapMessage())
		}
	} else {
		accountClosedPL, err = repo.TradeRepo.GetAccountClosedPLWithDaily(logins, groups, startTime, endTime, symbol)
		if err != nil {
			return nil, errors.WithMessage(err, util.ErrWrapMessage())
		}
	}

	AccountClosedPLMap := map[int64]decimal.Decimal{}
	for _, v := range accountClosedPL {
		AccountClosedPLMap[v.Login] = v.ClosedPL
	}
	return AccountClosedPLMap, nil
}

// 获取账号的平仓盈亏（不关联组）
func (s *tradeService) GetAccountClosedPLWithoutGroup(logins []int64) ([]model.AccountClosedPL, error) {
	return repo.TradeRepo.GetAccountClosedPL(logins)
}

func (s *tradeService) GetSymbolClosedPL(groups []string, startTime, endTime time.Time) (map[string]map[string]decimal.Decimal, error) {
	todayZero := carbon.Now().StartOfDay().ToDateTimeString()
	var symbolClosedPL []model.SymbolClosedPL
	var err error
	if carbon.Time2Carbon(endTime).Gte(carbon.Parse(todayZero)) {
		symbolClosedPL, err = repo.TradeRepo.GetSymbolClosedPLWithAccount(groups, startTime, endTime)
		if err != nil {
			return nil, errors.WithMessage(err, util.ErrWrapMessage())
		}
	} else {
		symbolClosedPL, err = repo.TradeRepo.GetSymbolClosedPLWithDaily(groups, startTime, endTime)
		if err != nil {
			return nil, errors.WithMessage(err, util.ErrWrapMessage())
		}
	}

	result := map[string]map[string]decimal.Decimal{}
	SymbolLeverage, err := SymbolLeverageService.GetAllMap()
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}
	symbols, err := SymbolOriginalService.GetAllMap(false)
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}
	for k, _ := range symbols {
		closedPL := map[string]decimal.Decimal{
			"profit":     decimal.Zero,
			"swaps":      decimal.Zero,
			"commission": decimal.Zero,
			"closedPL":   decimal.Zero,
		}
		for _, v := range symbolClosedPL {
			if SymbolLeverage[v.Symbol] == k {
				closedPL["profit"] = closedPL["profit"].Add(v.Profit)
				closedPL["swaps"] = closedPL["swaps"].Add(v.Swaps)
				closedPL["commission"] = closedPL["commission"].Add(v.Commission)
				closedPL["closedPL"] = closedPL["closedPL"].Add(v.ClosedPL())
			}
		}
		result[k] = closedPL
	}
	return result, nil
}

func (s *tradeService) GetClosedPLMap(logins []int64, startTime, endTime string) (map[int64]decimal.Decimal, error) {
	closedPLMap, err := repo.TradeRepo.GetClosedPLMap(logins, startTime, endTime)
	return closedPLMap, errors.WithMessage(err, util.ErrWrapMessage())
}

func (s *tradeService) GetCloseTypeDesc(tradeType int8) string {
	switch tradeType {
	case 10:
		return "S / L"
	case 11:
		return "T / P"
	case 12:
		return "System Closed"
	default:
		return "Manual Close"
	}
}

func (s *tradeService) GetOpenTypeDesc(openType int8) string {
	switch openType {
	case 1:
		return "Market order"
	case 2:
		return "Pending order"
	default:
		return ""
	}
}

func (s *tradeService) GetCmdDesc(cmd int8) string {
	switch cmd {
	case 0:
		return "Buy"
	case 1:
		return "Sell"
	default:
		return ""
	}
}

func (s *tradeService) prepareDefaultTime(start, end time.Time) (startTime, endTime time.Time) {
	startTime = start
	if start.IsZero() {
		startTime = carbon.Parse(model.DEFAULT_CLOSE_TIME).Carbon2Time()
	}
	endTime = end
	if end.IsZero() {
		endTime = carbon.Now().Carbon2Time()
	}
	return
}
