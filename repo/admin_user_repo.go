package repo

import (
	"backend-rms/core/database"
	"backend-rms/util"

	"github.com/pkg/errors"
)

type adminUserRepo struct{}

var AdminUserRepo = &adminUserRepo{}

func (r *adminUserRepo) GetUsernameByID(id int64) (username string, err error) {
	err = database.CRMDB.Table("users").
		Select("user_name").
		Where("id", id).Scan(&username).Error
	return username, errors.Wrapf(err, "%s failed", util.GetCurrentFuncName())
}
