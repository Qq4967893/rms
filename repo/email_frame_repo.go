package repo

import (
	"backend-rms/core/database"
	"backend-rms/model"
	"backend-rms/util"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type emailFrameRepo struct{}

var EmailFrameRepo = &emailFrameRepo{}

func (r *emailFrameRepo) GetInfoById(id int64) (emailFrame *model.EmailFrame, err error) {
	err = database.APIDB.Table("email_frame").Where("id", id).First(&emailFrame).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return emailFrame, errors.Wrapf(err, "%s failed", util.GetCurrentFuncName())
}
