package service

import (
	"backend-rms/grpc/tsp"
	"backend-rms/model"
	"backend-rms/repo"
	"backend-rms/util"
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type symbolService struct{}

type SymbolHoldingAccount struct {
	Login            int64           `json:"login"`
	BuyLots          decimal.Decimal `json:"buy_lots"`
	SellLots         decimal.Decimal `json:"sell_lots"`
	NetLots          decimal.Decimal `json:"net_lots"`
	SymbolFloatingPL decimal.Decimal `json:"symbol_floating_pl"`
	FloatingPL       decimal.Decimal `json:"floating_pl"`
}

type SymbolTrade struct {
	Ticket      int64           `json:"ticket"`
	Login       int64           `json:"login"`
	OpenTime    time.Time       `json:"open_time"`
	OpenPrice   decimal.Decimal `json:"open_price"`
	SymbolPrice decimal.Decimal `json:"symbol_price"`
	Cmd         int8            `json:"cmd"`
	Lots        decimal.Decimal `json:"lots"`
	FloatingPL  decimal.Decimal `json:"floating_pl"`
	Profit      decimal.Decimal `json:"profit"`
	Swap        decimal.Decimal `json:"swaps"`
	Commission  decimal.Decimal `json:"commission"`
	SL          decimal.Decimal `json:"sl"`
	TP          decimal.Decimal `json:"tp"`
}

var SymbolService = &symbolService{}

func (ss *symbolService) GetList() ([]*model.Symbol, error) {
	symbols, err := repo.SymbolRepo.GetAll()
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}
	return symbols, nil
}

// 交易品种当前是否可以交易
func (ss *symbolService) IsPossibleTrade(symbol string) (bool, error) {
	symbolQuoteTimes, err := tsp.GetSymbolQuoteTime(context.Background(), "test request id", symbol)
	if err != nil {
		return false, errors.Wrap(err, util.ErrWrapMessage())
	}

	return ss.inTimeHorizon(symbolQuoteTimes[symbol].Data), nil
}

// 根据时间判断当前是否有报价是否可以交易
func (ss *symbolService) inTimeHorizon(weekTimes map[int32]string) bool {
	now := time.Now()
	weekday := int32(now.Weekday())
	dayQuoteTimes := strings.Split(weekTimes[weekday], ",")
	for _, dayQuoteTime := range dayQuoteTimes {
		hourMins := strings.Split(dayQuoteTime, "-")

		startHourMin := strings.Split(hourMins[0], ":")
		startHour, _ := strconv.Atoi(startHourMin[0])
		startMin, _ := strconv.Atoi(startHourMin[1])
		startDateTime := time.Date(now.Year(), now.Month(), now.Day(), startHour, startMin, 0, 0, now.Location())

		endHourMin := strings.Split(hourMins[1], ":")
		endHour, _ := strconv.Atoi(endHourMin[0])
		endMin, _ := strconv.Atoi(endHourMin[1])
		endDateTime := time.Date(now.Year(), now.Month(), now.Day(), endHour, endMin, 0, 0, now.Location())

		if now.Unix() > startDateTime.Unix() && now.Unix() < endDateTime.Unix() {
			return true
		}
	}
	return false
}

func (ss *symbolService) GetHoldingAccountList(symbol string, groups []string) (map[int64]SymbolHoldingAccount, error) {
	trades, err := repo.TradeRepo.GetOpeningTrades(groups, 0, symbol, -1)
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}
	if len(trades) == 0 {
		return nil, nil
	}

	tickets := make([]int64, 0)
	logins := make([]int64, 0)
	for _, trade := range trades {
		tickets = append(tickets, trade.Ticket)
		if !util.InSlice(trade.Login, logins) {
			logins = append(logins, trade.Login)
		}
	}
	floatingPLMap, err := CacheService.GetTradeFloatingPLMap(tickets)
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}

	accountFloatingPLMap, err := CacheService.GetAccountFloatingPLMap(logins)
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}

	list := make(map[int64]SymbolHoldingAccount)
	for _, trade := range trades {
		//获取订单的买卖手数
		buyLot := decimal.Zero
		sellLot := decimal.Zero
		if trade.Cmd == 0 {
			buyLot = decimal.NewFromInt(trade.Volume).Div(decimal.NewFromInt(100))
		} else if trade.Cmd == 1 {
			sellLot = decimal.NewFromInt(trade.Volume).Div(decimal.NewFromInt(100))
		}

		//获取订单的浮动盈亏
		symbolFloatingPL := floatingPLMap[trade.Ticket]

		if _, ok := list[trade.Login]; ok {
			middle := list[trade.Login]
			list[trade.Login] = SymbolHoldingAccount{
				Login:            trade.Login,
				BuyLots:          middle.BuyLots.Add(buyLot),
				SellLots:         middle.SellLots.Add(sellLot),
				NetLots:          middle.NetLots.Add(buyLot).Sub(sellLot),
				SymbolFloatingPL: middle.SymbolFloatingPL.Add(symbolFloatingPL),
				FloatingPL:       middle.FloatingPL,
			}
			continue
		}

		list[trade.Login] = SymbolHoldingAccount{
			Login:            trade.Login,
			BuyLots:          buyLot,
			SellLots:         sellLot,
			NetLots:          buyLot.Sub(sellLot),
			SymbolFloatingPL: symbolFloatingPL,
			FloatingPL:       accountFloatingPLMap[trade.Login],
		}
	}
	return list, nil
}

func (ss *symbolService) GetTradeList(symbol string, groups []string, cmd int8) ([]SymbolTrade, error) {
	trades, err := repo.TradeRepo.GetOpeningTrades(groups, 0, symbol, cmd)
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}
	if len(trades) == 0 {
		return nil, nil
	}

	tickets := make([]int64, 0)
	for _, trade := range trades {
		tickets = append(tickets, trade.Ticket)
	}
	floatingPLMap, err := CacheService.GetTradeFloatingPLMap(tickets)
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}

	quote, err := QuoteService.GetQuote(symbol)
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}

	list := make([]SymbolTrade, 0)
	for _, trade := range trades {
		var symbolPrice decimal.Decimal
		if trade.Cmd == 0 {
			symbolPrice = quote.Bid
		} else if trade.Cmd == 1 {
			symbolPrice = quote.Ask
		}

		floatingPL := decimal.Zero
		if _, ok := floatingPLMap[trade.Ticket]; ok {
			floatingPL = floatingPLMap[trade.Ticket]
		}
		middle := SymbolTrade{
			Ticket:      trade.Ticket,
			Login:       trade.Login,
			OpenTime:    trade.OpenTime,
			OpenPrice:   trade.OpenPrice,
			SymbolPrice: symbolPrice,
			Cmd:         trade.Cmd,
			Lots:        decimal.NewFromInt(trade.Volume).Div(decimal.NewFromInt(100)),
			FloatingPL:  floatingPL,
			Profit:      floatingPL.Sub(trade.Swaps).Sub(trade.Commission),
			Swap:        trade.Swaps,
			Commission:  trade.Commission,
			SL:          trade.SL,
			TP:          trade.TP,
		}
		list = append(list, middle)
	}
	return list, nil
}

func (ss *symbolService) GetPendingTradeList(symbol string, groups []string, cmd int8) ([]model.PendingTrade, error) {
	trades, err := repo.PendingTradeRepo.GetPendingTrades(groups, symbol, cmd)
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}
	if len(trades) == 0 {
		return nil, nil
	}
	return trades, nil
}
