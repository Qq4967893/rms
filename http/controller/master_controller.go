package controller

import (
	"backend-rms/http/request"
	"backend-rms/http/response"
	"backend-rms/model"
	"backend-rms/service"
	"backend-rms/util"
	"backend-rms/util/formatter"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type MasterController struct{}

func (mc *MasterController) GetList(c *gin.Context) {
	req := request.GetMasterListRequest{}
	if ok := request.Validate(c, &req); !ok {
		return
	}

	groups, err := service.GroupService.GetGroupsByBook(req.Book)
	if err != nil {
		response.Error(c, err)
		return
	}

	login := req.SearchText
	paging := &model.Paging{Limit: req.Limit, Offset: req.Offset, Sort: req.Sort, Order: req.Order}
	requestID := c.GetString("request_id")

	masterList, err := service.MasterService.GetMasterList(requestID, req.IsPresent, groups, login, req.StartTime, req.EndTime)
	if err != nil {
		response.Error(c, err)
		return
	}

	totalCount := len(masterList)
	sort.Slice(masterList, func(i, j int) bool {
		switch paging.Sort {
		case "period_pl":
			if paging.Order == "asc" {
				return masterList[i].NetPeriodPL.Cmp(masterList[j].NetPeriodPL) < 0
			} else {
				return masterList[i].NetPeriodPL.Cmp(masterList[j].NetPeriodPL) > 0
			}
		case "closed_pl":
			if paging.Order == "asc" {
				return masterList[i].PeriodClosedPL.Cmp(masterList[j].PeriodClosedPL) < 0
			} else {
				return masterList[i].PeriodClosedPL.Cmp(masterList[j].PeriodClosedPL) > 0
			}
		case "floating_pl":
			if paging.Order == "asc" {
				return masterList[i].FloatingPL.Cmp(masterList[j].FloatingPL) < 0
			} else {
				return masterList[i].FloatingPL.Cmp(masterList[j].FloatingPL) > 0
			}
		case "follower_count":
			if paging.Order == "asc" {
				return masterList[i].FollowerCount < masterList[j].FollowerCount
			} else {
				return masterList[i].FollowerCount > masterList[j].FollowerCount
			}
		case "funds":
			if paging.Order == "asc" {
				return masterList[i].CopiedAmount.Cmp(masterList[j].CopiedAmount) < 0
			} else {
				return masterList[i].CopiedAmount.Cmp(masterList[j].CopiedAmount) > 0
			}
		case "margin_level":
			if paging.Order == "asc" {
				return masterList[i].MarginLevel.Cmp(masterList[j].MarginLevel) < 0
			} else {
				return masterList[i].MarginLevel.Cmp(masterList[j].MarginLevel) > 0
			}
		default:
			return false
		}
	})

	start := paging.Offset
	if start > len(masterList) {
		start = len(masterList)
	}

	end := paging.Offset + paging.Limit
	if end > len(masterList) {
		end = len(masterList)
	}

	masterList = masterList[start:end]
	masters, err := mc.buildMasters(masterList)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SucceedWithData(c, "", gin.H{
		"total_count": totalCount,
		"masters":     masters,
	})
}

func (mc *MasterController) buildMasters(masterList []*service.Master) (masters []response.Master, err error) {
	groupMap, err := service.GroupService.GetAllMap()
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}

	for _, m := range masterList {
		periodPL := formatter.Amount(m.NetPeriodPL)
		if m.IsLimited {
			periodPL = "*"
		}
		masters = append(masters, response.Master{
			Login:         m.Login,
			Passport:      m.Passport,
			PeriodPL:      periodPL,
			ClosedPL:      formatter.Amount(m.PeriodClosedPL),
			FloatingPL:    formatter.Amount(m.FloatingPL),
			Balance:       formatter.Amount(m.Balance),
			FollowerCount: m.FollowerCount,
			CopiedAmount:  formatter.Amount(m.CopiedAmount),
			MarginLevel:   formatter.Amount(m.MarginLevel),
			Book:          groupMap[m.Group],
			IsPresent:     formatter.Bool(m.IsPresent),
		})
	}
	return
}

func (mc *MasterController) Download(c *gin.Context) {
}
