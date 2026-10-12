package service

import (
	"backend-rms/model"
	"backend-rms/repo"
	"backend-rms/util"
	"time"

	"github.com/pkg/errors"
)

type dailyClosedTradeService struct {
}

var DailyClosedTradeService = &dailyClosedTradeService{}

func (ds *dailyClosedTradeService) GetSymbolList(groups []string, symbol string, startDate, endDate time.Time) (symbolList []model.SymbolDailyClosedTrade, err error) {
	list, err := repo.DailyClosedTradeRepo.GetSymbolList(groups, symbol, startDate, endDate)
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}

	symbolListMap := make(map[string]model.SymbolDailyClosedTrade)
	for _, v := range list {
		symbolOrigin, err := SymbolLeverageService.GetSymbolOrigin(v.Symbol)
		if err != nil {
			return nil, errors.WithMessage(err, util.ErrWrapMessage())
		}

		if _, ok := symbolListMap[symbolOrigin]; ok {
			symbolListMap[symbolOrigin] = model.SymbolDailyClosedTrade{
				Symbol:             symbolOrigin,
				Profit:             symbolListMap[symbolOrigin].Profit.Add(v.Profit),
				Swaps:              symbolListMap[symbolOrigin].Swaps.Add(v.Swaps),
				Commission:         symbolListMap[symbolOrigin].Commission.Add(v.Commission),
				RawSpreadProfit:    symbolListMap[symbolOrigin].RawSpreadProfit.Add(v.RawSpreadProfit),
				MarkUpSpreadProfit: symbolListMap[symbolOrigin].MarkUpSpreadProfit.Add(v.MarkUpSpreadProfit),
				Lots:               symbolListMap[symbolOrigin].Lots.Add(v.Lots),
			}
		} else {
			symbolListMap[symbolOrigin] = model.SymbolDailyClosedTrade{
				Symbol:             symbolOrigin,
				Profit:             v.Profit,
				Swaps:              v.Swaps,
				Commission:         v.Commission,
				RawSpreadProfit:    v.RawSpreadProfit,
				MarkUpSpreadProfit: v.MarkUpSpreadProfit,
				Lots:               v.Lots,
			}
		}
	}

	for _, v := range symbolListMap {
		middle := model.SymbolDailyClosedTrade{
			Symbol:             v.Symbol,
			Profit:             v.Profit,
			Swaps:              v.Swaps,
			Commission:         v.Commission,
			RawSpreadProfit:    v.RawSpreadProfit,
			MarkUpSpreadProfit: v.MarkUpSpreadProfit,
			Lots:               v.Lots,
		}
		symbolList = append(symbolList, middle)
	}
	return symbolList, nil
}

func (ds *dailyClosedTradeService) GetAccountList(groups []string, symbol string, startDate, endDate time.Time, paging *model.Paging) (accountList []model.AccountDailyClosedTrade, err error) {
	return repo.DailyClosedTradeRepo.GetAccountList(groups, symbol, startDate, endDate, paging)
}

func (ds *dailyClosedTradeService) GetAccountCount(groups []string, symbol string, startDate, endDate time.Time) (int64, error) {
	return repo.DailyClosedTradeRepo.GetAccountCount(groups, symbol, startDate, endDate)
}

func (ds *dailyClosedTradeService) GetTradeList(groups []string, symbol string, startDate, endDate time.Time, paging *model.Paging) (tradeList []model.DailyClosedTradeInfo, err error) {
	return repo.DailyClosedTradeRepo.GetTradeList(groups, symbol, startDate, endDate, paging)
}

func (ds *dailyClosedTradeService) GetTradeCount(groups []string, symbol string, startDate, endDate time.Time) (int64, error) {
	return repo.DailyClosedTradeRepo.GetTradeCount(groups, symbol, startDate, endDate)
}

func (ds *dailyClosedTradeService) Add(params []model.DailyClosedTrade) error {
	err := repo.DailyClosedTradeRepo.Add(params)
	if err != nil {
		return errors.WithMessage(err, util.ErrWrapMessage())
	}
	return nil
}
