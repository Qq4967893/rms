package controller

import (
	"backend-rms/core/jwt"
	"backend-rms/http/request"
	"backend-rms/http/response"
	"backend-rms/model"
	"backend-rms/service"
	"backend-rms/util"
	"backend-rms/util/crypto"

	"github.com/gin-gonic/gin"
)

type ManagerController struct{}

func (mc *ManagerController) GetLoggedDetails(c *gin.Context) {
	loginManagerID := c.GetInt64("login_manager_id")
	result, err := service.ManagerService.GetInfoById(loginManagerID)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.SucceedWithData(c, "", result)
}

func (mc *ManagerController) GetList(c *gin.Context) {
	req := request.GetManagerListRequest{}
	if ok := request.Validate(c, &req); !ok {
		return
	}

	username := c.DefaultQuery("username", "")
	totalCount, err := service.ManagerService.GetListCount(username)
	if err != nil {
		response.Error(c, err)
		return
	}

	paging := model.NewPaging(req.Limit, req.Offset)
	managerList, err := service.ManagerService.GetList(username, paging)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.SucceedWithData(c, "success", gin.H{
		"manager_list": managerList,
		"total_count":  totalCount,
	})
}

func (mc *ManagerController) Add(c *gin.Context) {
	req := request.AddManagerRequest{}
	if ok := request.Validate(c, &req); !ok {
		return
	}

	checkMessage := service.ManagerService.CheckPassword(req.Password)
	if checkMessage != "" {
		response.FailWithMessage(c, checkMessage)
		return
	}

	//验证角色id是否存在
	isRoleIdExisted, isRoleIdExistedErr := service.RoleService.IsIdExisted(req.RoleId)
	if isRoleIdExistedErr != nil {
		response.Error(c, isRoleIdExistedErr)
		return
	}
	if !isRoleIdExisted {
		response.FailWithMessage(c, "The role ID not exists")
		return
	}

	var managerParams model.Manager
	managerParams.Username = req.UserName
	managerParams.Password, _ = crypto.PasswordHash(req.Password)
	managerParams.RoleId = req.RoleId
	managerParams.Email = req.Email
	err := service.ManagerService.Add(managerParams)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Succeed(c)
}

func (mc *ManagerController) Update(c *gin.Context) {
	req := request.UpdateManagerRequest{}
	if ok := request.Validate(c, &req); !ok {
		return
	}

	isManagerIdExisted, isManagerIdExistedErr := service.ManagerService.IsIdExisted(req.Id)
	if isManagerIdExistedErr != nil {
		response.Error(c, isManagerIdExistedErr)
		return
	}
	if !isManagerIdExisted {
		response.FailWithMessage(c, "The Manager ID not exists")
		return
	}

	isRoleIdExisted, isRoleIdExistedErr := service.RoleService.IsIdExisted(req.RoleId)
	if isRoleIdExistedErr != nil {
		response.Error(c, isRoleIdExistedErr)
		return
	}
	if !isRoleIdExisted {
		response.FailWithMessage(c, "The role ID not exists")
		return
	}

	var managerParams model.Manager
	managerParams.Email = req.Email
	managerParams.RoleId = req.RoleId

	err := service.ManagerService.Update(req.Id, managerParams)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Succeed(c)
}

func (mc *ManagerController) Delete(c *gin.Context) {
	req := request.DeleteManagerRequest{}
	if ok := request.Validate(c, &req); !ok {
		return
	}

	err := service.ManagerService.Delete(req.Id)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Succeed(c)
}

func (mc *ManagerController) ResetPassword(c *gin.Context) {
	req := request.ResetManagerPasswordRequest{}
	if ok := request.Validate(c, &req); !ok {
		return
	}

	isManagerIdExisted, isManagerIdExistedErr := service.ManagerService.IsIdExisted(req.Id)
	if isManagerIdExistedErr != nil {
		response.Error(c, isManagerIdExistedErr)
		return
	}
	if !isManagerIdExisted {
		response.FailWithMessage(c, "The Manager ID not exists")
		return
	}

	resetPassword := util.GenerateRandomString(8)
	password, passwordHashErr := crypto.PasswordHash(resetPassword)
	if passwordHashErr != nil {
		response.Error(c, passwordHashErr)
		return
	}
	var managerParams model.Manager
	managerParams.Password = password
	err := service.ManagerService.Update(req.Id, managerParams)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.SucceedWithData(c, "", map[string]string{"new_password": resetPassword})
}

func (mc *ManagerController) ModifyPassword(c *gin.Context) {
	req := request.ModifyManagerPasswordRequest{}
	if ok := request.Validate(c, &req); !ok {
		return
	}

	if req.NewPassword != req.RepeatNewPassword {
		response.FailWithMessage(c, "the two passwords are inconsistent")
		return
	}

	checkMessage := service.ManagerService.CheckPassword(req.NewPassword)
	if checkMessage != "" {
		response.FailWithMessage(c, checkMessage)
		return
	}

	loginManagerID := c.GetInt64("login_manager_id")
	//验证旧密码是否正确
	managerInfo, getInfoByIdErr := service.ManagerService.GetInfoById(loginManagerID)
	if getInfoByIdErr != nil {
		response.Error(c, getInfoByIdErr)
		return
	}
	if !crypto.PasswordVerify(req.OldPassword, managerInfo.Password) {
		response.FailWithMessage(c, "manager password input error")
		return
	}

	password, passwordHashErr := crypto.PasswordHash(req.NewPassword)
	if passwordHashErr != nil {
		response.Error(c, passwordHashErr)
		return
	}
	var managerParams model.Manager
	managerParams.Password = password
	err := service.ManagerService.Update(loginManagerID, managerParams)
	if err != nil {
		response.Error(c, err)
		return
	}

	accessToken, err := jwt.Instance.GetTokenFromHeader(c.Request)
	if err != nil {
		response.Error(c, err)
		return
	}

	managerID := c.GetInt64("login_manager_id")
	if err = service.AuthService.DeleteToken(managerID, accessToken); err != nil {
		response.Error(c, err)
		return
	}
	response.Succeed(c)
}
