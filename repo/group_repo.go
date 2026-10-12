package repo

import (
	"backend-rms/core/database"
	"backend-rms/model"
	"backend-rms/util"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type groupRepo struct{}

var GroupRepo = &groupRepo{}

func (r *groupRepo) IsNameExisted(name string, excludeId int64) (bool, error) {
	var count int64
	tx := database.RMSDB.Model(model.Group{}).Where("name = ?", name)
	if excludeId != 0 {
		tx.Where("id != ?", excludeId)
	}
	result := tx.Count(&count)

	return count > 0, errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
}

func (r *groupRepo) GetInfoById(id int64) (group *model.Group, err error) {
	err = database.RMSDB.Where("Id", id).First(&group).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return group, errors.Wrapf(err, "%s failed", util.GetCurrentFuncName())
}

func (r *groupRepo) GetInfoByName(name string) (group *model.Group, err error) {
	err = database.RMSDB.Where("name", name).First(&group).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return group, errors.Wrapf(err, "%s failed", util.GetCurrentFuncName())
}

func (r *groupRepo) GetAll() (groups []*model.Group, err error) {
	result := database.RMSDB.Find(&groups)
	return groups, errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
}

func (r *groupRepo) GetAllByBook(book string) (groups []*model.Group, err error) {
	tx := database.RMSDB.Model(model.Group{})
	if book != "ALL" {
		tx.Where("book", book)
	}
	result := tx.Find(&groups)
	return groups, errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
}

func (r *groupRepo) Add(group *model.Group) error {
	result := database.RMSDB.Create(group)

	return errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
}

func (r *groupRepo) Update(group *model.Group) error {
	result := database.RMSDB.Model(group).Select("name", "book").Updates(group)
	return errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
}
