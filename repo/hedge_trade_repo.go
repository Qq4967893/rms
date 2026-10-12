package repo

import (
	"backend-rms/core/database"
	"backend-rms/model"
	"backend-rms/util"

	"github.com/pkg/errors"
)

type hedgeTradeRepo struct{}

var HedgeTradeRepo = &hedgeTradeRepo{}

func (cr *hedgeTradeRepo) GetAll(login, managerId int64, paging *model.Paging) (hedgeTrades []*model.HedgeTrade, err error) {
	tx := database.RMSDB.Table("hedge_trades")
	if login != 0 {
		tx.Where("login", login)
	}

	if managerId != 0 {
		tx.Where("manager_id", managerId)
	}

	if paging != nil {
		tx.Order(paging.Sort + " " + paging.Order).Limit(paging.Limit).Offset(paging.Offset)
	}

	result := tx.Find(&hedgeTrades)
	return hedgeTrades, errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
}

func (cr *hedgeTradeRepo) GetCount(login, managerId int64) (totalCount int64, err error) {
	tx := database.RMSDB.Table("hedge_trades")
	if login != 0 {
		tx.Where("login", login)
	}

	if managerId != 0 {
		tx.Where("manager_id", managerId)
	}

	result := tx.Count(&totalCount)
	return totalCount, errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
}

func (cr *hedgeTradeRepo) Add(hedgeTrades []model.HedgeTrade) error {
	result := database.RMSDB.Create(&hedgeTrades)
	return errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
}
