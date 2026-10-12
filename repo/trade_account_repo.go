package repo

import (
	"backend-rms/core/database"
	"backend-rms/model"
	"backend-rms/util"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type tradeAccountRepo struct{}

var TradeAccountRepo = &tradeAccountRepo{}

func (r *tradeAccountRepo) GetTradeAccountIDs(passport int64) (tradeAccountIDs []int64, err error) {
	err = database.APIDB.Table("trade_account").
		Scopes(r.liveTradeAccount).
		Where("status", 1).
		Where("user_code", passport).
		Pluck("trade_account_id", &tradeAccountIDs).Error
	return tradeAccountIDs, errors.Wrapf(err, "%s failed", util.GetCurrentFuncName())
}

func (r *tradeAccountRepo) GetPassportByLogin(login int64) (passport int64, err error) {
	err = database.APIDB.Table("trade_account").
		Select("user_code").
		Scopes(r.liveTradeAccount).
		Where("status", 1).
		Where("trade_account_id", login).
		Scan(&passport).Error
	return passport, errors.Wrapf(err, "%s failed", util.GetCurrentFuncName())
}

func (r *tradeAccountRepo) GetLoginPassportMap(logins []int64) (map[int64]int64, error) {
	var tradeAccounts []*model.TradeAccount
	result := database.APIDB.Table("trade_account").
		Where("status", 1).
		Where("trade_account_id IN ?", logins).
		Select("trade_account_id, user_code as passport").Find(&tradeAccounts)
	if result.Error != nil {
		return nil, errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
	}

	loginPassportMap := make(map[int64]int64, len(tradeAccounts))
	for _, tradeAccount := range tradeAccounts {
		loginPassportMap[tradeAccount.TradeAccountID] = tradeAccount.Passport
	}

	return loginPassportMap, nil
}

func (r *tradeAccountRepo) GetLoginTypeMap(logins []int64) (map[int64]int8, error) {
	var tradeAccounts []*model.TradeAccount
	result := database.APIDB.Table("trade_account").
		Where("trade_account_id IN ?", logins).
		Select("trade_account_id, trade_account_type").Find(&tradeAccounts)
	if result.Error != nil {
		return nil, errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
	}

	loginTypeMap := make(map[int64]int8, len(tradeAccounts))
	for _, tradeAccount := range tradeAccounts {
		loginTypeMap[tradeAccount.TradeAccountID] = tradeAccount.TradeAccountType
	}

	return loginTypeMap, nil
}

func (r *tradeAccountRepo) liveTradeAccount(db *gorm.DB) *gorm.DB {
	return db.Where("trade_account_type > ?", 0)
}
