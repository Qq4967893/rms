package repo

import (
	"backend-rms/core/database"
	"backend-rms/model"
	"backend-rms/util"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type userRepo struct{}

var UserRepo = &userRepo{}

func (r *userRepo) GetInfoByPassport(passport int64) (user *model.User, err error) {
	err = database.APIDB.Table("user").
		Select("id, user_code as passport, username, avatar, world_code as country, language, area_id, status, user_affiliation as affiliation_id, type, create_date as created_at").
		Where("user_code", passport).First(&user).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return user, errors.Wrapf(err, "%s failed", util.GetCurrentFuncName())
}
