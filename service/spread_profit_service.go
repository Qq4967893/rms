package service

import (
	"backend-rms/model"
	"backend-rms/util"
	"context"
	"time"

	"github.com/golang-module/carbon/v2"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type spreadProfitService struct {
}

type SpreadProfit struct {
	SpreadProfit           decimal.Decimal
	WithoutRawSpreadProfit decimal.Decimal
	CloseDate              time.Time
}

var SpreadProfitService = &spreadProfitService{}

func (sps *spreadProfitService) CalcSpreadProfit(trades []model.ClosedTrade) (map[int64]SpreadProfit, error) {
	groupMap, err := GroupService.GetAllMap()
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}

	//初始化品种对应的LP品种的点差修改记录
	err = SpreadService.InitSpreadList()
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}

	tradeSpreadProfit := make(map[int64]SpreadProfit)
	for _, tradeInfo := range trades {
		if tradeInfo.CloseTime.Equal(carbon.Parse(model.DEFAULT_CLOSE_TIME).Carbon2Time()) {
			continue
		}

		if _, ok := groupMap[tradeInfo.TransCloseGroup]; !ok {
			continue
		}

		ticket := tradeInfo.Ticket
		symbol := tradeInfo.Symbol
		cmd := tradeInfo.Cmd
		lots := decimal.NewFromInt(tradeInfo.Volume).Mul(decimal.NewFromFloat(0.01))
		openTime := tradeInfo.OpenTime
		closeTime := tradeInfo.CloseTime
		openPrice := tradeInfo.OpenPrice
		closePrice := tradeInfo.ClosePrice
		transOpenPrice := tradeInfo.TransOpenPrice
		transClosePrice := tradeInfo.TransClosePrice
		depositCurrency := tradeInfo.DepositCurrency

		//品种去杠杆
		symbolOrigin, err := SymbolLeverageService.GetSymbolOrigin(symbol)
		if err != nil {
			return nil, errors.WithMessage(err, util.ErrWrapMessage())
		}

		symbolInfo, err := CacheService.GetSymbolInfo(symbolOrigin)
		if err != nil {
			return nil, errors.WithMessage(err, util.ErrWrapMessage())
		}
		openSpreadOffset := SpreadService.GetSpreadOffset(cmd, "open", openPrice, transOpenPrice, symbolInfo.Digits)
		openSpread, err := SpreadService.GetSpreadByDate(symbolOrigin, openSpreadOffset, openTime)
		if err != nil {
			return nil, errors.WithMessage(err, util.ErrWrapMessage())
		}

		closeSpreadOffset := SpreadService.GetSpreadOffset(cmd, "close", closePrice, transClosePrice, symbolInfo.Digits)
		closeSpread, err := SpreadService.GetSpreadByDate(symbolOrigin, closeSpreadOffset, closeTime)
		if err != nil {
			return nil, errors.WithMessage(err, util.ErrWrapMessage())
		}

		//获取开平仓时间点品种的14天平均点差
		rawCloseSpread, err := SpreadService.GetRawSpread(symbolOrigin, closeTime)
		if err != nil {
			return nil, errors.WithMessage(err, util.ErrWrapMessage())
		}
		rawOpenSpread, err := SpreadService.GetRawSpread(symbolOrigin, openTime)
		if err != nil {
			return nil, errors.WithMessage(err, util.ErrWrapMessage())
		}
		if rawOpenSpread == -1 {
			rawOpenSpread = rawCloseSpread
		}

		grossProfit, err := sps.GetGrossProfit(symbolOrigin, cmd, closeTime, depositCurrency)
		if err != nil {
			return nil, errors.WithMessage(err, util.ErrWrapMessage())
		}

		withoutRawOpenSpreadProfit := grossProfit.Mul(lots).Mul(openSpread.Div(decimal.NewFromInt(2)))
		withoutRawCloseSpreadProfit := grossProfit.Mul(lots).Mul(closeSpread.Div(decimal.NewFromInt(2)))

		openSpread = openSpread.Add(decimal.NewFromInt(rawOpenSpread))
		closeSpread = closeSpread.Add(decimal.NewFromInt(rawCloseSpread))

		openSpreadProfit := grossProfit.Mul(lots).Mul(openSpread.Div(decimal.NewFromInt(2)))
		closeSpreadProfit := grossProfit.Mul(lots).Mul(closeSpread.Div(decimal.NewFromInt(2)))

		tradeSpreadProfit[ticket] = SpreadProfit{
			SpreadProfit:           openSpreadProfit.Add(closeSpreadProfit),
			WithoutRawSpreadProfit: withoutRawOpenSpreadProfit.Add(withoutRawCloseSpreadProfit),
			CloseDate:              tradeInfo.CloseTime,
		}
	}
	return tradeSpreadProfit, nil
}

func (sps *spreadProfitService) GetGrossProfit(symbol string, cmd int8, dateTime time.Time, depositCurrency string) (decimal.Decimal, error) {
	var quotes map[string]*Quote

	timeStamp := dateTime.UnixMilli()
	symbolInfo, err := CacheService.GetSymbolInfo(symbol)
	if err != nil {
		return decimal.Zero, errors.WithMessage(err, util.ErrWrapMessage())
	}

	profitCurrency := symbolInfo.ProfitCurrency

	var bid decimal.Decimal
	if profitCurrency != "USD" {
		quote, err := TickService.GetHistoryTick(context.Background(), profitCurrency, dateTime)
		if err != nil {
			return decimal.Zero, errors.WithMessage(err, util.ErrWrapMessage())
		}

		bid = quote.Bid
		if bid == decimal.Zero {
			cmd = 0
			bid, err = HistoryQuoteService.GetLastClosePrice(profitCurrency, timeStamp)
			if err != nil {
				return decimal.Zero, errors.WithMessage(err, util.ErrWrapMessage())
			}
			quote.Bid = bid
		}

		quotes = map[string]*Quote{
			profitCurrency: quote,
		}
	}

	pipProfit, err := ProfitService.CalcPipProfitByQuotes(*symbolInfo, cmd, quotes)
	if err != nil {
		return decimal.Zero, errors.WithMessage(err, util.ErrWrapMessage())
	}
	grossProfit := pipProfit.GrossProfit

	//如果入金币种是GBP，需要把USD转换为GBP，需要提供GBPUSD的报价
	if depositCurrency == "GBP" {
		quote, err := TickService.GetHistoryTick(context.Background(), "GBPUSD", dateTime)
		if err != nil {
			return decimal.Zero, errors.WithMessage(err, util.ErrWrapMessage())
		}

		bid = quote.Bid
		if bid == decimal.Zero {
			cmd = 0
			bid, err = HistoryQuoteService.GetLastClosePrice("GBPUSD", timeStamp)
			if err != nil {
				return decimal.Zero, errors.WithMessage(err, util.ErrWrapMessage())
			}
			quote.Bid = bid
		}

		quotes = map[string]*Quote{
			profitCurrency: quote,
		}
		grossProfit, err = ProfitService.usdToOtherCurrency(grossProfit, quotes, "GBP")
		if err != nil {
			return decimal.Zero, errors.WithMessage(err, util.ErrWrapMessage())
		}
	}
	return grossProfit, nil
}
