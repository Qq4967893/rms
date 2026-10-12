package repo

import (
	"backend-rms/core/database"
	"backend-rms/model"
	"backend-rms/util"

	"github.com/pkg/errors"
)

type notificationRepo struct{}

var NotificationRepo = &notificationRepo{}

func (r *notificationRepo) Add(params model.Notification) (int64, error) {
	result := database.APIDB.Model(model.Notification{}).Create(&params)
	return params.ID, errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
}
