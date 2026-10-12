package service

import (
	"backend-rms/model"
	"backend-rms/repo"
	"backend-rms/util"

	"github.com/pkg/errors"
)

type changeGroupRecordService struct{}

var ChangeGroupRecordService = &changeGroupRecordService{}

// 账号是否存在处理中的转组记录
func (s *changeGroupRecordService) IsProcessingExisted(login int64) (bool, error) {
	changeGroupRecord, err := repo.ChangeGroupRecordRepo.GetInfo(login, 0)
	if err != nil {
		return false, errors.WithMessage(err, util.ErrWrapMessage())
	}
	if changeGroupRecord == nil {
		return false, nil
	}
	return true, nil
}

func (s *changeGroupRecordService) GetOperators() ([]*model.Operator, error) {
	return repo.ChangeGroupRecordRepo.GetOperators()
}

func (s *changeGroupRecordService) GetList(logins []int64, managerId int64, paging *model.Paging) ([]*model.ChangeGroupRecord, error) {
	return repo.ChangeGroupRecordRepo.GetAll(logins, managerId, paging)
}

func (s *changeGroupRecordService) GetCount(logins []int64, managerId int64) (int64, error) {
	return repo.ChangeGroupRecordRepo.GetCount(logins, managerId)
}

func (s *changeGroupRecordService) Add(changeGroupRecord model.ChangeGroupRecord) (int64, error) {
	return repo.ChangeGroupRecordRepo.Add(changeGroupRecord)
}

func (s *changeGroupRecordService) CompleteChangeGroup(id int64) error {
	changeGroupRecord := model.ChangeGroupRecord{
		Status: 1,
	}
	return repo.ChangeGroupRecordRepo.Update(id, changeGroupRecord)
}
