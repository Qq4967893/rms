package service

import (
	"backend-rms/core/logger"
	"backend-rms/grpc/tsp"
	"backend-rms/model"
	"backend-rms/util"
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

type accountHedgingService struct{}

var AccountHedgingService = &accountHedgingService{}

type HedgingTrade struct {
	Symbol string
	Cmd    int8
	Volume int64
}

type AfterSplitTrade struct {
	Symbol string
	Cmd    int8
	Volume int64
}

// 进行转组，订单对冲
func (s *accountHedgingService) ChangeGroupHedging(login int64, group, toGroup string, hedgeAccount, managerId int64) (bool, error) {
	book, err := GroupService.GetBookByGroup(group)
	if err != nil {
		return false, errors.WithMessage(err, util.ErrWrapMessage())
	}
	toBook, err := GroupService.GetBookByGroup(toGroup)
	if err != nil {
		return false, errors.WithMessage(err, util.ErrWrapMessage())
	}

	openingTrades, err := TradeService.GetOpeningTrades(nil, login)
	if err != nil {
		return false, errors.WithMessage(err, util.ErrWrapMessage())
	}

	hedgingTrades, err := s.GetHedgingTrades(openingTrades, book, toBook)
	if err != nil {
		return false, errors.WithMessage(err, util.ErrWrapMessage())
	}

	afterSplitTrades, err := s.SplitTradeLots(hedgingTrades)
	if err != nil {
		return false, errors.WithMessage(err, util.ErrWrapMessage())
	}

	//监测转组记录中账号的转组状态（如果还有账号转组中的记录则不允许此次转组）
	isProcessingExisted, err := ChangeGroupRecordService.IsProcessingExisted(login)
	if err != nil {
		return false, errors.WithMessage(err, util.ErrWrapMessage())
	}
	if isProcessingExisted {
		return false, errors.New("The login has not been processed completely")
	}

	openingTradesOfJson, err := json.Marshal(openingTrades)
	if err != nil {
		return false, errors.Wrapf(err, util.GetCurrentFuncName())
	}
	logger.Infof("service.account_hedging_service.ChangeGroupHedging: login:%v, group:%v, to group:%v, opening trades:%v \n", login, group, toGroup, string(openingTradesOfJson))

	//进行转组
	err = s.ChangeGroup(login, group, toGroup)
	if err != nil {
		return false, errors.WithMessage(err, util.ErrWrapMessage())
	}

	//获取账号的浮动盈亏
	accountFloatingPLs, err := CacheService.GetAccountFloatingPLMap([]int64{login})
	if err != nil {
		return false, errors.WithMessage(err, util.ErrWrapMessage())
	}

	//添加转组记录
	openTicket := ""
	for _, openingTrade := range openingTrades {
		if openTicket == "" {
			openTicket = strconv.FormatInt(openingTrade.Ticket, 10)
		} else {
			openTicket = openTicket + "," + strconv.FormatInt(openingTrade.Ticket, 10)
		}
	}
	changeGroupRecord := model.ChangeGroupRecord{
		Login:       login,
		BeforeGroup: group,
		AfterGroup:  toGroup,
		OpenTicket:  openTicket,
		FloatingPl:  accountFloatingPLs[login],
		ManagerId:   managerId,
		Status:      0,
	}
	changeGroupRecordId, err := ChangeGroupRecordService.Add(changeGroupRecord)
	if err != nil {
		return false, errors.WithMessage(err, util.ErrWrapMessage())
	}

	//对冲账号的持仓单
	err = s.HedgeTrades(hedgeAccount, afterSplitTrades, changeGroupRecordId, login, managerId)
	if err != nil {
		return false, errors.WithMessage(err, util.ErrWrapMessage())
	}

	//更改账号转组记录状态
	err = ChangeGroupRecordService.CompleteChangeGroup(changeGroupRecordId)
	if err != nil {
		return false, errors.WithMessage(err, util.ErrWrapMessage())
	}
	return true, nil
}

// 获取对冲订单
func (s *accountHedgingService) GetHedgingTrades(openingTrades []model.OpeningTrade, book, toBook string) (hedgingTrades []HedgingTrade, err error) {
	//聚合同一品种订单
	symbolVolumes := make(map[string]int64)
	for _, openingTrade := range openingTrades {
		volume := openingTrade.Volume
		if openingTrade.Cmd == 1 {
			volume = openingTrade.Volume * -1
		}

		if _, ok := symbolVolumes[openingTrade.Symbol]; ok {
			symbolVolumes[openingTrade.Symbol] = symbolVolumes[openingTrade.Symbol] + volume
		} else {
			symbolVolumes[openingTrade.Symbol] = volume
		}
	}

	for symbol, volume := range symbolVolumes {
		//判断品种是否在可交易时间内
		isPossibleTrade, err := SymbolService.IsPossibleTrade(symbol)
		if err != nil {
			return nil, errors.WithMessage(err, util.ErrWrapMessage())
		}
		if !isPossibleTrade {
			return nil, errors.New(symbol + "%s closing, unable to transfer group")
		}

		hedgingVolume := volume
		//A组转B组，需要上层的LP进行反向下单对冲，B组转A组，需要上层的LP进行正向下单对冲
		if book == "A" && toBook == "B" {
			hedgingVolume = volume * int64(-1)
		}

		if hedgingVolume == 0 {
			delete(symbolVolumes, symbol)
		}

		cmd := 0
		if hedgingVolume < 0 {
			cmd = 1
			hedgingVolume = hedgingVolume * -1
		}
		hedgingTrades = append(hedgingTrades, HedgingTrade{
			Symbol: symbol,
			Cmd:    int8(cmd),
			Volume: hedgingVolume,
		})
	}
	return hedgingTrades, nil
}

// 拆分手数过大的订单为多个订单
func (s *accountHedgingService) SplitTradeLots(hedgingTrades []HedgingTrade) (afterSplitTrades []AfterSplitTrade, err error) {
	for _, trade := range hedgingTrades {
		volume := trade.Volume
		symbolInfo, err := CacheService.GetSymbolInfo(trade.Symbol)
		if err != nil {
			return nil, errors.WithMessage(err, util.ErrWrapMessage())
		}
		maxVolume := int64(symbolInfo.LotMax * 100)

		//判断手数是否大于最大交易手数（大于：拆分，不大于不拆分）
		if volume > maxVolume {
			//获取手数是最大手数的倍数（向下取整）
			times := volume / maxVolume
			//根据倍数添加最大手数订单
			for i := int64(0); i < times; i++ {
				afterSplitTrades = append(afterSplitTrades, AfterSplitTrade{
					Symbol: trade.Symbol,
					Cmd:    trade.Cmd,
					Volume: maxVolume,
				})
			}

			//添加取余的手数订单
			if volume%maxVolume > 0 {
				afterSplitTrades = append(afterSplitTrades, AfterSplitTrade{
					Symbol: trade.Symbol,
					Cmd:    trade.Cmd,
					Volume: volume % maxVolume,
				})
			}
		} else {
			afterSplitTrades = append(afterSplitTrades, AfterSplitTrade{
				Symbol: trade.Symbol,
				Cmd:    trade.Cmd,
				Volume: volume,
			})
		}
	}
	return afterSplitTrades, nil
}

func (s *accountHedgingService) ChangeGroup(login int64, group, toGroup string) error {
	err := tsp.ModifyGroups(context.Background(), "test request id", login, group, toGroup)
	if err != nil {
		return errors.WithMessage(err, util.ErrWrapMessage())
	}
	return nil
}

func (s *accountHedgingService) HedgeTrades(hedgeAccount int64, tickets []AfterSplitTrade, changeGroupRecordId, login, managerId int64) error {
	ticketsOfJson, err := json.Marshal(tickets)
	if err != nil {
		return errors.Wrapf(err, util.GetCurrentFuncName())
	}
	logger.Infof("service.account_hedging_service.hedgeTrades:all ticket:hedge account:v%, all ticket:v%", hedgeAccount, string(ticketsOfJson))
	hedgingTrades := make([]model.HedgeTrade, 0)
	for _, ticket := range tickets {
		orderInfo, err := tsp.OpenTrade(context.Background(), "test request id", hedgeAccount, ticket.Symbol, ticket.Cmd, ticket.Volume)
		if err != nil {
			logger.Infof("service.account_hedging_service.hedgeTrades:error ticket:hedge account:v%, symbol:v%, cmd:v%, volume:v%", hedgeAccount, ticket.Symbol, ticket.Cmd, ticket.Volume)
			return errors.WithMessage(err, util.ErrWrapMessage())
		}
		hedgingTrade := model.HedgeTrade{
			Login:               login,
			HedgeAccount:        hedgeAccount,
			Ticket:              orderInfo.Ticket,
			Symbol:              orderInfo.Symbol,
			Cmd:                 int8(orderInfo.Cmd),
			Volume:              int64(orderInfo.Volume),
			OpenTime:            time.Unix(orderInfo.OpenTime, 0),
			ChangeGroupRecordId: changeGroupRecordId,
			ManagerId:           managerId,
		}
		hedgingTrades = append(hedgingTrades, hedgingTrade)
	}

	if len(hedgingTrades) == 0 {
		return nil
	}
	err = HedgeTradeService.BatchesAdd(hedgingTrades)
	if err != nil {
		return errors.WithMessage(err, util.ErrWrapMessage())
	}
	return nil
}
