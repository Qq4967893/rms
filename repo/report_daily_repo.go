package repo

import (
	"backend-rms/core/database"
	"backend-rms/model"
	"backend-rms/util"
	"time"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type reportDailyRepo struct {
}

var ReportDailyRepo = &reportDailyRepo{}

func (r *reportDailyRepo) GetBalance(groups []string, dateTime time.Time) (balance decimal.Decimal, err error) {
	result := database.ET6DB.Model(model.ReportDaily{}).Debug().Where("time = ?", dateTime).
		Where("`group` IN ?", groups).
		Select("COALESCE(SUM(balance), 0) as balance").Scan(&balance)
	if result.Error != nil {
		return decimal.Zero, errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
	}
	return
}

func (r *reportDailyRepo) GetLoginBalanceMapByTime(dateTime string, logins []int64) (loginBalanceMap map[int64]decimal.Decimal, err error) {
	tx := database.ET6DB.Select("login, balance").
		Where("time = ?", dateTime)
	if len(logins) > 0 {
		tx.Where("login IN ?", logins)
	}

	var reportDailies []*model.ReportDaily
	result := tx.Find(&reportDailies)
	if result.Error != nil {
		return nil, errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
	}

	loginBalanceMap = make(map[int64]decimal.Decimal, len(reportDailies))
	for _, reportDaily := range reportDailies {
		loginBalanceMap[reportDaily.Login] = reportDaily.Balance
	}

	return loginBalanceMap, nil
}

func (r *reportDailyRepo) GetReportDailyTime(dateTime string) (reportDailyTime *time.Time, err error) {
	tx := database.ET6DB.Model(model.ReportDaily{}).Select("time").
		Where("time = ?", dateTime)
	result := tx.Scan(&reportDailyTime)
	return reportDailyTime, errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
}
