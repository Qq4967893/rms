package service

import (
	pb "backend-rms/grpc/et6api"
	"backend-rms/grpc/tsp"
	"backend-rms/http/response"
	"backend-rms/model"
	"backend-rms/util"
	"backend-rms/util/formatter"
	"context"
	"sort"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type passportDetailService struct{}

var PassportDetailService = &passportDetailService{}

func (s *passportDetailService) GetInfo(passport int64) (map[string]any, error) {
	user, err := UserService.GetInfoByPassport(passport)
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}

	if user == nil || user.ID == 0 {
		return nil, errors.New("passport not found")
	}

	managerName, err := AdminUserService.GetUsernameByID(user.AffiliationID)
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}

	deposit, err := WalletRecordService.GetDepositByPassport(passport)
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}

	withdrawal, err := WalletRecordService.GetWithdrawalByPassport(passport)
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}

	balance, err := WalletUserService.GetBalanceByPassport(passport)
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}

	ret := map[string]any{
		"passport": response.Passport{
			Passport:    user.Passport,
			Username:    user.Username,
			RegTime:     formatter.Datetime(user.CreatedAt),
			ManagerName: managerName,
			Deposit:     formatter.Amount(deposit),
			Withdrawal:  formatter.Amount(withdrawal),
			NetDeposit:  formatter.Amount(deposit.Add(withdrawal)),
			Balance:     formatter.Amount(balance),
		},
	}

	return ret, nil
}

func (s *passportDetailService) GetAccountListByLogin(requestID string, login int64) (map[string]any, error) {
	passport, err := TradeAccountService.GetPassportByLogin(login)
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}

	logins, err := TradeAccountService.GetTradeAccountIDs(passport)
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}

	loginTypeMap, err := TradeAccountService.GetLoginTypeMap(logins)
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}

	cashInflowMap, err := CashFlowService.GetCashInflowMap(logins, "", "")
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}

	cashOutflowMap, err := CashFlowService.GetCashOutflowMap(logins, "", "")
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}

	closedPLMap, err := TradeService.GetClosedPLMap(logins, "", "")
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}

	equityInfos, err := tsp.GetUserEquity(context.Background(), requestID, logins)
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}

	equityMap := make(map[int64]*pb.EquityInfo, len(equityInfos))
	for _, equityInfo := range equityInfos {
		equityMap[int64(equityInfo.Login)] = equityInfo
	}

	groupMap, err := GroupService.GetAllMap()
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}

	var groups []string
	for group, _ := range groupMap {
		groups = append(groups, group)
	}

	groupPairs, err := GroupPairService.GetToGroups()
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}

	accountList, err := AccountService.GetListByLogins(logins, groups)
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}

	var accounts []*response.Account
	for _, account := range accountList {
		if account.Login == login {
			equityInfo := equityMap[account.Login]
			accounts = append(accounts, s.buildAccount(account, loginTypeMap, groupMap, equityInfo, cashInflowMap, cashOutflowMap, closedPLMap, groupPairs))
			break
		}
	}

	for _, account := range accountList {
		if account.Login == login {
			continue
		}
		equityInfo := equityMap[account.Login]
		accounts = append(accounts, s.buildAccount(account, loginTypeMap, groupMap, equityInfo, cashInflowMap, cashOutflowMap, closedPLMap, groupPairs))
	}

	return map[string]any{
		"accounts": accounts,
	}, nil
}

