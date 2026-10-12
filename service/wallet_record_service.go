package service

import (
	"backend-rms/repo"
	"backend-rms/util"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type walletRecordService struct{}

var WalletRecordService = &walletRecordService{}

func (s *walletRecordService) GetDepositByPassport(passport int64) (deposit decimal.Decimal, err error) {
	deposit, err = repo.WalletRecordRepo.GetDepositByPassport(passport)
	return deposit, errors.WithMessage(err, util.ErrWrapMessage())
}

func (s *walletRecordService) GetWithdrawalByPassport(passport int64) (withdrawal decimal.Decimal, err error) {
	withdrawal, err = repo.WalletRecordRepo.GetWithdrawalByPassport(passport)
	return withdrawal, errors.WithMessage(err, util.ErrWrapMessage())
}
