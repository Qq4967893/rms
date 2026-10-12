package service

import (
	myCasbin "backend-rms/core/casbin"
	"backend-rms/model"
	"backend-rms/repo"
	"backend-rms/util"
	"regexp"
	"strconv"

	"github.com/pkg/errors"
)

type managerService struct {
}

var ManagerService = &managerService{}

func (s *managerService) IsIdExisted(id int64) (bool, error) {
	role, err := repo.ManagerRepo.GetInfoById(id)
	if role == nil {
		return false, errors.WithMessage(err, util.ErrWrapMessage())
	}
	return true, errors.WithMessage(err, util.ErrWrapMessage())
}

func (s *managerService) GetInfoById(id int64) (manager *model.Manager, err error) {
	manager, err = repo.ManagerRepo.GetInfoById(id)
	return manager, errors.WithMessage(err, util.ErrWrapMessage())
}

func (s *managerService) GetInfoByUsername(username string) (manager *model.Manager, err error) {
	manager, err = repo.ManagerRepo.GetInfoByUsername(username)
	return manager, errors.WithMessage(err, util.ErrWrapMessage())
}

func (s *managerService) GetInfoByRoleId(roleId int64) (manager *model.Manager, err error) {
	manager, err = repo.ManagerRepo.GetInfoByRoleId(roleId)
	return manager, errors.WithMessage(err, util.ErrWrapMessage())
}

func (s *managerService) GetList(username string, paging model.Paging) (managers []model.ManagerGetList, err error) {
	managers, err = repo.ManagerRepo.GetList(username, paging)
	return managers, errors.WithMessage(err, util.ErrWrapMessage())
}

func (s *managerService) GetListCount(username string) (totalCount int64, err error) {
	totalCount, err = repo.ManagerRepo.GetListCount(username)
	return totalCount, errors.WithMessage(err, util.ErrWrapMessage())
}

func (s *managerService) Add(params model.Manager) error {
	id, err := repo.ManagerRepo.Add(params)
	if err != nil {
		return errors.WithMessage(err, util.ErrWrapMessage())
	}

	policy := []string{myCasbin.CASBIN_G_PREFIX + strconv.Itoa(int(id)), myCasbin.CASBIN_P_PREFIX + strconv.Itoa(int(params.RoleId))}
	err = myCasbin.AddPolicy("g", policy)
	if err != nil {
		return errors.WithMessage(err, util.ErrWrapMessage())
	}
	return nil
}

func (s *managerService) Update(id int64, params model.Manager) error {
	//获取原manager的角色id
	manager, err := repo.ManagerRepo.GetInfoById(id)
	if err != nil {
		return errors.WithMessage(err, util.ErrWrapMessage())
	}

	err = repo.ManagerRepo.Update(id, params)
	if err != nil {
		return errors.WithMessage(err, util.ErrWrapMessage())
	}

	//如果没有修改用户角色，不需要做处理
	if (manager.RoleId != params.RoleId) && (params.RoleId != 0) {
		removePolicy := []string{myCasbin.CASBIN_G_PREFIX + strconv.Itoa(int(id)), myCasbin.CASBIN_P_PREFIX + strconv.Itoa(int(manager.RoleId))}
		err := myCasbin.RemovePolicy("g", removePolicy)
		if err != nil {
			return errors.WithMessage(err, util.ErrWrapMessage())
		}

		addPolicy := []string{myCasbin.CASBIN_G_PREFIX + strconv.Itoa(int(id)), myCasbin.CASBIN_P_PREFIX + strconv.Itoa(int(params.RoleId))}
		err = myCasbin.AddPolicy("g", addPolicy)
		if err != nil {
			return errors.WithMessage(err, util.ErrWrapMessage())
		}
	}
	return nil
}

func (s *managerService) Delete(id int64) error {
	//获取原manager的角色id
	manager, err := repo.ManagerRepo.GetInfoById(id)
	if err != nil {
		return errors.WithMessage(err, util.ErrWrapMessage())
	}

	err = repo.ManagerRepo.Delete(id)
	if err != nil {
		return errors.WithMessage(err, util.ErrWrapMessage())
	}

	policy := []string{myCasbin.CASBIN_G_PREFIX + strconv.Itoa(int(id)), myCasbin.CASBIN_P_PREFIX + strconv.Itoa(int(manager.RoleId))}
	err = myCasbin.RemovePolicy("g", policy)
	if err != nil {
		return errors.WithMessage(err, util.ErrWrapMessage())
	}
	return nil
}

func (s *managerService) CheckPassword(password string) string {
	if len(password) < 8 || len(password) > 24 {
		return "The password should be more than 8 characters but less than 24 characters in length"
	}
	hasLower := regexp.MustCompile(`[a-z]`)
	hasUpper := regexp.MustCompile(`[A-Z]`)
	hasDigit := regexp.MustCompile(`[0-9]`)
	if hasLower.MatchString(password) && hasUpper.MatchString(password) && hasDigit.MatchString(password) {
		return ""
	} else {
		return "The password must include uppercase and lowercase letters, numbers"
	}
}
