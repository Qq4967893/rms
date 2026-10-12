package repo

import (
	"backend-rms/core/database"
	"backend-rms/model"
	"backend-rms/util"
	"fmt"
	"time"

	"github.com/golang-module/carbon/v2"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type tradeRepo struct{}

var TradeRepo = &tradeRepo{}

func (r *tradeRepo) GetTradeListByLogins(logins []int64, ticket int64, symbol string) (trades []*model.Trade, err error) {
	tx := database.ET6DB.Where("login IN ?", logins).
		Where("cmd < ?", 2)
	if ticket > 0 {
		tx.Where("ticket", ticket)
	}
	if symbol != "" {
		tx.Where("symbol LIKE ?", symbol+"%")
	}
	result := tx.Find(&trades)
	return trades, errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
}

func (r *tradeRepo) GetOpeningTrades(groups []string, login int64, symbol string, cmd int8) (trades []model.OpeningTrade, err error) {
	tx := r.joinAccount().Where("t.close_time = ?", model.DEFAULT_CLOSE_TIME)

	if groups != nil {
		tx.Where("a.group IN ?", groups)
	}

	if cmd == -1 {
		tx.Where("t.cmd < ?", 2)
	} else {
		tx.Where("t.cmd = ?", cmd)
	}

	if symbol != "" {
		tx.Where("t.symbol LIKE ?", symbol+"%")
	}

	if login != 0 {
		tx.Where("a.login", login)
	}

	result := tx.Select("t.ticket, t.symbol, t.cmd, t.volume, t.open_price, t.commission, t.swaps, t.login, t.open_time, t.sl, t.tp, a.group").
		Find(&trades)
	return trades, errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
}

func (r *tradeRepo) GetPointOpeningTradeWithAccount(groups []string, logins []int64, dateTime time.Time, symbol string) (trades []model.OpeningTrade, err error) {
	tx := r.joinAccount().
		Where("t.cmd < ?", 2).
		Where("t.open_time <= ?", dateTime).
		Where("t.close_time >= ? OR t.close_time = ?", dateTime, model.DEFAULT_CLOSE_TIME).
		Select("t.ticket, t.symbol, t.cmd, t.volume, t.open_price, t.commission, t.swaps, t.login, a.group")
	if len(groups) > 0 {
		tx.Where("a.group IN ?", groups)
	}
	if len(logins) > 0 {
		tx.Where("t.login IN ?", logins)
	}
	if symbol != "" {
		tx.Where("t.symbol like ?", symbol+"%")
	}
	result := tx.Find(&trades)
	return trades, errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
}

func (r *tradeRepo) GetPointOpeningTradeWithDaily(groups []string, logins []int64, dateTime time.Time, symbol string) (trades []model.OpeningTrade, err error) {
	tx, joinDailyErr := r.joinDaily(dateTime)
	if joinDailyErr != nil {
		return nil, errors.WithMessage(joinDailyErr, util.ErrWrapMessage())
	}
	tx.Where("t.cmd < ?", 2).
		Where("t.open_time <= ?", dateTime).
		Where("t.close_time >= ? OR t.close_time = ?", dateTime, model.DEFAULT_CLOSE_TIME).
		Select("t.ticket, t.symbol, t.cmd, t.volume, t.open_price, t.commission, t.swaps, t.login, r.group")
	if len(groups) > 0 {
		tx.Where("r.group IN ?", groups)
	}
	if len(logins) > 0 {
		tx.Where("t.login IN ?", logins)
	}
	if symbol != "" {
		tx.Where("t.symbol like ?", symbol+"%")
	}
	result := tx.Find(&trades)
	return trades, errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
}

func (r *tradeRepo) GetReportTradeList(tradeStatus, symbol string, groups []string, startTime, endTime time.Time, login []int64, ticket int64, paging model.Paging) (trades []*model.ReportTrade, err error) {
	tx, tableAlias, err := r.getTxByTime(endTime)
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}

	fields := r.getFieldsByTableAlias(tableAlias)
	tx = r.buildReportTradeTx(tx, tableAlias, tradeStatus, symbol, groups, startTime, endTime, login, ticket)
	result := tx.Select(fields).
		Order("t." + paging.Sort + " " + paging.Order).Order("t.ticket desc").
		Limit(paging.Limit).Offset(paging.Offset).Find(&trades)

	return trades, errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
}

