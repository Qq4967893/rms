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
	"github.com/shopspring/decimal"
)

type HedgeTradeController struct{}

func (hc *HedgeTradeController) GetList(c *gin.Context) {
	req := request.GetHedgeTradeList{}
	if ok := request.Validate(c, &req); !ok {
		return
	}

	paging := &model.Paging{Limit: req.Limit, Offset: req.Offset, Sort: "created_at", Order: "desc"}
	hedgeTrades, err := service.HedgeTradeService.GetList(req.Login, req.ManagerId, paging)
	if err != nil {
		response.Error(c, err)
		return
	}

	totalCount, err := service.HedgeTradeService.GetCount(req.Login, req.ManagerId)
	if err != nil {
		response.Error(c, err)
		return
	}

	list, err := hc.buildList(hedgeTrades)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SucceedWithData(c, "", gin.H{
		"total_count":  totalCount,
		"hedge_trades": list,
	})
}

func (hc *HedgeTradeController) buildList(hedgeTrades []*model.HedgeTrade) (list []response.HedgeTradeList, err error) {
	for _, hedgeTrade := range hedgeTrades {
		managerInfo, err := service.ManagerService.GetInfoById(hedgeTrade.ManagerId)
		if err != nil {
			return nil, errors.WithMessage(err, util.ErrWrapMessage())
		}

		cmd := "Sell"
		if hedgeTrade.Cmd == 0 {
			cmd = "Buy"
		}

		lot := decimal.NewFromInt(hedgeTrade.Volume).Div(decimal.NewFromInt(100))
		middle := response.HedgeTradeList{
			Account:    hedgeTrade.Login,
			HedgeOrder: hedgeTrade.Ticket,
			OpenTime:   formatter.Datetime(hedgeTrade.OpenTime),
			Market:     hedgeTrade.Symbol,
			Cmd:        cmd,
			Lot:        formatter.Amount(lot),
			Operator:   managerInfo.Username,
		}
		list = append(list, middle)
	}
	return
}
