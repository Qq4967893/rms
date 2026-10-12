package service

import (
	myCasbin "backend-rms/core/casbin"
	"backend-rms/model"
	"backend-rms/repo"
	"backend-rms/util"

	"github.com/pkg/errors"
)

type roleService struct {
}

var RoleService = &roleService{}

func (s *roleService) IsIdExisted(id int64) (bool, error) {
	role, err := repo.RoleRepo.GetInfoById(id)
	if role == nil {
		return false, errors.WithMessage(err, util.ErrWrapMessage())
	}
	return true, err
}

func (s *roleService) GetList() (roles []*model.Role, err error) {
	roles, err = repo.RoleRepo.GetList()
	return roles, errors.WithMessage(err, util.ErrWrapMessage())
}

func (s *roleService) Add(role *model.Role, menuIds []int) (err error) {
	id, err := repo.RoleRepo.Add(role)
	if err != nil {
		return errors.WithMessage(err, util.ErrWrapMessage())
	}
	if len(menuIds) == 0 {
		return
	}
	err = RoleMenuService.Add(id, menuIds)
	if err != nil {
		return errors.WithMessage(err, util.ErrWrapMessage())
	}

	policies, err := RoleMenuService.BuildRoleMenuPolicies([]int64{id}, menuIds)
	if err != nil {
		return errors.WithMessage(err, util.ErrWrapMessage())
	}
	err = myCasbin.AddPolicies("p", policies)
	if err != nil {
		return errors.WithMessage(err, util.ErrWrapMessage())
	}
	return
}

func (s *roleService) Update(id int64, role model.Role) (err error) {
	return repo.RoleRepo.Update(id, role)
}

func (s *roleService) Delete(id int64) (err error) {
	err = repo.RoleRepo.Delete(id)
	if err != nil {
		return errors.WithMessage(err, util.ErrWrapMessage())
	}

	err = RoleMenuService.Delete(id)
	err = myCasbin.RemoveRolePolicy(id)
	if err != nil {
		return errors.WithMessage(err, util.ErrWrapMessage())
	}
	return
}
