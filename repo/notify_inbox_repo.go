package repo

import (
	"backend-rms/core/database"
	"backend-rms/model"
	"backend-rms/util"

	"github.com/pkg/errors"
)

type notifyInboxRepo struct{}

var NotifyInboxRepo = &notifyInboxRepo{}

func (r *notifyInboxRepo) Add(params model.NotifyInbox) error {
	result := database.APIDB.Table("notify_inbox").Create(&params)
	return errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
}
