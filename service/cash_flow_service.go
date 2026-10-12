package service

import (
	"backend-rms/repo"
	"backend-rms/util"
	"time"

	"github.com/golang-module/carbon/v2"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type cashFlowService struct{}

var CashFlowService = &cashFlowService{}

func (s *cashFlowService) GetNetCashFlowMap(logins []int64, startTime, endTime string) (map[int64]decimal.Decimal, error) {
	netCashFlowMap, err := repo.CashFlowRepo.GetCashFlowMap(logins, startTime, endTime, 0)
	return netCashFlowMap, errors.WithMessage(err, util.ErrWrapMessage())
}

func (s *cashFlowService) GetCashInflowMap(logins []int64, startTime, endTime string) (map[int64]decimal.Decimal, error) {
	cashInflowMap, err := repo.CashFlowRepo.GetCashFlowMap(logins, startTime, endTime, 1)
	return cashInflowMap, errors.WithMessage(err, util.ErrWrapMessage())
}

func (s *cashFlowService) GetCashOutflowMap(logins []int64, startTime, endTime string) (map[int64]decimal.Decimal, error) {
	cashOutflowMap, err := repo.CashFlowRepo.GetCashFlowMap(logins, startTime, endTime, -1)
	return cashOutflowMap, errors.WithMessage(err, util.ErrWrapMessage())
}

func (s *cashFlowService) GetCashInflow(groups []string, startTime, endTime time.Time) (amount decimal.Decimal, err error) {
	amount, err = s.getCashFlow(groups, startTime, endTime, 1)
	if err != nil {
		return decimal.Zero, errors.WithMessage(err, util.ErrWrapMessage())
	}
	return
}

func (s *cashFlowService) GetCashOutflow(groups []string, startTime, endTime time.Time) (amount decimal.Decimal, err error) {
	untreatedAmount, err := s.getCashFlow(groups, startTime, endTime, 2)
	//转换为正值
	amount = untreatedAmount.Mul(decimal.NewFromInt(-1))
	if err != nil {
		return decimal.Zero, errors.WithMessage(err, util.ErrWrapMessage())
	}
	return
}

func (s *cashFlowService) GetNetCashFlow(groups []string, startTime, endTime time.Time) (amount decimal.Decimal, err error) {
	amount, err = s.getCashFlow(groups, startTime, endTime, 0)
	if err != nil {
		return decimal.Zero, errors.WithMessage(err, util.ErrWrapMessage())
	}
	return
}

func (s *cashFlowService) getCashFlow(groups []string, startTime, endTime time.Time, paymentType int8) (amount decimal.Decimal, err error) {
	//判断endTime，如果是当天连接user表否则连接daily表
	todayZero := carbon.Now().StartOfDay().ToDateTimeString()
	if carbon.Time2Carbon(endTime).Gte(carbon.Parse(todayZero)) {
		amount, err = repo.CashFlowRepo.GetCashFlowWithAccount(groups, startTime, endTime, paymentType)
		if err != nil {
			return decimal.Zero, errors.WithMessage(err, util.ErrWrapMessage())
		}
		return
	} else {
		amount, err = repo.CashFlowRepo.GetCashFlowWithDaily(groups, startTime, endTime, paymentType)
		if err != nil {
			return decimal.Zero, errors.WithMessage(err, util.ErrWrapMessage())
		}
		return
	}
}
