package repo

import (
	"backend-rms/core/database"
	"backend-rms/model"
	"backend-rms/util"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type hedgeAccountRepo struct {
}

var HedgeAccountRepo = &hedgeAccountRepo{}

func (hr *hedgeAccountRepo) GetInfo() (hedgeAccount *model.HedgeAccount, err error) {
	err = database.RMSDB.Table("hedge_accounts").
		Select("id, hedge_account, role_ids, `desc`").
		First(&hedgeAccount).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return hedgeAccount, errors.Wrapf(err, "%s failed", util.GetCurrentFuncName())
}

func (hr *hedgeAccountRepo) GetList() (HedgeAccounts []*model.HedgeAccount, err error) {
	result := database.RMSDB.Find(&HedgeAccounts)
	return HedgeAccounts, errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
}

func (hr *hedgeAccountRepo) Update(id int64, params model.HedgeAccount) error {
	result := database.RMSDB.Where("id", id).Updates(&params)
	return errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
}