func (r *tradeRepo) GetReportTradeCount(tradeStatus, symbol string, groups []string, startTime, endTime time.Time, login []int64, ticket int64) (totalCount int64, err error) {
	tx, tableAlias, err := r.getTxByTime(endTime)
	if err != nil {
		return 0, errors.WithMessage(err, util.ErrWrapMessage())
	}

	tx = r.buildReportTradeTx(tx, tableAlias, tradeStatus, symbol, groups, startTime, endTime, login, ticket)
	result := tx.Count(&totalCount)

	return totalCount, errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
}

func (r *tradeRepo) GetReportTradeListByTickets(tickets []int64, groups []string) (trades []*model.ReportTrade, err error) {
	if len(tickets) == 0 {
		return nil, errors.New("tickets can not be empty")
	}

	tx, tableAlias, err := r.getTxByTime(time.Now())
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}

	tx.Where("t.cmd < 2").Where("t.ticket IN ?", tickets)
	if len(groups) > 0 {
		tx.Where(tableAlias+".group IN ?", groups)
	}

	fields := r.getFieldsByTableAlias(tableAlias)
	result := tx.Select(fields).Find(&trades)

	return trades, errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
}

func (r *tradeRepo) GetAccountClosedPL(logins []int64) (accountClosedPL []model.AccountClosedPL, err error) {
	tx := r.joinAccount().Scopes(r.closedTrade)
	if len(logins) > 0 {
		tx.Where("t.login IN ?", logins)
	}

	result := tx.Select("t.login, COALESCE(SUM(t.profit + t.swaps + t.commission),0) AS ClosedPL").
		Group("t.login").Find(&accountClosedPL)
	return accountClosedPL, errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
}

func (r *tradeRepo) GetClosedPLWithAccount(groups []string, startTime, endTime time.Time) (closedPL decimal.Decimal, err error) {
	tx := r.joinAccount().Where("a.group IN ?", groups)
	tx.Scopes(r.closedTrade)
	if !startTime.IsZero() {
		tx.Where("t.close_time >= ?", startTime)
	}

	if !endTime.IsZero() {
		tx.Where("t.close_time <= ?", endTime)
	}

	result := tx.Select("COALESCE(SUM(t.profit + t.swaps + t.commission),0) AS closedPL").Scan(&closedPL)
	return closedPL, errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
}

func (r *tradeRepo) GetClosedTrades(groups []string, startTime, endTime time.Time) (trades []model.ClosedTrade, err error) {
	tx := r.joinAccount().Where("a.group IN ?", groups).
		Where("t.cmd < ?", 2)

	if endTime.IsZero() && startTime.IsZero() {
		tx.Where("t.close_time != ?", model.DEFAULT_CLOSE_TIME)
	} else {
		tx.Where("t.close_time > ?", startTime).Where("t.close_time <= ?", endTime)
	}

	result := tx.Select("t.ticket, t.login, t.symbol, t.cmd, t.volume, t.open_price, t.close_price, t.open_time, t.close_time, t.commission, t.swaps, t.login, a.group, t.trans_open_price, t.trans_close_price, t.trans_open_group, t.trans_close_group, t.deposit_currency").
		Find(&trades)
	return trades, errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
}

func (r *tradeRepo) GetClosedPLWithDaily(groups []string, startTime, endTime time.Time) (closedPL decimal.Decimal, err error) {
	tx, joinDailyErr := r.joinDaily(endTime)
	if joinDailyErr != nil {
		return decimal.Zero, errors.WithMessage(joinDailyErr, util.ErrWrapMessage())
	}

	tx.Where("r.group IN ?", groups)
	tx.Scopes(r.closedTrade)
	if !startTime.IsZero() {
		tx.Where("t.close_time >= ?", startTime)
	}

	if !endTime.IsZero() {
		tx.Where("t.close_time <= ?", endTime)
	}

	result := tx.Select("COALESCE(SUM(t.profit + t.swaps + t.commission),0) AS closedPL").Scan(&closedPL)
	return closedPL, errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
}