func (s *passportDetailService) GetTradeListByPassport(requestID string, passport, login, ticket int64, symbol string, paging *model.Paging) (map[string]any, error) {
	var err error
	logins := []int64{login}
	if login <= 0 {
		logins, err = TradeAccountService.GetTradeAccountIDs(passport)
		if err != nil {
			return nil, errors.WithMessage(err, util.ErrWrapMessage())
		}
	}

	tradeList, err := TradeService.GetTradeListByLogins(logins, ticket, symbol)
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}

	var tickets []int64
	openingTrades := make(map[int64][]int64)
	for _, trade := range tradeList {
		tickets = append(tickets, trade.Ticket)
		if trade.CloseTime.IsZero() {
			openingTrades[trade.Login] = append(openingTrades[trade.Login], trade.Ticket)
		}
	}

	openTypeMap, err := EventService.GetOpenTypeMapByTickets(tickets)
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}

	openingTradeMap := make(map[int64]*pb.Order)
	for login, tickets := range openingTrades {
		res, err := tsp.GetTradeList(context.Background(), requestID, login, tickets)
		if err != nil {
			return nil, errors.WithMessage(err, util.ErrWrapMessage())
		}

		for _, o := range res.GetOrder() {
			openingTradeMap[o.Ticket] = o
		}
	}

	var (
		totalVolume     int64
		totalFloatingPL decimal.Decimal
		totalClosedPL   decimal.Decimal
	)

	var trades []response.Trade
	for _, trade := range tradeList {
		profit := trade.Profit
		openCloseTime := trade.OpenTime
		closeTime := ""
		closePrice := ""
		closeTypeDesc := ""
		if trade.CloseTime.IsZero() {
			profit = decimal.NewFromFloat(openingTradeMap[trade.Ticket].Profit)
			totalFloatingPL = totalFloatingPL.Add(profit).Add(trade.Swaps).Add(trade.Commission)
		} else {
			closeTime = formatter.Datetime(trade.CloseTime)
			closePrice = trade.ClosePrice.String()
			openCloseTime = trade.CloseTime
			closeTypeDesc = TradeService.GetCloseTypeDesc(trade.Type)
			totalClosedPL = totalClosedPL.Add(trade.Profit).Add(trade.Swaps).Add(trade.Commission)
		}

		totalVolume += trade.Volume

		trades = append(trades, response.Trade{
			Ticket:        trade.Ticket,
			Login:         trade.Login,
			Symbol:        trade.Symbol,
			CmdDesc:       TradeService.GetCmdDesc(trade.Cmd),
			Lots:          decimal.NewFromFloat(float64(trade.Volume) / 100).StringFixed(2),
			OpenPrice:     trade.OpenPrice.String(),
			OpenTime:      formatter.Datetime(trade.OpenTime),
			ClosePrice:    closePrice,
			CloseTime:     closeTime,
			SL:            trade.SL.String(),
			TP:            trade.TP.String(),
			Profit:        formatter.Amount(profit),
			Swaps:         formatter.Amount(trade.Swaps),
			Commission:    formatter.Amount(trade.Commission),
			PL:            formatter.Amount(profit.Add(trade.Swaps).Add(trade.Commission)),
			OpenTypeDesc:  TradeService.GetOpenTypeDesc(openTypeMap[trade.Ticket]),
			CloseTypeDesc: closeTypeDesc,
			OpenCloseTime: formatter.Datetime(openCloseTime),
		})
	}

	sort.Slice(trades, func(i, j int) bool {
		switch paging.Sort {
		case "open_close_time":
			if paging.Order == "asc" {
				return trades[i].OpenCloseTime < trades[j].OpenCloseTime
			} else {
				return trades[i].OpenCloseTime > trades[j].OpenCloseTime
			}
		case "open_time":
			if paging.Order == "asc" {
				return trades[i].OpenTime < trades[j].OpenTime
			} else {
				return trades[i].OpenTime > trades[j].OpenTime
			}
		case "close_time":
			if paging.Order == "asc" {
				return trades[i].CloseTime < trades[j].CloseTime
			} else {
				return trades[i].CloseTime > trades[j].CloseTime
			}
		case "open_price":
			if paging.Order == "asc" {
				return trades[i].OpenPrice < trades[j].OpenPrice
			} else {
				return trades[i].OpenPrice > trades[j].OpenPrice
			}
		default:
			return false
		}
	})

	start := paging.Offset
	if start > len(trades) {
		start = len(trades)
	}

	end := paging.Offset + paging.Limit
	if end > len(trades) {
		end = len(trades)
	}

	return map[string]any{
		"summary": response.TradeSummary{
			TradeCount: len(trades),
			Lots:       decimal.NewFromFloat(float64(totalVolume) / 100).StringFixed(2),
			FloatingPL: formatter.Amount(totalFloatingPL),
			ClosedPL:   formatter.Amount(totalClosedPL),
		},
		"trades":      trades[start:end],
		"total_count": len(trades),
	}, nil
}

