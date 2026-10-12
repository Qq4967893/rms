package controller

import (
	"backend-rms/http/request"
	"backend-rms/http/response"
	"backend-rms/model"
	"backend-rms/service"
	"backend-rms/util/formatter"
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
)

type RestrictedPassportController struct{}

func (rc *RestrictedPassportController) GetList(c *gin.Context) {
	req := request.GetRestrictedPassportListRequest{}
	if ok := request.Validate(c, &req); !ok {
		return
	}

	paging := &model.Paging{Limit: req.Limit, Offset: req.Offset, Sort: req.Sort, Order: req.Order}
	restrictedPassports, err := service.TradingLimitService.GetList(0, model.TRADING_LIMIT, paging)
	if err != nil {
		response.Error(c, err)
		return
	}

	result := make([]response.GetRestrictedPassportList, 0)
	for _, v := range restrictedPassports {
		middle := response.GetRestrictedPassportList{
			Id:           v.Id,
			Passport:     v.Passport,
			Symbol:       v.Symbol,
			TriggerValue: v.TriggerValue,
			LimitValue:   v.LimitValue,
			CurrentValue: v.CurrentValue,
			TriggerTime:  formatter.Datetime(v.TriggerTime),
		}
		result = append(result, middle)
	}
	response.SucceedWithData(c, "", result)
	return
}

func (rc *RestrictedPassportController) GetRecordList(c *gin.Context) {
	req := request.GetRestrictedPassportRecordListRequest{}
	if ok := request.Validate(c, &req); !ok {
		return
	}

	paging := &model.Paging{Limit: req.Limit, Offset: req.Offset, Sort: req.Sort, Order: req.Order}
	restrictedPassports, err := service.TradingLimitService.GetList(req.Passport, model.REMOVE_TRADING_LIMIT, paging)
	if err != nil {
		response.Error(c, err)
		return
	}
	result := make([]response.GetRestrictedPassportRecord, 0)
	for _, v := range restrictedPassports {
		middle := response.GetRestrictedPassportRecord{
			Passport:     v.Passport,
			Symbol:       v.Symbol,
			TriggerValue: v.TriggerValue,
			LimitValue:   v.LimitValue,
			RemoveValue:  v.CurrentValue,
			TriggerTime:  formatter.Datetime(v.TriggerTime),
			RemoveTime:   formatter.Datetime(v.RemoveTime),
			Operator:     v.Operator,
		}
		result = append(result, middle)
	}
	response.SucceedWithData(c, "", result)
	return
}

func (rc *RestrictedPassportController) Remove(c *gin.Context) {
	req := request.RemoveRestrictedPassportRequest{}
	if ok := request.Validate(c, &req); !ok {
		return
	}

	restrictedInfo, err := service.TradingLimitService.GetInfoById(req.Id)
	if err != nil {
		response.Error(c, err)
		return
	}

	//判断当前通行证品种是否限制交易
	if restrictedInfo.RemoveTime.Format(formatter.DateTime) != "0001-01-01 00:00:00" {
		response.FailWithMessage(c, "The passport symbol is not limit")
		return
	}

	if restrictedInfo.CurrentValue >= restrictedInfo.LimitValue {
		response.FailWithMessage(c, "The user's Exposure value is currently too high to perform this operation")
		return
	}

	loginManagerID := c.GetInt64("login_manager_id")
	managerInfo, err := service.ManagerService.GetInfoById(loginManagerID)
	if err != nil {
		response.Error(c, err)
		return
	}
	err = service.TradingLimitService.RemoveRestricted(req.Id, restrictedInfo.CurrentValue, time.Now(), managerInfo.Username)
	if err != nil {
		response.Error(c, err)
		return
	}

	content, _ := json.Marshal(map[string][]any{
		"title":   {},
		"content": {},
	})

	notificationParams := model.Notification{
		Title:    "api_exposure_symbol_cancel",
		Content:  string(content),
		Format:   3,
		Type:     1,
		Resource: "system",
		Partner:  "lite",
	}
	notificationId, err := service.NotificationService.Add(notificationParams)
	if err != nil {
		response.Error(c, err)
		return
	}

	notifyInboxParams := model.NotifyInbox{
		UserCode: restrictedInfo.Passport,
		Mid:      notificationId,
		Title:    "api_exposure_symbol_cancel",
		Content:  string(content),
		Format:   3,
		Type:     1,
		Resource: "system",
	}
	err = service.NotifyInboxService.Add(notifyInboxParams)
	if err != nil {
		response.Error(c, err)
		return
	}

	service.EmailService.SendTemplate("exposure_symbol_cancel", restrictedInfo.Email, map[string]string{"user_name": restrictedInfo.Username}, 0, restrictedInfo.WorldCode, 2)
	response.Succeed(c)
	return
}
