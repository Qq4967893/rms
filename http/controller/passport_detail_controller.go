package controller

import (
	"backend-rms/http/request"
	"backend-rms/http/response"
	"backend-rms/model"
	"backend-rms/service"

	"github.com/gin-gonic/gin"
)

type PassportDetailController struct{}

func (pdc *PassportDetailController) GetInfo(c *gin.Context) {
	req := request.GetPassportDetailInfoRequest{}
	if ok := request.Validate(c, &req); !ok {
		return
	}

	passportInfo, err := service.PassportDetailService.GetInfo(req.Passport)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SucceedWithData(c, "", passportInfo)
}

func (pdc *PassportDetailController) GetAccountList(c *gin.Context) {
	req := request.GetAccountListRequest{}
	if ok := request.Validate(c, &req); !ok {
		return
	}

	requestID := c.GetString("request_id")
	accounts, err := service.PassportDetailService.GetAccountListByLogin(requestID, req.Login)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SucceedWithData(c, "", accounts)
}

func (pdc *PassportDetailController) GetTradeList(c *gin.Context) {
	req := request.GetTradeListRequest{}
	if ok := request.Validate(c, &req); !ok {
		return
	}

	paging := &model.Paging{Limit: req.Limit, Offset: req.Offset, Sort: req.Sort, Order: req.Order}

	requestID := c.GetString("request_id")
	trades, err := service.PassportDetailService.GetTradeListByPassport(requestID, req.Passport, req.Login, req.Ticket, req.Symbol, paging)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SucceedWithData(c, "", trades)
}

func (pdc *PassportDetailController) GetSymbolList(c *gin.Context) {
	req := request.GetSymbolListRequest{}
	if ok := request.Validate(c, &req); !ok {
		return
	}

	requestID := c.GetString("request_id")
	symbols, err := service.PassportDetailService.GetSymbolListByPassport(requestID, req.Passport, req.Login)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SucceedWithData(c, "", symbols)
}
