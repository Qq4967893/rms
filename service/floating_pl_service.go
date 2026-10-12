package service

import (
	"backend-rms/grpc/tsp"
	"backend-rms/model"
	"backend-rms/util"
	"context"
	"time"

	"github.com/golang-module/carbon/v2"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type FloatingPLInfo struct {
	FloatingPL         decimal.Decimal
	FloatingProfit     decimal.Decimal
	FloatingSwap       decimal.Decimal
	FloatingCommission decimal.Decimal
}

type floatingPLService struct{}

var FloatingPLService = &floatingPLService{}

func (s *floatingPLService) GetPointAccountPLMap(groups []string, logins []int64, dateTime time.Time, requestID string) (map[int64]decimal.Decimal, error) {
	var trades []model.OpeningTrade
	chunkSize := 1000
	if len(logins) > chunkSize {
		for i := 0; i < len(logins); i += chunkSize {
			end := i + chunkSize
			if end > len(logins) {
				end = len(logins)
			}
			results, err := TradeService.GetPointOpeningTradeList(groups, logins[i:end], dateTime, "")
			if err != nil {
				return nil, errors.WithMessage(err, util.ErrWrapMessage())
			}
			trades = append(trades, results...)
		}
	} else {
		var err error
		trades, err = TradeService.GetPointOpeningTradeList(groups, logins, dateTime, "")
		if err != nil {
			return nil, errors.WithMessage(err, util.ErrWrapMessage())
		}
	}
	return s.GetPointAccountFloatingPL(trades, dateTime, requestID)
}

func (s *floatingPLService) GetPointAccountFloatingPL(trades []model.OpeningTrade, dateTime time.Time, requestID string) (map[int64]decimal.Decimal, error) {
	historyFloatingPLs, err := s.GetPointFloatingPLInfo(trades, dateTime, requestID)
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}

	accountHistoryFloatingPL := map[int64]decimal.Decimal{}
	for _, v := range trades {
		if _, ok := accountHistoryFloatingPL[v.Login]; ok {
			accountHistoryFloatingPL[v.Login] = accountHistoryFloatingPL[v.Login].Add(historyFloatingPLs[v.Ticket].FloatingPL)
		} else {
			accountHistoryFloatingPL[v.Login] = historyFloatingPLs[v.Ticket].FloatingPL
		}
	}
	return accountHistoryFloatingPL, nil
}

func (s *floatingPLService) GetPointFloatingPLInfo(trades []model.OpeningTrade, dateTime time.Time, requestID string) (map[int64]FloatingPLInfo, error) {
	allHistoryTicks, err := TickService.GetAllHistoryTicks(context.Background(), dateTime)
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}

	symbolInfo, err := tsp.GetSpreadInfo(context.Background(), requestID)
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}

	floatingProfit, err := ProfitService.CalcProfitBatch(trades, symbolInfo, allHistoryTicks)
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}

	date := carbon.Time2Carbon(dateTime).ToDateString()
	historySwaps, err := SwapService.GetPointSwaps(date)
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}

	historyFloatingPLs := map[int64]FloatingPLInfo{}
	for _, v := range trades {
		historySwap := decimal.Zero
		if _, ok := historySwaps[v.Ticket]; ok {
			historySwap = historySwaps[v.Ticket]
		}
		historyFloatingPL := v.Commission.Add(floatingProfit[v.Ticket].GrossProfit).Add(historySwap)
		historyFloatingPLs[v.Ticket] = FloatingPLInfo{
			FloatingPL:         historyFloatingPL,
			FloatingProfit:     floatingProfit[v.Ticket].GrossProfit,
			FloatingSwap:       historySwap,
			FloatingCommission: v.Commission,
		}
	}
	return historyFloatingPLs, nil
}
