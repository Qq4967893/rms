package repo

import (
	"backend-rms/core/database"
	"backend-rms/model"
	"backend-rms/util"
	"time"

	"github.com/pkg/errors"
)

type copyRecordRepo struct{}

var CopyRecordRepo = &copyRecordRepo{}

func (r *copyRecordRepo) GetAll(followerId, masterId int64) (copiesRecords []model.CopyRecord, err error) {
	result := database.APIDB.Table("copy_records").
		Where("mt4_to", followerId).
		Where("mt4_from", masterId).
		Find(&copiesRecords)
	return copiesRecords, errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
}

func (r *copyRecordRepo) GetCopies(masterId int64, groups []string) (followerList []model.CopiesRecord, err error) {
	result := database.APIDB.Table("copy_records as c").
		Joins("left join trade_account as t on t.trade_account_id = c.mt4_to").
		Where("t.group IN ?", groups).
		Where("c.mt4_from = ?", masterId).
		Group("c.mt4_to").
		Select("c.mt4_to, ANY_VALUE(t.user_code) as user_code, ANY_VALUE(t.group) as `group`").
		Find(&followerList)
	return followerList, errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
}

// 获取跟随者们跟随的高手们
func (r *copyRecordRepo) GetCopiesMasters(dateTime time.Time, followers []int64, groups []string) (copiesMasters []model.CopiesMasters, err error) {
	result := database.APIDB.Table("copy_records as c").
		Joins("left join trade_account as t on t.trade_account_id = c.mt4_to").
		Where("t.group IN ?", groups).
		Where("c.mt4_to IN ?", followers).
		Where("c.start_time < ?", dateTime).Where("c.end_time >= ? OR c.end_time = ?", dateTime, model.Copy_Record_Default_EndTime).
		Group("c.mt4_to").
		Group("c.mt4_from").
		Select("c.mt4_to, c.mt4_from, sum(c.amount) as amount").
		Find(&copiesMasters)
	return copiesMasters, errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
}

// 获取正在跟随的账号
func (r *copyRecordRepo) GetFollowingLogins(masterId int64, groups []string) ([]int64, error) {
	followingLogins := make([]struct {
		Mt4To int64 `json:"mt4_to"`
	}, 0)

	result := database.APIDB.Table("copy_records as c").
		Joins("left join trade_account as t on t.trade_account_id = c.mt4_to").
		Where("t.group IN ?", groups).
		Where("c.mt4_from", masterId).
		Where("c.end_time", model.Copy_Record_Default_EndTime).
		Select("c.mt4_to").
		Find(&followingLogins)
	if result.Error != nil {
		return nil, errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
	}

	var logins []int64
	for _, v := range followingLogins {
		logins = append(logins, v.Mt4To)
	}
	return logins, nil
}
