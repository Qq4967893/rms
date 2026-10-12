package repo

import (
	"backend-rms/core/database"
	"backend-rms/model"
	"backend-rms/util"
	"time"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type dailyClosedTradeRepo struct {
}

var DailyClosedTradeRepo = &dailyClosedTradeRepo{}

func (dr *dailyClosedTradeRepo) GetSymbolList(groups []string, symbol string, startDate, endDate time.Time) (symbolList []model.SymbolDailyClosedTrade, err error) {
	tx := dr.joinTrade().Where("t.trans_close_group IN ?", groups).
		Where("d.close_date >= ?", startDate).
		Where("d.close_date <= ?", endDate)

	if symbol != "" {
		tx.Where("t.symbol LIKE ?", symbol+"%")
	}

	tx.Select("t.symbol, sum(t.profit) as profit, sum(t.swaps) as swaps, sum(t.commission) as commission, sum(d.spread_profit - d.without_raw_spread_profit) as raw_spread_profit, sum(d.without_raw_spread_profit) as mark_up_spread_profit, (sum(t.volume) / 100) as lots").
		Group("t.symbol")

	result := tx.Find(&symbolList)
	return symbolList, errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
}

func (dr *dailyClosedTradeRepo) GetAccountList(groups []string, symbol string, startDate, endDate time.Time, paging *model.Paging) (accountList []model.AccountDailyClosedTrade, err error) {
	tx := dr.joinTrade().Where("t.trans_close_group IN ?", groups).
		Where("d.close_date >= ?", startDate).
		Where("d.close_date <= ?", endDate)

	if symbol != "" {
		tx.Where("t.symbol LIKE ?", symbol+"%")
	}

	tx.Select("t.login, sum(t.profit) as profit, sum(t.swaps) as swaps, sum(t.commission) as commission, sum(d.spread_profit - d.without_raw_spread_profit) as raw_spread_profit, sum(d.without_raw_spread_profit) as mark_up_spread_profit, (sum(t.volume) / 100) as lots").
		Group("t.login")

	if paging != nil {
		tx.Order(paging.Sort + " " + paging.Order).
			Limit(paging.Limit).
			Offset(paging.Offset)
	}
	result := tx.Find(&accountList)
	return accountList, errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
}

func (dr *dailyClosedTradeRepo) GetAccountCount(groups []string, symbol string, startDate, endDate time.Time) (totalCount int64, err error) {
	tx := dr.joinTrade().Where("t.trans_close_group IN ?", groups).
		Where("d.close_date >= ?", startDate).
		Where("d.close_date <= ?", endDate)

	if symbol != "" {
		tx.Where("t.symbol LIKE ?", symbol+"%")
	}

	result := tx.Group("t.login").Count(&totalCount)
	return totalCount, errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
}

func (dr *dailyClosedTradeRepo) GetTradeList(groups []string, symbol string, startDate, endDate time.Time, paging *model.Paging) (tradeList []model.DailyClosedTradeInfo, err error) {
	tx := dr.joinTrade().Where("t.trans_close_group IN ?", groups).
		Where("d.close_date >= ?", startDate).
		Where("d.close_date <= ?", endDate)

	if symbol != "" {
		tx.Where("t.symbol LIKE ?", symbol+"%")
	}

	tx.Select("t.login, t.symbol, d.ticket, t.profit, t.swaps, t.commission, (d.spread_profit - d.without_raw_spread_profit) as raw_spread_profit, d.without_raw_spread_profit as mark_up_spread_profit, (t.volume / 100) as lots, t.open_price, t.close_price, t.open_time, t.close_time")

	if paging != nil {
		tx.Order(paging.Sort + " " + paging.Order).
			Limit(paging.Limit).
			Offset(paging.Offset)
	}
	result := tx.Find(&tradeList)
	return tradeList, errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
}

func (dr *dailyClosedTradeRepo) GetTradeCount(groups []string, symbol string, startDate, endDate time.Time) (totalCount int64, err error) {
	tx := dr.joinTrade().Where("t.trans_close_group IN ?", groups).
		Where("d.close_date >= ?", startDate).
		Where("d.close_date <= ?", endDate)

	if symbol != "" {
		tx.Where("t.symbol LIKE ?", symbol+"%")
	}

	result := tx.Count(&totalCount)
	return totalCount, errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
}

func (dr *dailyClosedTradeRepo) Add(params []model.DailyClosedTrade) error {
	result := database.RMSDB.CreateInBatches(&params, model.BatchSize)
	return result.Error
}

func (dr *dailyClosedTradeRepo) joinTrade() (tx *gorm.DB) {
	tx = database.RMSDB.Table("daily_closed_trades as d").
		Joins("left join trading_system_user.trade_trans_result as t on d.ticket = t.ticket")
	return
}
