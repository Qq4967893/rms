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

type balanceService struct{}

var BalanceService = &balanceService{}

func (s *balanceService) GetPointBalance(groups []string, dateTime time.Time) (decimal.Decimal, error) {
	balanceDateTime, err := s.buildPointBalanceDateTime(dateTime)
	if err != nil {
		return decimal.Zero, errors.WithMessage(err, util.ErrWrapMessage())
	}

	balance, err := repo.ReportDailyRepo.GetBalance(groups, balanceDateTime["startTime"])
	if err != nil {
		return decimal.Zero, errors.WithMessage(err, util.ErrWrapMessage())
	}

	paymentAmount, err := CashFlowService.GetNetCashFlow(groups, balanceDateTime["startTime"], balanceDateTime["endTime"])
	if err != nil {
		return decimal.Zero, errors.WithMessage(err, util.ErrWrapMessage())
	}

	closedPL, err := TradeService.GetClosedPL(groups, balanceDateTime["startTime"], balanceDateTime["endTime"])
	if err != nil {
		return decimal.Zero, errors.WithMessage(err, util.ErrWrapMessage())
	}

	return balance.Add(paymentAmount).Add(closedPL), nil
}

func (s *balanceService) buildPointBalanceDateTime(dateTime time.Time) (map[string]time.Time, error) {
	yesterday := carbon.Time2Carbon(dateTime).Yesterday().Carbon2Time()
	dailyDateTime := model.GetDailyDateTime(yesterday)
	balanceDateTime := map[string]time.Time{
		"startTime": dailyDateTime.Carbon2Time(),
		"endTime":   dateTime,
	}
	return balanceDateTime, nil
}
