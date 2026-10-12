package service

import (
	"backend-rms/model"
	"backend-rms/repo"
	"backend-rms/util"

	"github.com/pkg/errors"
)

type groupPairService struct {
}

var GroupPairService = &groupPairService{}

func (s *groupPairService) IsAnyGroupExisted(groupA, groupB string, excludeId int64) (bool, error) {
	isExisted, err := repo.GroupPairRepo.IsAnyGroupExisted(groupA, groupB, excludeId)
	return isExisted, errors.WithMessage(err, util.ErrWrapMessage())
}

func (s *groupPairService) GetInfoById(id int64) (groupPair *model.GroupPair, err error) {
	groupPair, err = repo.GroupPairRepo.GetInfoById(id)
	return groupPair, errors.WithMessage(err, util.ErrWrapMessage())
}

func (s *groupPairService) GetAll() (groupPairs []*model.GroupPair, err error) {
	groupPairs, err = repo.GroupPairRepo.GetAll()
	return groupPairs, errors.WithMessage(err, util.ErrWrapMessage())
}

// 获取允许转到的组
func (s *groupPairService) GetToGroups() (map[string]string, error) {
	groupPairs, err := repo.GroupPairRepo.GetAll()
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}

	newGroupPairs := make(map[string]string, 0)
	for _, groupPair := range groupPairs {
		newGroupPairs[groupPair.GroupA] = groupPair.GroupB
		newGroupPairs[groupPair.GroupB] = groupPair.GroupA
	}
	return newGroupPairs, nil
}

func (s *groupPairService) Add(groupPair *model.GroupPair) error {
	return repo.GroupPairRepo.Add(groupPair)
}

func (s *groupPairService) Update(groupPair *model.GroupPair) error {
	err := repo.GroupPairRepo.Update(groupPair)
	return errors.WithMessage(err, util.ErrWrapMessage())
}

func (s *groupPairService) Delete(id int64) (int64, error) {
	rows, err := repo.GroupPairRepo.Delete(id)
	return rows, errors.WithMessage(err, util.ErrWrapMessage())
}