func (s *passportDetailService) GetSymbolListByPassport(requestID string, passport, login int64) (map[string]any, error) {
	var err error
	logins := []int64{login}
	if login <= 0 {
		logins, err = TradeAccountService.GetTradeAccountIDs(passport)
		if err != nil {
			return nil, errors.WithMessage(err, util.ErrWrapMessage())
		}
	}

	symbolMap, err := SymbolLeverageService.GetAllMap()
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}

	tradeList, err := TradeService.GetTradeListByLogins(logins, 0, "")
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}

	type Volume struct {
		Volume     int64
		BuyVolume  int64
		SellVolume int64
	}

	symbolVolumeMap := make(map[string]*Volume)
	for _, trade := range tradeList {
		symbol := symbolMap[trade.Symbol]
		if _, ok := symbolVolumeMap[symbol]; !ok {
			symbolVolumeMap[symbol] = &Volume{}
		}
		symbolVolumeMap[symbol].Volume += trade.Volume
		if trade.Cmd == 0 {
			symbolVolumeMap[symbol].BuyVolume += trade.Volume
		} else {
			symbolVolumeMap[symbol].SellVolume += trade.Volume
		}
	}

	var symbols []*response.SymbolLots
	for symbol, v := range symbolVolumeMap {
		symbols = append(symbols, &response.SymbolLots{
			Symbol:   symbol,
			Lots:     decimal.NewFromFloat(float64(v.Volume) / 100).StringFixed(2),
			BuyLots:  decimal.NewFromFloat(float64(v.BuyVolume) / 100).StringFixed(2),
			SellLots: decimal.NewFromFloat(float64(v.SellVolume) / 100).StringFixed(2),
			NetLots:  decimal.NewFromFloat(float64(v.BuyVolume-v.SellVolume) / 100).StringFixed(2),
		})
	}

	return map[string]any{
		"symbols": symbols,
	}, nil
}

func (s *passportDetailService) buildAccount(account *model.Account, loginTypeMap map[int64]int8, groupMap map[string]string,
	equityInfo *pb.EquityInfo, cashInflowMap, cashOutflowMap, closedPLMap map[int64]decimal.Decimal, groupPairs map[string]string) *response.Account {
	floatingPL := decimal.NewFromFloat(equityInfo.Equity).Sub(decimal.NewFromFloat(equityInfo.Balance))

	toGroups := ""
	if _, ok := groupPairs[account.Group]; ok {
		toGroups = groupPairs[account.Group]
	}

	return &response.Account{
		Login:       account.Login,
		Group:       account.Group,
		Book:        groupMap[account.Group],
		RegTime:     formatter.Datetime(account.Regdate),
		Balance:     formatter.Amount(decimal.NewFromFloat(equityInfo.Balance)),
		Equity:      formatter.Amount(decimal.NewFromFloat(equityInfo.Equity)),
		MarginLevel: decimal.NewFromFloat(equityInfo.MarginLevel).Round(2).String(),
		FloatingPL:  formatter.Amount(floatingPL),
		ClosedPL:    formatter.Amount(closedPLMap[account.Login]),
		GrossPL:     formatter.Amount(floatingPL.Add(closedPLMap[account.Login])),
		CashInflow:  formatter.Amount(cashInflowMap[account.Login]),
		CashOutflow: formatter.Amount(cashOutflowMap[account.Login]),
		NetCashFlow: formatter.Amount(cashInflowMap[account.Login].Add(cashOutflowMap[account.Login])),
		IsMaster:    formatter.Bool(loginTypeMap[account.Login] == 3),
		ToGroups:    toGroups,
	}
}
