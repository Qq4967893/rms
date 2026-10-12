package controller

import (
	"backend-rms/http/request"
	"backend-rms/http/response"
	"backend-rms/model"
	"backend-rms/service"
	"backend-rms/util/formatter"

	"github.com/gin-gonic/gin"
)

type MasterHistoryController struct{}

func (mhc *MasterHistoryController) GetList(c *gin.Context) {
	req := request.GetMasterHistoryListRequest{}
	if ok := request.Validate(c, &req); !ok {
		return
	}
	paging := &model.Paging{Limit: req.Limit, Offset: req.Offset, Sort: req.Sort, Order: req.Order}
	historicalRecords, totalCount, err := service.MasterService.GetMasterHistoricalList(req.MasterId, paging)
	if err != nil {
		response.Error(c, err)
		return
	}

	records, err := mhc.buildHistoricalRecords(historicalRecords)

	response.SucceedWithData(c, "", gin.H{
		"total_count":    totalCount,
		"master_records": records,
	})
}

func (mhc *MasterHistoryController) buildHistoricalRecords(historicalRecords []*service.MasterHistoricalRecord) (records []response.MasterHistoricalRecord, err error) {
	for _, m := range historicalRecords {
		endTime := ""
		if !m.EndTime.IsZero() {
			endTime = formatter.Datetime(m.EndTime)
		}
		records = append(records, response.MasterHistoricalRecord{
			Login:         m.Login,
			Commission:    formatter.Amount(m.Commission),
			ClosedPL:      formatter.Amount(m.ClosedPL),
			FollowerCount: m.FollowerCount,
			FollowedTimes: m.FollowedTimes,
			StartTime:     formatter.Datetime(m.StartTime),
			EndTime:       endTime,
		})
	}
	return
}