func (r *tradeRepo) GetAccountClosedPLWithAccount(logins []int64, groups []string, startTime, endTime time.Time, symbol string) (accountClosedPL []model.AccountClosedPL, err error) {
	tx := r.joinAccount()
	tx.Scopes(r.closedTrade)
	if len(logins) > 0 {
		tx.Where("t.login IN ?", logins)
	}

	if len(groups) > 0 {
		tx.Where("a.group IN ?", groups)
	}

	if !startTime.IsZero() {
		tx.Where("t.close_time >= ?", startTime)
	}

	if !endTime.IsZero() {
		tx.Where("t.close_time <= ?", endTime)
	}

	if symbol != "" {
		tx.Where("t.symbol like ?", symbol+"%")
	}

	result := tx.Select("t.login, COALESCE(SUM(t.profit + t.swaps + t.commission),0) AS ClosedPL").
		Group("t.login").Find(&accountClosedPL)
	return accountClosedPL, errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
}

func (r *tradeRepo) GetAccountClosedPLWithDaily(logins []int64, groups []string, startTime, endTime time.Time, symbol string) (accountClosedPL []model.AccountClosedPL, err error) {
	tx, joinDailyErr := r.joinDaily(endTime)
	if joinDailyErr != nil {
		return nil, errors.WithMessage(joinDailyErr, util.ErrWrapMessage())
	}

	tx.Scopes(r.closedTrade)
	if len(logins) > 0 {
		tx.Where("t.login IN ?", logins)
	}

	if len(groups) > 0 {
		tx.Where("r.group IN ?", groups)
	}

	if !startTime.IsZero() {
		tx.Where("t.close_time >= ?", startTime)
	}

	if !endTime.IsZero() {
		tx.Where("t.close_time <= ?", endTime)
	}

	if symbol != "" {
		tx.Where("t.symbol like ?", symbol+"%")
	}

	result := tx.Select("t.login, COALESCE(SUM(t.profit + t.swaps + t.commission),0) AS ClosedPL").
		Group("t.login").Find(&accountClosedPL)
	return accountClosedPL, errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
}

func (r *tradeRepo) GetSymbolClosedPLWithAccount(groups []string, startTime, endTime time.Time) (SymbolClosedPL []model.SymbolClosedPL, err error) {
	tx := r.joinAccount().Debug().Where("a.group IN ?", groups)
	tx.Scopes(r.closedTrade)
	if !startTime.IsZero() {
		tx.Where("t.close_time >= ?", startTime)
	}

	if !endTime.IsZero() {
		tx.Where("t.close_time <= ?", endTime)
	}
	result := tx.Select("t.symbol, sum(t.profit) as profit, sum(t.swaps) as swaps, sum(t.commission) as commission").
		Group("t.symbol").Find(&SymbolClosedPL)
	return SymbolClosedPL, errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
}

func (r *tradeRepo) GetSymbolClosedPLWithDaily(groups []string, startTime, endTime time.Time) (SymbolClosedPL []model.SymbolClosedPL, err error) {
	tx, joinDailyErr := r.joinDaily(endTime)
	if joinDailyErr != nil {
		return nil, errors.WithMessage(joinDailyErr, util.ErrWrapMessage())
	}

	tx.Where("r.group IN ?", groups)
	tx.Scopes(r.closedTrade)
	if !startTime.IsZero() {
		tx.Where("t.close_time >= ?", startTime)
	}

	if !endTime.IsZero() {
		tx.Where("t.close_time <= ?", endTime)
	}
	result := tx.Select("t.symbol, sum(t.profit) as profit, sum(t.swaps) as swaps, sum(t.commission) as commission").
		Group("t.symbol").Find(&SymbolClosedPL)
	return SymbolClosedPL, errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
}

func (r *tradeRepo) GetClosedPL(login int64, startTime, endTime time.Time) (PL decimal.Decimal, err error) {
	tx := database.ET6DB.Model(model.Trade{}).
		Where("cmd < ?", 2).
		Where("login = ?", login)
	if startTime.IsZero() {
		tx.Where("close_time > ?", model.DEFAULT_CLOSE_TIME)
	} else {
		tx.Where("close_time >= ?", startTime)
	}

	if !endTime.IsZero() {
		tx.Where("close_time <= ?", endTime)
	}

	result := tx.Select("COALESCE(SUM(profit+swaps+commission), 0) as pl").Scan(&PL)
	if result.Error != nil {
		return decimal.Zero, errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
	}
	return
}

