package service

import (
	"backend-rms/core/logger"
	"backend-rms/model"

	"github.com/shopspring/decimal"
)

type summaryService struct{}

var SummaryService = &summaryService{}

func (s *summaryService) HandleForOrder(trade *model.Trade, typeOrder string) {
	var source string
	if symbolInfo, ok := InitializedData.Symbols[trade.Symbol]; ok {
		source = symbolInfo.Source
	} else {
		// 新增订单品种不在初始化的品种内
		logger.Errorf("The order symbol %v is not within the initialized symbol", trade.Symbol)
		return
	}
	logger.Infof("%v order ticket is %v, symbol is %v ,volume is %v", typeOrder, trade.Ticket, trade.Symbol, trade.Volume)

	book := InitializedData.Groups[InitializedData.Accounts[trade.Login].Group].Book
	switch trade.Cmd {
	case 0:
		if typeOrder == "open" {
			SummaryBook[book][source].BuyLots = SummaryBook[book][source].BuyLots.Add(decimal.NewFromInt(trade.Volume).Div(decimal.NewFromInt(100)))
			SummaryBook[book][source].Et6NetLots = SummaryBook[book][source].Et6NetLots.Add(decimal.NewFromInt(trade.Volume).Div(decimal.NewFromInt(100)))
		} else {
			SummaryBook[book][source].BuyLots = SummaryBook[book][source].BuyLots.Sub(decimal.NewFromInt(trade.Volume).Div(decimal.NewFromInt(100)))
			SummaryBook[book][source].Et6NetLots = SummaryBook[book][source].Et6NetLots.Sub(decimal.NewFromInt(trade.Volume).Div(decimal.NewFromInt(100)))
		}
		break
	case 1:
		if typeOrder == "open" {
			SummaryBook[book][source].SellLots = SummaryBook[book][source].SellLots.Add(decimal.NewFromInt(trade.Volume).Div(decimal.NewFromInt(100)))
			SummaryBook[book][source].Et6NetLots = SummaryBook[book][source].Et6NetLots.Sub(decimal.NewFromInt(trade.Volume).Div(decimal.NewFromInt(100)))
		} else {
			SummaryBook[book][source].SellLots = SummaryBook[book][source].SellLots.Sub(decimal.NewFromInt(trade.Volume).Div(decimal.NewFromInt(100)))
			SummaryBook[book][source].Et6NetLots = SummaryBook[book][source].Et6NetLots.Add(decimal.NewFromInt(trade.Volume).Div(decimal.NewFromInt(100)))
		}
		break
	}
	if typeOrder == "open" {
		SummaryBook[book][source].Orders++
	} else {
		SummaryBook[book][source].Orders--
	}
}
