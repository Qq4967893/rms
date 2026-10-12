package repo

import (
	"backend-rms/core/database"
	"backend-rms/model"
	"backend-rms/util"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type groupPairRepo struct{}

var GroupPairRepo = &groupPairRepo{}

func (r *groupPairRepo) IsAnyGroupExisted(groupA, groupB string, excludeId int64) (bool, error) {
	var count int64
	tx := database.RMSDB.Model(model.GroupPair{}).Where("group_a = ? or group_b = ?", groupA, groupB)
	if excludeId != 0 {
		tx.Where("id != ?", excludeId)
	}
	result := tx.Count(&count)

	return count > 0, errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
}

func (r *groupPairRepo) GetInfoById(id int64) (groupPair *model.GroupPair, err error) {
	err = database.RMSDB.Where("id", id).First(&groupPair).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return groupPair, errors.Wrapf(err, "%s failed", util.GetCurrentFuncName())
}

func (r *groupPairRepo) GetAll() (groupPairs []*model.GroupPair, err error) {
	result := database.RMSDB.Find(&groupPairs)
	return groupPairs, errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
}

func (r *groupPairRepo) Add(groupPair *model.GroupPair) error {
	result := database.RMSDB.Create(groupPair)
	return errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
}

func (r *groupPairRepo) Update(groupPair *model.GroupPair) error {
	result := database.RMSDB.Model(groupPair).Select("group_a", "group_b").Updates(groupPair)
	return errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
}

func (r *groupPairRepo) Delete(id int64) (int64, error) {
	result := database.RMSDB.Delete(&model.GroupPair{}, id)
	return result.RowsAffected, errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
}
