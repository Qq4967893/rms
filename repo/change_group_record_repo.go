package repo

import (
	"backend-rms/core/database"
	"backend-rms/model"
	"backend-rms/util"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type changeGroupRecordRepo struct{}

var ChangeGroupRecordRepo = &changeGroupRecordRepo{}

func (cr *changeGroupRecordRepo) GetInfo(login int64, status int8) (changeGroupRecord *model.ChangeGroupRecord, err error) {
	err = database.RMSDB.Where("login", login).Where("status", status).First(&changeGroupRecord).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return changeGroupRecord, errors.Wrapf(err, "%s failed", util.GetCurrentFuncName())
}

func (cr *changeGroupRecordRepo) GetOperators() (operators []*model.Operator, err error) {
	result := database.RMSDB.Table("change_group_records as c").Joins("left join managers as m on c.manager_id = m.id").
		Group("c.manager_id").Select("m.username, c.manager_id").Find(&operators)
	return operators, errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
}

func (cr *changeGroupRecordRepo) GetAll(logins []int64, managerId int64, paging *model.Paging) (changeGroupRecords []*model.ChangeGroupRecord, err error) {
	tx := database.RMSDB.Table("change_group_records")
	if len(logins) > 0 {
		tx.Where("login in ?", logins)
	}

	if managerId != 0 {
		tx.Where("manager_id", managerId)
	}

	if paging != nil {
		tx.Order(paging.Sort + " " + paging.Order).Limit(paging.Limit).Offset(paging.Offset)
	}

	result := tx.Find(&changeGroupRecords)
	return changeGroupRecords, errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
}

func (cr *changeGroupRecordRepo) GetCount(logins []int64, managerId int64) (totalCount int64, err error) {
	tx := database.RMSDB.Table("change_group_records")
	if len(logins) > 0 {
		tx.Where("login in ?", logins)
	}

	if managerId != 0 {
		tx.Where("manager_id", managerId)
	}

	result := tx.Count(&totalCount)
	return totalCount, errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
}

func (cr *changeGroupRecordRepo) Add(changeGroupRecord model.ChangeGroupRecord) (int64, error) {
	result := database.RMSDB.Create(&changeGroupRecord)
	return changeGroupRecord.ID, errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
}

func (cr *changeGroupRecordRepo) Update(id int64, changeGroupRecord model.ChangeGroupRecord) error {
	result := database.RMSDB.Model(changeGroupRecord).Where("id", id).Updates(changeGroupRecord)
	return errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
}
