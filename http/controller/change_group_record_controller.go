package controller

import (
	"backend-rms/http/request"
	"backend-rms/http/response"
	"backend-rms/model"
	"backend-rms/service"
	"backend-rms/util"
	"backend-rms/util/formatter"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type ChangeGroupRecordController struct{}

func (cc *ChangeGroupRecordController) GetList(c *gin.Context) {
	req := request.GetChangeGroupRecordListRequest{}
	if ok := request.Validate(c, &req); !ok {
		return
	}

	paging := &model.Paging{Limit: req.Limit, Offset: req.Offset, Sort: "created_at", Order: "desc"}

	logins := make([]int64, 0)
	if req.Login != 0 {
		logins = []int64{req.Login}
	}
	changeGroupRecords, err := service.ChangeGroupRecordService.GetList(logins, req.ManagerId, paging)
	if err != nil {
		response.Error(c, err)
		return
	}
	totalCount, err := service.ChangeGroupRecordService.GetCount(logins, req.ManagerId)
	if err != nil {
		response.Error(c, err)
		return
	}

	list, err := cc.buildList(changeGroupRecords)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.SucceedWithData(c, "", gin.H{
		"total_count":          totalCount,
		"change_group_records": list,
	})
}

func (cc *ChangeGroupRecordController) GetOperators(c *gin.Context) {
	operators, err := service.ChangeGroupRecordService.GetOperators()
	if err != nil {
		response.Error(c, err)
		return
	}
	response.SucceedWithData(c, "", operators)
}

func (cc *ChangeGroupRecordController) buildList(changeGroupRecords []*model.ChangeGroupRecord) (list []response.ChangeGroupRecordList, err error) {
	groupMaps, err := service.GroupService.GetAllMap()
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}

	for _, changeGroupRecord := range changeGroupRecords {
		operate := groupMaps[changeGroupRecord.BeforeGroup] + "-" + groupMaps[changeGroupRecord.AfterGroup]

		hedge := "YES"
		if changeGroupRecord.Status == 0 {
			hedge = "NO"
		}

		managerInfo, err := service.ManagerService.GetInfoById(changeGroupRecord.ManagerId)
		if err != nil {
			return nil, errors.WithMessage(err, util.ErrWrapMessage())
		}

		middle := response.ChangeGroupRecordList{
			Login:      changeGroupRecord.Login,
			Operate:    operate,
			GroupOld:   changeGroupRecord.BeforeGroup,
			TransferTo: changeGroupRecord.AfterGroup,
			Hedge:      hedge,
			Time:       formatter.Datetime(changeGroupRecord.CreatedAt),
			Operator:   managerInfo.Username,
		}
		list = append(list, middle)
	}
	return
}
