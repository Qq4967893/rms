package repo

import (
	"backend-rms/core/database"
	"backend-rms/model"
	"backend-rms/util"

	"github.com/pkg/errors"
)

type spreadRecordRepo struct{}

var SpreadRecordRepo = &spreadRecordRepo{}

func (r *spreadRecordRepo) GetList() (spreads []model.SpreadRecord, err error) {
	result := database.APIDB.Select("symbol, before_spread, after_spread, change_time").Order("change_time desc").Find(&spreads)
	return spreads, errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
}
