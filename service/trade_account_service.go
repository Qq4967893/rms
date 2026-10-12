package service

import (
	"backend-rms/repo"
	"backend-rms/util"

	"github.com/pkg/errors"
)

type tradeAccountService struct{}

var TradeAccountService = &tradeAccountService{}

func (s *tradeAccountService) GetTradeAccountIDs(passport int64) (tradeAccountIDs []int64, err error) {
	tradeAccountIDs, err = repo.TradeAccountRepo.GetTradeAccountIDs(passport)
	return tradeAccountIDs, errors.WithMessage(err, util.ErrWrapMessage())
}

func (s *tradeAccountService) GetPassportByLogin(login int64) (passport int64, err error) {
	passport, err = repo.TradeAccountRepo.GetPassportByLogin(login)
	return passport, errors.WithMessage(err, util.ErrWrapMessage())
}

func (s *tradeAccountService) GetLoginPassportMap(logins []int64) (map[int64]int64, error) {
	loginPassportMap, err := repo.TradeAccountRepo.GetLoginPassportMap(logins)
	return loginPassportMap, errors.WithMessage(err, util.ErrWrapMessage())
}

func (s *tradeAccountService) GetLoginTypeMap(logins []int64) (map[int64]int8, error) {
	loginTypeMap, err := repo.TradeAccountRepo.GetLoginTypeMap(logins)
	return loginTypeMap, errors.WithMessage(err, util.ErrWrapMessage())
}
