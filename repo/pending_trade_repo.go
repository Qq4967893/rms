package repo

import (
	"backend-rms/core/database"
	"backend-rms/model"
	"backend-rms/util"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type pendingTradeRepo struct{}

var PendingTradeRepo = &pendingTradeRepo{}

func (r *pendingTradeRepo) GetPendingTrades(groups []string, symbol string, cmd int8) (pendingTrades []model.PendingTrade, err error) {
	tx := r.joinAccount().Where("a.group IN ?", groups).
		Where("t.close_time = ?", model.DEFAULT_CLOSE_TIME).
		Where("t.cmd IN ?", model.GetPendingCmd(cmd))

	if symbol != "" {
		tx.Where("t.symbol LIKE ?", symbol+"%")
	}
	result := tx.Find(&pendingTrades)
	return pendingTrades, errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
}

func (r *pendingTradeRepo) joinAccount() (tx *gorm.DB) {
	tx = database.ET6DB.Table("trade_trans_result as t").
		Joins("left join account as a on t.login = a.login")
	return
}
