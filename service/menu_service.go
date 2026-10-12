package service

import (
	myCasbin "backend-rms/core/casbin"
	"backend-rms/model"
	"backend-rms/repo"
	"backend-rms/util"

	"github.com/pkg/errors"
)

type menuService struct {
}

var MenuService = &menuService{}

type MenuTree struct {
	Id            int64       `json:"id"`
	Name          string      `json:"name"`
	Pid           int64       `json:"pid"`
	Path          string      `json:"path"`
	Display       int8        `json:"display"`
	Sort          int8        `json:"sort"`
	NextLevelMenu []*MenuTree `json:"nextLevelMenu"`
}

func (s *menuService) IsIdExisted(id int64) (bool, error) {
	role, err := repo.MenuRepo.GetInfoById(id)
	if role == nil {
		return false, errors.WithMessage(err, util.ErrWrapMessage())
	}
	return true, errors.WithMessage(err, util.ErrWrapMessage())
}

func (s *menuService) GetMenuList() (menuList []*MenuTree, err error) {
	menuTree := MenuTree{}
	err = s.getNextLevelMenu(0, &menuTree)
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}
	menuList = menuTree.NextLevelMenu
	return
}

func (s *menuService) Add(menu *model.Menu) (err error) {
	err = repo.MenuRepo.Add(menu)
	return errors.WithMessage(err, util.ErrWrapMessage())
}

func (s *menuService) Update(id int64, menu map[string]interface{}) (err error) {
	oldMenu, err := repo.MenuRepo.GetInfoById(id)
	if err != nil {
		return errors.WithMessage(err, util.ErrWrapMessage())
	}

	err = repo.MenuRepo.Update(id, menu)
	if err != nil {
		return errors.WithMessage(err, util.ErrWrapMessage())
	}

	//判断该菜单是否和角色有关联
	roleMenus, err := repo.RoleMenuRepo.GetAllByMenuId(id)
	if err != nil {
		return errors.WithMessage(err, util.ErrWrapMessage())
	}

	if len(roleMenus) == 0 {
		return nil
	}

	//casbin的删除操作
	err = myCasbin.RemoveMenuPolicy(oldMenu.Path)
	if err != nil {
		return errors.WithMessage(err, util.ErrWrapMessage())
	}

	//判断该菜单的path字段内是否有路由地址，有路由地址进行casbin添加操作
	if menu["path"] == "[]" {
		return nil
	}

	var roleIds []int64
	for _, v := range roleMenus {
		roleIds = append(roleIds, v.Role_id)
	}
	policies, err := RoleMenuService.BuildRoleMenuPolicies(roleIds, []int{int(id)})
	if err != nil {
		return errors.WithMessage(err, util.ErrWrapMessage())
	}
	err = myCasbin.AddPolicies("p", policies)
	if err != nil {
		return errors.WithMessage(err, util.ErrWrapMessage())
	}
	return nil
}

func (s *menuService) Delete(id int64) error {
	oldMenu, err := repo.MenuRepo.GetInfoById(id)
	if err != nil {
		return errors.WithMessage(err, util.ErrWrapMessage())
	}

	err = repo.MenuRepo.Delete(id)
	if err != nil {
		return errors.WithMessage(err, util.ErrWrapMessage())
	}

	//查看是否有菜单角色关联关系，有再进行删除casbin操作
	roleMenus, err := repo.RoleMenuRepo.GetAllByMenuId(id)
	if err != nil {
		return errors.WithMessage(err, util.ErrWrapMessage())
	}

	if len(roleMenus) == 0 {
		return nil
	}

	err = repo.RoleMenuRepo.Delete(0, 20)
	if err != nil {
		return errors.WithMessage(err, util.ErrWrapMessage())
	}

	err = myCasbin.RemoveMenuPolicy(oldMenu.Path)
	if err != nil {
		return errors.WithMessage(err, util.ErrWrapMessage())
	}
	return nil
}

func (s *menuService) getNextLevelMenu(pid int64, menuTree *MenuTree) (err error) {
	var menus []model.Menu
	menus, err = repo.MenuRepo.GetInfoByPid(pid)
	if err != nil {
		return errors.WithMessage(err, util.ErrWrapMessage())
	}

	for i := 0; i < len(menus); i++ {
		nextLevelMenu := MenuTree{menus[i].ID, menus[i].Name, menus[i].Pid, menus[i].Path, menus[i].Display, menus[i].Sort, []*MenuTree{}}
		menuTree.NextLevelMenu = append(menuTree.NextLevelMenu, &nextLevelMenu)
		err = s.getNextLevelMenu(menus[i].ID, &nextLevelMenu)
		if err != nil {
			return errors.WithMessage(err, util.ErrWrapMessage())
		}
	}
	return
}
