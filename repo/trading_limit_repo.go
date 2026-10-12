package repo

import (
	"backend-rms/core/database"
	"backend-rms/model"
	"backend-rms/util"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type tradingLimitRepo struct{}

var TradingLimitRepo = &tradingLimitRepo{}

func (r *tradingLimitRepo) GetList(passport int64, tradingLimit int8, paging *model.Paging) (tradingLimits []*model.TradingLimitList, err error) {
	tx := database.APIDB.Table("trading_limits")

	if passport != 0 {
		tx.Where("passport", passport)
	}

	if tradingLimit == 1 {
		tx.Where("remove_time", model.TRADING_REMOVE_TIME)
	} else {
		tx.Where("remove_time != ?", model.TRADING_REMOVE_TIME)
	}

	if paging != nil {
		tx.Limit(paging.Limit).Offset(paging.Offset)
	}
	result := tx.Distinct("id, passport, symbol, trigger_value, limit_value, current_value, remove_value, trigger_time, remove_time, operator").Find(&tradingLimits)
	return tradingLimits, errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
}

func (r *tradingLimitRepo) GetInfoById(id int64) (tradingLimitInfo *model.TradingLimitInfo, err error) {
	err = database.APIDB.Table("trading_limits as t").
		Joins("left join user as u on t.passport = u.user_code").
		Where("t.id", id).
		Select("u.email, u.username, u.real_name, u.world_code, t.passport, t.current_value, t.limit_value, t.symbol, t.remove_time").
		First(&tradingLimitInfo).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return tradingLimitInfo, errors.Wrapf(err, "%s failed", util.GetCurrentFuncName())
}

func (r *tradingLimitRepo) Update(id int64, params model.TradingLimit) error {
	result := database.APIDB.Where("id", id).Updates(&params)
	return errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
}
