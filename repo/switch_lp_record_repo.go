package repo

import (
	"backend-rms/core/database"
	"backend-rms/model"
	"backend-rms/util"

	"github.com/pkg/errors"
)

type switchLpRecordRepo struct {
}

var SwitchLpRecordRepo = &switchLpRecordRepo{}

func (r *switchLpRecordRepo) GetAll(symbol string) (switchLpRecord []*model.SwitchLpRecord, err error) {
	result := database.ETCDB.Where("source", symbol).Order("switch_time desc").Find(&switchLpRecord)
	return switchLpRecord, errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
}
