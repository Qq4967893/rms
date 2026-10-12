package controller

import (
	"backend-rms/http/request"
	"backend-rms/http/response"
	"backend-rms/model"
	"backend-rms/service"

	"github.com/gin-gonic/gin"
)

type PermissionController struct {
}

func (rc *PermissionController) GetInfo(c *gin.Context) {
	req := request.GetPermissionInfoRequest{}
	if ok := request.Validate(c, &req); !ok {
		return
	}

	result, err := service.RoleMenuService.GetAllByRoleId(req.Id)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.SucceedWithData(c, "", result)
}

func (rc *PermissionController) Update(c *gin.Context) {
	req := request.UpdatePermissionRequest{}
	if ok := request.Validate(c, &req); !ok {
		return
	}

	isRoleIdExisted, isRoleIdExistedErr := service.RoleService.IsIdExisted(req.Id)
	if isRoleIdExistedErr != nil {
		response.Error(c, isRoleIdExistedErr)
		return
	}
	if !isRoleIdExisted {
		response.FailWithMessage(c, "The role ID not exists")
		return
	}

	var role model.Role
	role.Name = req.Name
	err := service.RoleService.Update(req.Id, role)
	if err != nil {
		response.Error(c, err)
		return
	}

	err = service.RoleMenuService.UpdateRoleMenus(req.Id, req.MenuIds)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Succeed(c)
}
