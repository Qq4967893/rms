package service

import (
	myCasbin "backend-rms/core/casbin"
	"backend-rms/model"
	"backend-rms/repo"
	"backend-rms/util"
	"encoding/json"
	"strconv"

	"github.com/pkg/errors"
)

type roleMenuService struct {
}

var RoleMenuService = &roleMenuService{}

func (s *roleMenuService) GetDisplayMenus(roleId int64) ([]any, error) {
	roleMenus, err := repo.RoleMenuRepo.GetAllByRoleId(roleId)
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}
	var menuIds []int
	for _, v := range roleMenus {
		menuIds = append(menuIds, v.Menu_id)
	}

	roleMenuInfo, _ := repo.MenuRepo.GetAllByIds(menuIds, 1)
	//处理Display Menus
	var result []any
	for _, v := range roleMenuInfo {
		if v.Pid == 0 {
			displayMenus := map[string]any{
				"name": v.Name,
			}
			var nextLevelMenu []any
			for _, vn := range roleMenuInfo {
				if vn.Pid == v.ID {
					nextLevelMenuInfo := map[string]any{
						"name": vn.Name,
					}
					nextLevelMenu = append(nextLevelMenu, nextLevelMenuInfo)
				}
			}
			displayMenus["nextLevelMenu"] = nextLevelMenu
			result = append(result, displayMenus)
		}
	}
	return result, nil
}

func (s *roleMenuService) GetAllByRoleId(roleId int64) (roleMenus []*model.RoleMenu, err error) {
	roleMenus, err = repo.RoleMenuRepo.GetAllByRoleId(roleId)
	return roleMenus, errors.WithMessage(err, util.ErrWrapMessage())
}

func (s *roleMenuService) Add(roleId int64, menuIds []int) error {
	var roleMenus []model.RoleMenu
	for _, menuId := range menuIds {
		var roleMenu = model.RoleMenu{
			Role_id: roleId,
			Menu_id: menuId,
		}
		roleMenus = append(roleMenus, roleMenu)
	}

	err := repo.RoleMenuRepo.AddBatches(roleMenus)
	if err != nil {
		return errors.WithMessage(err, util.ErrWrapMessage())
	}
	return nil
}

func (s *roleMenuService) UpdateRoleMenus(roleId int64, menuIds []int) error {
	err := s.Delete(roleId)
	if err != nil {
		return errors.WithMessage(err, util.ErrWrapMessage())
	}

	err = myCasbin.RemoveRolePolicy(roleId)
	if err != nil {
		return errors.WithMessage(err, util.ErrWrapMessage())
	}

	err = s.Add(roleId, menuIds)
	if err != nil {
		return errors.WithMessage(err, util.ErrWrapMessage())
	}

	policies, err := s.BuildRoleMenuPolicies([]int64{roleId}, menuIds)
	if err != nil {
		return errors.WithMessage(err, util.ErrWrapMessage())
	}
	err = myCasbin.AddPolicies("p", policies)
	if err != nil {
		return errors.WithMessage(err, util.ErrWrapMessage())
	}
	return nil
}

func (s *roleMenuService) Delete(roleId int64) (err error) {
	err = repo.RoleMenuRepo.Delete(roleId, 0)
	return errors.WithMessage(err, util.ErrWrapMessage())
}

func (s *roleMenuService) BuildRoleMenuPolicies(roleIds []int64, menuIds []int) ([][]string, error) {
	//获取model.Menu.Path
	menus, err := repo.MenuRepo.GetAllByIds(menuIds, 0)
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}
	var policies [][]string
	for _, roleId := range roleIds {
		for _, vm := range menus {
			//如果菜单内没有接口信息
			if vm.Path == "[]" {
				continue
			}

			var path []map[string]string
			err := json.Unmarshal([]byte(vm.Path), &path)
			if err != nil {
				return nil, errors.Wrap(err, "json: unmarshal failed")
			}
			for _, vp := range path {
				policy := []string{
					"role_" + strconv.FormatInt(roleId, 10),
					vp["api"],
					vp["method"],
				}
				policies = append(policies, policy)
			}
		}
	}
	return policies, nil
}
