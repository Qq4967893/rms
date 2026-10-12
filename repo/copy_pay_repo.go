package repo

import (
	"backend-rms/core/database"
	"backend-rms/model"
	"backend-rms/util"

	"github.com/pkg/errors"
)

type copyPayRepo struct{}

var CopyPayRepo = &copyPayRepo{}

func (r *copyPayRepo) GetAll(followerId, masterId int64) (copyPays []model.CopyPay, err error) {
	result := database.APIDB.Table("copy_pay").
		Where("mt4_from", masterId).
		Where("mt4_to", followerId).
		Select("mt4_to, mt4_from, pay_zone, pay_time, pay_master").
		Find(&copyPays)
	return copyPays, errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
}
