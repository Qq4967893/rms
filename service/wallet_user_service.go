package service

import (
	"backend-rms/repo"
	"backend-rms/util"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type walletUserService struct{}

var WalletUserService = &walletUserService{}

func (s *walletUserService) GetBalanceByPassport(passport int64) (balance decimal.Decimal, err error) {
	balance, err = repo.WalletUserRepo.GetBalanceByPassport(passport)
	return balance, errors.WithMessage(err, util.ErrWrapMessage())
}