func (r *tradeRepo) GetClosedPLMap(logins []int64, startTime, endTime string) (map[int64]decimal.Decimal, error) {
	tx := database.ET6DB.Model(model.Trade{}).Where("cmd < ?", 2)
	if len(logins) > 0 {
		tx.Where("login IN ?", logins)
	}

	if startTime == "" {
		tx.Where("close_time > ?", model.DEFAULT_CLOSE_TIME)
	} else {
		tx.Where("close_time >= ?", startTime)
	}

	if endTime != "" {
		tx.Where("close_time <= ?", endTime)
	}

	type Row struct {
		Login int64
		PL    decimal.Decimal
	}

	var rows []*Row
	result := tx.Select("login, sum(profit+swaps+commission) as pl").Group("login").Find(&rows)
	if result.Error != nil {
		return nil, errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
	}

	closedPLMap := make(map[int64]decimal.Decimal, len(rows))
	for _, v := range rows {
		closedPLMap[v.Login] = v.PL
	}

	return closedPLMap, nil
}

func (r *tradeRepo) closedTrade(db *gorm.DB) *gorm.DB {
	return db.Where("cmd < ?", 2).
		Where("close_time != ?", model.DEFAULT_CLOSE_TIME)
}

func (r *tradeRepo) joinAccount() (tx *gorm.DB) {
	tx = database.ET6DB.Table("trade_trans_result as t").
		Joins("left join account as a on t.login = a.login")
	return
}

func (r *tradeRepo) joinDaily(dateTime time.Time) (tx *gorm.DB, err error) {
	dailyDateTime := model.GetDailyDateTime(dateTime)
	tx = database.ET6DB.Table("trade_trans_result as t").
		Joins("left join report_daily as r on t.login = r.login").
		Where("r.time = ?", dailyDateTime)
	return
}

func (r *tradeRepo) buildReportTradeTx(db *gorm.DB, tableAlias, tradeStatus, symbol string, groups []string, startTime, endTime time.Time, login []int64, ticket int64) (tx *gorm.DB) {
	switch tradeStatus {
	case "open":
		tx = db.Where("t.open_time >= ? AND t.open_time <= ?", startTime, endTime)
	case "closed":
		tx = db.Where("t.close_time >= ? AND t.close_time <= ?", startTime, endTime)
	case "holding":
		tx = db.Where("(t.open_time <= ? AND t.close_time >= ?) OR (t.open_time <= ? AND t.close_time <= ?)",
			startTime, endTime, startTime, model.DEFAULT_CLOSE_TIME)
	default:
		tx = db.Where("(t.open_time >= ? AND t.open_time <= ?) OR (t.close_time >= ? AND t.close_time <= ?) OR ((t.open_time <= ? AND t.close_time >= ?) OR (t.open_time <= ? AND t.close_time <= ?))",
			startTime, endTime, startTime, endTime, startTime, endTime, startTime, model.DEFAULT_CLOSE_TIME)
	}

	if len(login) > 0 {
		if len(login) == 1 {
			tx.Where("t.login = ?", login)
		} else {
			tx.Where("t.login IN ?", login)
		}
	} else if ticket > 0 {
		tx.Where("t.ticket = ?", ticket)
	}

	if symbol != "" {
		tx.Where("t.symbol LIKE ?", symbol+"%")
	}

	if len(groups) > 0 {
		tx.Where(tableAlias+".group IN ?", groups)
	}

	tx.Where("t.cmd < 2")
	return
}

func (r *tradeRepo) getTxByTime(dateTime time.Time) (tx *gorm.DB, tableAlias string, err error) {
	today := carbon.Now().StartOfDay().Carbon2Time()
	if dateTime.Before(today) {
		tx, err = r.joinDaily(dateTime)
		tableAlias = "r"
	} else {
		tx = r.joinAccount()
		tableAlias = "a"
	}
	return tx, tableAlias, errors.WithMessage(err, util.ErrWrapMessage())
}

func (r *tradeRepo) getFieldsByTableAlias(tableAlias string) string {
	return fmt.Sprintf("t.type, t.ticket, t.login, t.symbol,"+
		"t.cmd, t.volume, t.sl, t.tp, t.profit, t.swaps, t.commission,"+
		"t.open_price, t.open_time, t.close_price, t.close_time, %s.group",
		tableAlias)
}
