package repo

import (
	"backend-rms/core/database"
	"backend-rms/model"
	"backend-rms/util"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type roleRepo struct {
}

var RoleRepo = &roleRepo{}

func (r *roleRepo) GetInfoById(id int64) (role *model.Role, err error) {
	err = database.RMSDB.Where("id", id).First(&role).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return role, errors.Wrapf(err, "%s failed", util.GetCurrentFuncName())
}

func (r *roleRepo) GetList() (roles []*model.Role, err error) {
	result := database.RMSDB.Find(&roles)
	return roles, errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
}

func (r *roleRepo) Add(role *model.Role) (id int64, err error) {
	result := database.RMSDB.Create(&role)
	return role.ID, errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
}

func (r *roleRepo) Update(id int64, role model.Role) error {
	result := database.RMSDB.Where("id", id).Updates(&role)
	return errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
}

func (r *roleRepo) Delete(id int64) (err error) {
	result := database.RMSDB.Delete(&model.Role{}, id)
	return errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
}
