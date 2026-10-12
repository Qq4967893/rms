package repo

import (
	"backend-rms/core/database"
	"backend-rms/util"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type walletRecordRepo struct{}

var WalletRecordRepo = &walletRecordRepo{}

func (r *walletRecordRepo) GetDepositByPassport(passport int64) (deposit decimal.Decimal, err error) {
	err = database.APIDB.Table("wallet").
		Select("COALESCE(SUM(amount),0) AS deposit").
		Where("type", 3).
		Where("direction", 1).
		Where("status", 1).
		Where("payment_platform not in ?", []int{0, 3}).
		Where("user_code", passport).Scan(&deposit).Error
	return deposit, errors.Wrapf(err, "%s failed", util.GetCurrentFuncName())
}

func (r *walletRecordRepo) GetWithdrawalByPassport(passport int64) (withdrawal decimal.Decimal, err error) {
	err = database.APIDB.Table("wallet").
		Select("COALESCE(SUM(amount),0) AS withdrawal").
		Where("type", 3).
		Where("direction", -1).
		Where("status", 1).
		Where("payment_platform != ?", 3).
		Where("finance_status", 1).
		Where("user_code", passport).Scan(&withdrawal).Error
	return withdrawal, errors.Wrapf(err, "%s failed", util.GetCurrentFuncName())
}
