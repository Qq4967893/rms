package controller

import (
	"backend-rms/http/request"
	"backend-rms/http/response"
	"backend-rms/service"

	"github.com/gin-gonic/gin"
)

type PassportController struct{}

func (pc *PassportController) GetAccounts(c *gin.Context) {
	req := request.GetAccountsRequest{}
	if ok := request.Validate(c, &req); !ok {
		return
	}

	accounts, err := service.PassportService.GetAccounts(req.Passport, req.Login)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SucceedWithData(c, "", accounts)
}
