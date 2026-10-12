package repo

import (
	"backend-rms/core/database"
	"backend-rms/model"
	"backend-rms/util"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type managerRepo struct{}

var ManagerRepo = &managerRepo{}

func (r *managerRepo) GetInfoById(id int64) (manager *model.Manager, err error) {
	err = database.RMSDB.Where("id", id).First(&manager).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return manager, errors.Wrapf(err, "%s failed", util.GetCurrentFuncName())
}

func (r *managerRepo) GetInfoByUsername(username string) (manager *model.Manager, err error) {
	err = database.RMSDB.Where("username", username).First(&manager).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return manager, errors.Wrapf(err, "%s failed", util.GetCurrentFuncName())
}

func (r *managerRepo) GetInfoByRoleId(roleId int64) (manager *model.Manager, err error) {
	err = database.RMSDB.Where("role_id", roleId).First(&manager).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return manager, errors.Wrapf(err, "%s failed", util.GetCurrentFuncName())
}

func (r *managerRepo) GetList(username string, paging model.Paging) (managers []model.ManagerGetList, err error) {
	tx := database.RMSDB.Table("managers as m").
		Joins("left join roles as r on m.role_id = r.id")
	if len(username) != 0 {
		tx.Where("m.username LIKE ?", "%"+username+"%")
	}
	result := tx.Select("m.id, m.username, m.email, m.role_id, r.name as role_name").Limit(paging.Limit).Offset(paging.Offset).Find(&managers)
	return managers, errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
}

func (r *managerRepo) GetListCount(username string) (totalCount int64, err error) {
	tx := database.RMSDB.Table("managers as m").
		Joins("left join roles as r on m.role_id = r.id")
	if len(username) != 0 {
		tx.Where("m.username LIKE ?", "%"+username+"%")
	}
	result := tx.Count(&totalCount)

	return totalCount, errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
}

func (r *managerRepo) Add(params model.Manager) (int64, error) {
	result := database.RMSDB.Model(model.Manager{}).Create(&params)
	return params.ID, errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
}

func (r *managerRepo) Update(id int64, params model.Manager) error {
	result := database.RMSDB.Where("id", id).Updates(&params)
	return errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
}

func (r *managerRepo) Delete(id int64) error {
	result := database.RMSDB.Delete(&model.Manager{}, id)
	return errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
}
