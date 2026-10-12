package service

import (
	"backend-rms/model"
	"backend-rms/repo"
	"backend-rms/util"
	"time"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type reportDailyService struct{}

var ReportDailyService = &reportDailyService{}

func (s *reportDailyService) GetLoginBalanceMapByTime(dateTime time.Time, logins []int64) (loginBalanceMap map[int64]decimal.Decimal, err error) {
	lastDateTime := model.GetDailyDateTime(dateTime)

	loginBalanceMap, err = repo.ReportDailyRepo.GetLoginBalanceMapByTime(lastDateTime.ToDateTimeString(), logins)
	return loginBalanceMap, errors.WithMessage(err, util.ErrWrapMessage())
}

func (s *reportDailyService) GetReportDailyTime(dateTime time.Time) (reportDailyTime *time.Time, err error) {
	lastDateTime := model.GetDailyDateTime(dateTime)

	reportDailyTime, err = repo.ReportDailyRepo.GetReportDailyTime(lastDateTime.ToDateTimeString())
	return reportDailyTime, errors.WithMessage(err, util.ErrWrapMessage())
}
