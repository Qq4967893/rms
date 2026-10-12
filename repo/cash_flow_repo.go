package repo

import (
	"backend-rms/core/database"
	"backend-rms/model"
	"backend-rms/util"
	"time"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type cashFlowRepo struct{}

var CashFlowRepo = &cashFlowRepo{}

func (r *cashFlowRepo) GetCashFlowMap(logins []int64, startTime, endTime string, direction int8) (map[int64]decimal.Decimal, error) {
	tx := database.ET6DB.Where("cmd = ?", 6)
	if len(logins) > 0 {
		tx.Where("login IN ?", logins)
	}

	if startTime != "" {
		tx.Where("open_time >= ?", startTime)
	}

	if endTime != "" {
		tx.Where("open_time <= ?", endTime)
	}

	if direction > 0 {
		tx.Where("profit > 0")
	} else if direction < 0 {
		tx.Where("profit < 0")
	}

	var trades []*model.Trade
	result := tx.Select("login, sum(profit) as profit").Group("login").Find(&trades)
	if result.Error != nil {
		return nil, errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
	}

	paymentAmountMap := make(map[int64]decimal.Decimal, len(trades))
	for _, trade := range trades {
		paymentAmountMap[trade.Login] = trade.Profit
	}

	return paymentAmountMap, nil
}

func (r *cashFlowRepo) GetCashFlowWithAccount(groups []string, startTime, endTime time.Time, paymentType int8) (amount decimal.Decimal, err error) {
	tx := TradeRepo.joinAccount().Where("a.group IN ?", groups).
		Where("t.cmd = ?", 6)

	if !startTime.IsZero() {
		tx.Where("t.open_time >= ?", startTime)
	}

	if !endTime.IsZero() {
		tx.Where("t.open_time <= ?", endTime)
	}

	if paymentType == 1 {
		tx.Where("t.profit > 0")
	} else if paymentType == 2 {
		tx.Where("t.profit < 0")
	}
	result := tx.Select("COALESCE(SUM(t.profit),0) AS amount").Scan(&amount)
	if result.Error != nil {
		return decimal.Zero, errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
	}
	return
}

func (r *cashFlowRepo) GetCashFlowWithDaily(groups []string, startTime, endTime time.Time, paymentType int8) (amount decimal.Decimal, err error) {
	tx, joinDailyErr := TradeRepo.joinDaily(endTime)
	if joinDailyErr != nil {
		return decimal.Zero, errors.WithMessage(joinDailyErr, util.ErrWrapMessage())
	}

	tx.Where("r.group IN ?", groups).
		Where("t.cmd = ?", 6)

	if !startTime.IsZero() {
		tx.Where("t.open_time >= ?", startTime)
	}

	if !endTime.IsZero() {
		tx.Where("t.open_time <= ?", endTime)
	}

	if paymentType == 1 {
		tx.Where("t.profit > 0")
	} else if paymentType == 2 {
		tx.Where("t.profit < 0")
	}

	result := tx.Select("COALESCE(SUM(t.profit),0) AS amount").Scan(&amount)
	if result.Error != nil {
		return decimal.Zero, errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
	}
	return
}
