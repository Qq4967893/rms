package repo

import (
	"backend-rms/core/database"
	"backend-rms/model"
	"backend-rms/util"
	"time"

	"github.com/pkg/errors"
)

type lpSpreadRecordRepo struct {
}

var LpSpreadRecordRepo = &lpSpreadRecordRepo{}

func (lr *lpSpreadRecordRepo) GetAvgSpreadLast14Days(symbol string, lpType int8, dateTime time.Time) (lpSpreadRecords []*model.LpSpreadRecord, err error) {
	result := database.ETCDB.Table("lp_spread_record").
		Where("lp_avg > ?", 0).
		Where("source", symbol).
		Where("lp_type", lpType).
		Where("lp_date <= ?", dateTime).
		Order("lp_date desc").
		Limit(14).
		Find(&lpSpreadRecords)
	return lpSpreadRecords, errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
}
