package repo

import (
	"backend-rms/core/database"
	"backend-rms/model"
	"backend-rms/util"

	"github.com/pkg/errors"
)

type roleMenuRepo struct {
}

var RoleMenuRepo = &roleMenuRepo{}

func (r *roleMenuRepo) GetAllByRoleId(roleId int64) (roleMenus []*model.RoleMenu, err error) {
	result := database.RMSDB.Where("role_id", roleId).Find(&roleMenus)
	return roleMenus, errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
}

func (r *roleMenuRepo) GetAllByMenuId(menuId int64) (roleMenus []*model.RoleMenu, err error) {
	result := database.RMSDB.Where("menu_id", menuId).Find(&roleMenus)
	return roleMenus, errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
}

func (r *roleMenuRepo) AddBatches(roleMenus []model.RoleMenu) (err error) {
	result := database.RMSDB.Create(&roleMenus)
	return errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
}

func (r *roleMenuRepo) Delete(roleId int64, menuId int64) (err error) {
	tx := database.RMSDB.Model(model.RoleMenu{})
	if roleId != 0 {
		tx.Where("role_id", roleId)
	}
	if menuId != 0 {
		tx.Where("menu_id", roleId)
	}
	result := tx.Delete(&model.RoleMenu{})
	return errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
}
