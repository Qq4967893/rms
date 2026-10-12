package repo

import (
	"backend-rms/core/database"
	"backend-rms/model"
	"backend-rms/util"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type walletUserRepo struct{}

var WalletUserRepo = &walletUserRepo{}

func (r *walletUserRepo) GetBalanceByPassport(passport int64) (balance decimal.Decimal, err error) {
	err = database.APIDB.Model(model.WalletUser{}).Table("wallet_user").
		Select("balance").
		Where("user_code", passport).Scan(&balance).Error
	return balance, errors.Wrapf(err, "%s failed", util.GetCurrentFuncName())
}
