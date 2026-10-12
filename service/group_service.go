package service

import (
	"backend-rms/model"
	"backend-rms/repo"
	"backend-rms/util"
	"strings"

	"github.com/pkg/errors"
)

type groupService struct {
}

var GroupService = &groupService{}

func (s *groupService) IsNameExisted(name string, excludeId int64) (bool, error) {
	return repo.GroupRepo.IsNameExisted(name, excludeId)
}

func (s *groupService) GetInfoById(id int64) (group *model.Group, err error) {
	group, err = repo.GroupRepo.GetInfoById(id)
	return group, errors.WithMessage(err, util.ErrWrapMessage())
}

func (s *groupService) GetBookByGroup(group string) (string, error) {
	groupInfo, err := repo.GroupRepo.GetInfoByName(group)
	if err != nil {
		return "", errors.WithMessage(err, util.ErrWrapMessage())
	}

	if groupInfo == nil {
		return "", errors.New("the group " + group + "is not found")
	}
	return groupInfo.Book, nil
}

func (s *groupService) GetGroupsByBook(book string) (groups []string, err error) {
	book = strings.ToUpper(book)
	groupInfos, err := repo.GroupRepo.GetAllByBook(book)
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}
	for _, v := range groupInfos {
		groups = append(groups, v.Name)
	}
	return groups, nil
}

func (s *groupService) GetAllList() (groups []*model.Group, err error) {
	groups, err = repo.GroupRepo.GetAll()
	return groups, errors.WithMessage(err, util.ErrWrapMessage())
}

func (s *groupService) GetAllMap() (map[string]string, error) {
	groups, err := repo.GroupRepo.GetAll()
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}

	var groupMap = make(map[string]string)
	for _, g := range groups {
		groupMap[g.Name] = g.Book
	}
	return groupMap, nil
}

func (s *groupService) Add(group *model.Group) error {
	return repo.GroupRepo.Add(group)
}

func (s *groupService) Update(group *model.Group) error {
	return repo.GroupRepo.Update(group)
}
