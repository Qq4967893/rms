package controller

import (
	"backend-rms/http/request"
	"backend-rms/http/response"
	"backend-rms/model"
	"backend-rms/service"

	"github.com/gin-gonic/gin"
)

type RoleController struct {
}

func (rc *RoleController) GetDisplayMenus(c *gin.Context) {
	req := request.GetDisplayMenusRequest{}
	if ok := request.Validate(c, &req); !ok {
		return
	}

	result, err := service.RoleMenuService.GetDisplayMenus(req.Id)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.SucceedWithData(c, "", result)
}

func (rc *RoleController) GetList(c *gin.Context) {
	result, err := service.RoleService.GetList()
	if err != nil {
		response.Error(c, err)
		return
	}
	response.SucceedWithData(c, "", result)
}

func (rc *RoleController) Add(c *gin.Context) {
	req := request.AddRoleRequest{}
	if ok := request.Validate(c, &req); !ok {
		return
	}

	role := &model.Role{
		Name: req.Name,
	}
	err := service.RoleService.Add(role, req.MenuIds)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Succeed(c)
}

func (rc *RoleController) Delete(c *gin.Context) {
	req := request.DeleteRoleRequest{}
	if ok := request.Validate(c, &req); !ok {
		return
	}

	//判断该角色是否有manager进行关联
	manager, getInfoByRoleIdErr := service.ManagerService.GetInfoByRoleId(req.Id)
	if getInfoByRoleIdErr != nil {
		response.Error(c, getInfoByRoleIdErr)
		return
	}
	if manager != nil {
		response.FailWithMessage(c, "This role is still used by users")
		return
	}

	err := service.RoleService.Delete(req.Id)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Succeed(c)
}
