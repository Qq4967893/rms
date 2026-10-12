package service

import (
	"backend-rms/grpc/tsp"
	"backend-rms/model"
	"backend-rms/repo"
	"backend-rms/util"
	"context"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type accountService struct{}

var AccountService = &accountService{}

type EquityInfo struct {
	Login       int64
	Balance     decimal.Decimal
	Equity      decimal.Decimal
	MarginLevel float64
	Margin      decimal.Decimal
	MarginFree  decimal.Decimal
}

func (s *accountService) GetInfoByLogin(login int64) (account *model.Account, err error) {
	account, err = repo.AccountRepo.GetInfoByLogin(login)
	return account, errors.WithMessage(err, util.ErrWrapMessage())
}

func (s *accountService) GetListByLogins(logins []int64, groups []string) (accounts []*model.Account, err error) {
	accounts, err = repo.AccountRepo.GetListByLogins(logins, groups)
	return accounts, errors.WithMessage(err, util.ErrWrapMessage())
}

func (s *accountService) GetList(groups []string) (accounts []*model.Account, err error) {
	accounts, err = repo.AccountRepo.GetAll(groups)
	return accounts, errors.WithMessage(err, util.ErrWrapMessage())
}

func (s *accountService) GetReportAccountList(groups []string, logins []int64, paging *model.Paging) (accounts []*model.ReportAccount, err error) {
	accounts, err = repo.AccountRepo.GetReportAccountList(groups, logins, paging)
	return accounts, errors.WithMessage(err, util.ErrWrapMessage())
}

func (s *accountService) GetAllReportAccounts(groups []string) (accounts []*model.ReportAccount, err error) {
	accounts, err = repo.AccountRepo.GetReportAccountList(groups, []int64{}, nil)
	return accounts, errors.WithMessage(err, util.ErrWrapMessage())
}

func (s *accountService) GetReportAccountCount(groups []string, login []int64) (totalCount int64, err error) {
	totalCount, err = repo.AccountRepo.GetReportAccountCount(groups, login)
	return totalCount, errors.WithMessage(err, util.ErrWrapMessage())
}

func (s *accountService) GetEquityMap(ctx context.Context, requestID string, logins []int64) (map[int64]*EquityInfo, error) {
	equityInfos, err := tsp.GetUserEquity(ctx, requestID, logins)
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}

	equityMap := make(map[int64]*EquityInfo, len(equityInfos))
	for _, equityInfo := range equityInfos {
		equityMap[int64(equityInfo.Login)] = &EquityInfo{
			Login:       int64(equityInfo.Login),
			Balance:     decimal.NewFromFloat(equityInfo.Balance),
			Equity:      decimal.NewFromFloat(equityInfo.Equity),
			MarginLevel: equityInfo.MarginLevel,
			Margin:      decimal.NewFromFloat(equityInfo.Margin),
			MarginFree:  decimal.NewFromFloat(equityInfo.MarginFree),
		}
	}
	return equityMap, nil
}
