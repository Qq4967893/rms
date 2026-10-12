package controller

import (
	"backend-rms/http/request"
	"backend-rms/http/response"
	"backend-rms/model"
	"backend-rms/service"

	"github.com/gin-gonic/gin"
)

type HedgeAccountController struct{}

func (hc *HedgeAccountController) GetList(c *gin.Context) {
	result, err := service.HedgeAccountService.GetList()
	if err != nil {
		response.Error(c, err)
		return
	}
	response.SucceedWithData(c, "", result)
}

func (hc *HedgeAccountController) Update(c *gin.Context) {
	req := request.UpdateHedgeAccountRequest{}
	if ok := request.Validate(c, &req); !ok {
		return
	}

	var hedgeAccountParams model.HedgeAccount
	hedgeAccountParams.RoleIds = req.RoleIds
	hedgeAccountParams.Desc = req.Desc

	err := service.HedgeAccountService.Update(req.Id, hedgeAccountParams)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Succeed(c)
}
