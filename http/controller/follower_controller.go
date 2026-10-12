package controller

import (
	"backend-rms/http/request"
	"backend-rms/http/response"
	"backend-rms/model"
	"backend-rms/service"
	"backend-rms/util"
	"backend-rms/util/formatter"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-module/carbon/v2"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type FollowerController struct{}

func (fc *FollowerController) GetList(c *gin.Context) {
	req := request.GetFollowerListRequest{}
	if ok := request.Validate(c, &req); !ok {
		return
	}

	groups, err := service.GroupService.GetGroupsByBook(req.Book)
	if err != nil {
		response.Error(c, err)
		return
	}

	requestID := c.GetString("request_id")
	followers, err := fc.buildFollowerList(req.MasterId, req.StartTime, req.EndTime, groups, req.FollowThanOne, req.IsFollowing, requestID)
	if err != nil {
		response.Error(c, err)
		return
	}

	totalCount := len(followers)
	if totalCount == 0 {
		response.SucceedWithData(c, "", gin.H{
			"total_count": totalCount,
			"followers":   followers,
		})
		return
	}

	err = util.DataOrderBy(followers, req.Sort, req.Order)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.SucceedWithData(c, "", gin.H{
		"total_count": totalCount,
		"followers":   followers,
	})
}

func (fc *FollowerController) GetHistoryList(c *gin.Context) {
	req := request.GetFollowerHistoryListRequest{}
	if ok := request.Validate(c, &req); !ok {
		return
	}

	copyRecords, err := service.CopyRecordService.GetAll(req.FollowerId, req.MasterId)
	if err != nil {
		response.Error(c, err)
		return
	}

	copyPays, err := service.CopyPayService.GetAll(req.FollowerId, req.MasterId)
	if err != nil {
		response.Error(c, err)
		return
	}

	followerRecords := make([]response.FollowerRecord, 0)
	for _, v := range copyRecords {
		var commission decimal.Decimal
		for _, vv := range copyPays {
			if (v.StartTime.Before(vv.PayTime) && v.EndTime.After(vv.PayTime)) || (v.StartTime.Before(vv.PayTime) && v.EndTime.IsZero()) {
				commission = commission.Add(vv.PayMaster)
			}
		}
		startTime := v.StartTime.Format(formatter.DateTime)
		var endTime string
		endTimeByClosedPl := v.EndTime
		if v.EndTime.IsZero() {
			endTime = "-"
			endTimeByClosedPl = time.Now()
		} else {
			endTime = v.EndTime.Format(formatter.DateTime)
		}

		var closePl decimal.Decimal
		closePls, err := service.TradeService.GetAccountClosedPL([]int64{v.Mt4To}, []string{}, v.StartTime, endTimeByClosedPl, "")
		if err != nil {
			response.Error(c, err)
			return
		}
		if len(closePls) > 0 {
			closePl = closePls[v.Mt4To]
		}
		middle := response.FollowerRecord{
			StartTime:  startTime,
			EndTime:    endTime,
			Commission: commission,
			ClosedPl:   closePl,
		}

		followerRecords = append(followerRecords, middle)
	}

	totalCount := len(followerRecords)
	if totalCount == 0 {
		response.SucceedWithData(c, "", gin.H{
			"total_count":     totalCount,
			"followerRecords": followerRecords,
		})
		return
	}

	err = util.DataOrderBy(followerRecords, req.Sort, req.Order)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.SucceedWithData(c, "", gin.H{
		"total_count":     totalCount,
		"followerRecords": followerRecords,
	})
}

func (fc *FollowerController) buildFollowerList(masterId int64, startTime, endTime time.Time, groups []string, followThanOne, isFollowing int, requestID string) (followerList []response.Follower, err error) {
	copies, err := service.CopyRecordService.GetCopies(masterId, groups)
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}

	if len(copies) == 0 {
		return
	}

	followers := make([]int64, 0)
	for _, v := range copies {
		followers = append(followers, v.Mt4To)
	}
	followMasterNums, err := service.CopyRecordService.GetPointFollowMasterNums(endTime, followers, groups)
	if err != nil {
		return
	}

	copiedAmounts, err := service.CopyRecordService.GetPointCopiedAmounts(endTime, followers, groups, masterId)
	if err != nil {
		return
	}

	copiedPLInfos, err := service.CopyRecordService.GetCopiedPLInfos(groups, followers, startTime, endTime, requestID)
	if err != nil {
		return
	}

	carbonEndTime := carbon.Time2Carbon(endTime)
	yesterday := carbonEndTime.Yesterday().Carbon2Time()
	today := model.GetDailyDateTime(yesterday).AddSeconds(1).ToDateTimeString()

	yesterdayBalances, err := service.ReportDailyService.GetLoginBalanceMapByTime(yesterday, followers)
	if err != nil {
		return
	}

	todayNetCashFlowMap, err := service.CashFlowService.GetNetCashFlowMap(followers, today, carbonEndTime.ToDateTimeString())
	if err != nil {
		return
	}

	todayClosedPLMap, err := service.TradeService.GetClosedPLMap(followers, today, carbonEndTime.ToDateTimeString())
	if err != nil {
		return
	}

	followingLogins, err := service.CopyRecordService.GetFollowingLogins(masterId, groups)
	if err != nil {
		return
	}

	for _, v := range copies {
		login := v.Mt4To

		followMasterNum := 0
		if _, ok := followMasterNums[login]; ok {
			followMasterNum = followMasterNums[login]
		}

		if followThanOne != -1 {
			if followThanOne >= 1 && followMasterNum <= 1 {
				continue
			} else if followThanOne < 1 && followMasterNum > 1 {
				continue
			}
		}

		if isFollowing != -1 {
			if isFollowing == 1 && !util.InSlice(login, followingLogins) {
				continue
			}
			if isFollowing == 0 && util.InSlice(login, followingLogins) {
				continue
			}
		}

		following := formatter.Bool(util.InSlice(login, followingLogins))
		followerInfo := response.Follower{
			Mt4Id:           login,
			UserCode:        v.UserCode,
			PeriodPl:        formatter.Amount(copiedPLInfos[login].PeriodPl),
			ClosedPl:        formatter.Amount(copiedPLInfos[login].ClosedPl),
			FloatingPl:      formatter.Amount(copiedPLInfos[login].FloatingPl),
			Balance:         formatter.Amount(yesterdayBalances[login].Add(todayNetCashFlowMap[login]).Add(todayClosedPLMap[login])),
			MarginLevel:     formatter.Amount(copiedPLInfos[login].MarginLevel),
			FollowMasterNum: followMasterNum,
			CopiedAmount:    formatter.Amount(copiedAmounts[login].Amount),
			CopiedAllAmount: formatter.Amount(copiedAmounts[login].OtherMasterAmount.Add(copiedAmounts[login].Amount)),
			Group:           v.Group,
			IsFollowing:     following,
			StartTime:       formatter.Datetime(v.StartTime),
		}
		followerList = append(followerList, followerInfo)
	}
	return
}
