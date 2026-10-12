package repo

import (
	"backend-rms/core/database"
	"backend-rms/model"
	"backend-rms/util"
	"time"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type masterRepo struct{}

var MasterRepo = &masterRepo{}

func (r *masterRepo) GetMasterList(groups []string, login int64, startTime, endTime time.Time) (masters []*model.Master, err error) {
	tx := database.APIDB.Table("master_record").
		Select("master_record.mt4_id as login, trade_account.user_code as passport, trade_account.group").
		Joins("left join trade_account on trade_account.trade_account_id = master_record.mt4_id").
		Where("master_record.created_at between ? and ?", startTime, endTime).
		Where("master_record.name = ?", "高手榜").
		Where("master_record.type IN ?", []int{0, 1}).
		Where("trade_account.trade_system_type = 1")
	if len(groups) > 0 {
		tx.Where("trade_account.`group` IN ?", groups)
	}
	if login > 0 {
		tx.Where("master_record.mt4_id = ?", login)
	}

	err = errors.Wrapf(tx.Find(&masters).Error, "%s failed", util.GetCurrentFuncName())
	return
}

func (r *masterRepo) GetMasterHistoricalList(login int64) (masterRecords []*model.MasterRecord, err error) {
	tx := database.APIDB.Table("master_record").
		Select("master_record.id, master_record.mt4_id as login, master_record.type, master_record.created_at").
		Joins("left join trade_account on trade_account.trade_account_id = master_record.mt4_id").
		Where("master_record.name = ?", "高手榜").
		Where("master_record.type IN ?", []int{0, 1}).
		Where("trade_account.trade_system_type = 1").
		Where("master_record.mt4_id = ?", login).
		Order("master_record.created_at asc")

	err = errors.Wrapf(tx.Find(&masterRecords).Error, "%s failed", util.GetCurrentFuncName())
	return
}

func (r *masterRepo) GetCurrentMasterList(groups []string, login int64) (masters []*model.Master, err error) {
	tx := database.APIDB.Table("master_new").
		Select("master_new.mt4_id as login, trade_account.user_code as passport, trade_account.group").
		Joins("left join trade_account on trade_account.trade_account_id = master_new.mt4_id").
		Where("trade_account.trade_system_type = 1")
	if len(groups) > 0 {
		tx.Where("trade_account.`group` IN ?", groups)
	}
	if login > 0 {
		tx.Where("master_new.mt4_id = ?", login)
	}

	err = errors.Wrapf(tx.Find(&masters).Error, "%s failed", util.GetCurrentFuncName())
	return
}

func (r *masterRepo) GetCopiedAmountMap(groups []string, logins []int64, dateTime time.Time) (map[int64]decimal.Decimal, error) {
	tx := database.APIDB.Table("copy_records").
		Select("copy_records.mt4_from as login, sum(copy_records.amount) as amount").
		Joins("left join trade_account on trade_account.trade_account_id = copy_records.mt4_to").
		Where("trade_account.trade_system_type = 1").
		Where("copy_records.start_time < ?", dateTime).
		Where("copy_records.end_time >= ? or copy_records.end_time = '0000-00-00 00:00:00'", dateTime).
		Group("copy_records.mt4_from")
	if len(groups) > 0 {
		tx.Where("trade_account.`group` IN ?", groups)
	}
	if len(logins) > 0 {
		tx.Where("copy_records.mt4_from IN ?", logins)
	}

	type Row struct {
		Login  int64
		Amount decimal.Decimal
	}

	var rows []*Row
	result := tx.Find(&rows)
	if result.Error != nil {
		return nil, errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
	}

	copiedAmountMap := make(map[int64]decimal.Decimal, len(rows))
	for _, v := range rows {
		copiedAmountMap[v.Login] = v.Amount
	}

	return copiedAmountMap, nil
}

// 获取高手当前时刻正在跟随的跟随者数
func (r *masterRepo) GetFollowerCountMap(groups []string, logins []int64) (map[int64]int64, error) {
	tx := database.APIDB.Table("copy_records").
		Select("copy_records.mt4_from as login, count(1) as count").
		Joins("left join trade_account on trade_account.trade_account_id = copy_records.mt4_to").
		Where("trade_account.trade_system_type = 1").
		Where("copy_records.end_time = ?", model.Copy_Record_Default_EndTime).
		Group("copy_records.mt4_from, copy_records.mt4_to")
	if len(groups) > 0 {
		tx.Where("trade_account.`group` IN ?", groups)
	}
	if len(logins) > 0 {
		tx.Where("copy_records.mt4_from IN ?", logins)
	}

	type Row struct {
		Login int64
		Count int64
	}

	var rows []*Row
	result := tx.Find(&rows)
	if result.Error != nil {
		return nil, errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
	}

	followerCountMap := make(map[int64]int64, len(rows))
	for _, v := range rows {
		if _, ok := followerCountMap[v.Login]; ok {
			followerCountMap[v.Login] = followerCountMap[v.Login] + v.Count
			continue
		}
		followerCountMap[v.Login] = v.Count
	}

	return followerCountMap, nil
}

func (r *masterRepo) GetCommission(login int64, startTime, endTime time.Time) (commission decimal.Decimal, err error) {
	tx := database.APIDB.Table("copy_pay_master").
		Where("mt4_id", login).
		Where("pay_time between ? and ?", startTime, endTime)
	result := tx.Select("COALESCE(SUM(sum_pay), 0) as pl").Scan(&commission)
	if result.Error != nil {
		return decimal.Zero, errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
	}
	return
}

func (r *masterRepo) GetFollowerCount(groups []string, login int64, startTime, endTime time.Time) (followerCount int64, err error) {
	tx := database.APIDB.Table("copy_records").
		Joins("left join trade_account on trade_account.trade_account_id = copy_records.mt4_to").
		Where("trade_account.trade_system_type = 1").
		Where("copy_records.mt4_from = ?", login).
		Where("copy_records.start_time >= ?", startTime).
		Where("copy_records.end_time <= ?", endTime).
		Group("copy_records.mt4_to")
	if len(groups) > 0 {
		tx.Where("trade_account.`group` IN ?", groups)
	}
	result := tx.Select("COALESCE(COUNT(1), 0) as count").Scan(&followerCount)
	if result.Error != nil {
		return 0, errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
	}
	return
}

func (r *masterRepo) GetFollowedTimes(groups []string, login int64, startTime, endTime time.Time) (followedTimes int64, err error) {
	tx := database.APIDB.Table("copy_records").
		Joins("left join trade_account on trade_account.trade_account_id = copy_records.mt4_to").
		Where("trade_account.trade_system_type = 1").
		Where("copy_records.mt4_from = ?", login).
		Where("copy_records.start_time >= ?", startTime).
		Where("copy_records.end_time <= ?", endTime)
	if len(groups) > 0 {
		tx.Where("trade_account.`group` IN ?", groups)
	}
	result := tx.Select("COALESCE(COUNT(1), 0) as count").Scan(&followedTimes)
	if result.Error != nil {
		return 0, errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
	}
	return
}
