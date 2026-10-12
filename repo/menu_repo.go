package repo

import (
	"backend-rms/core/database"
	"backend-rms/model"
	"backend-rms/util"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type menuRepo struct{}

var MenuRepo = &menuRepo{}

func (r *menuRepo) GetInfoById(id int64) (menu *model.Menu, err error) {
	err = database.RMSDB.Where("id", id).First(&menu).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return menu, errors.Wrapf(err, "%s failed", util.GetCurrentFuncName())
}

func (r *menuRepo) GetInfoByPid(pid int64) (menus []model.Menu, err error) {
	result := database.RMSDB.Where("pid", pid).Find(&menus)
	return menus, errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
}

func (r *menuRepo) GetAllByIds(ids []int, display int8) (menus []model.Menu, err error) {
	tx := database.RMSDB.Where("id IN ?", ids)
	if display == 1 {
		tx.Where("display", display).Order("sort asc")
	}
	result := tx.Find(&menus)
	return menus, errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
}

func (r *menuRepo) Add(menu *model.Menu) error {
	result := database.RMSDB.Create(&menu)
	return errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
}

func (r *menuRepo) Update(id int64, menu map[string]interface{}) error {
	result := database.RMSDB.Table("menus").Where("id", id).Updates(&menu)
	return errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
}

func (r *menuRepo) Delete(id int64) error {
	result := database.RMSDB.Delete(&model.Menu{}, id)
	return errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
}
