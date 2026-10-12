package report

import (
	"backend-rms/http/request"
	"backend-rms/http/request/report"
	"backend-rms/http/response"
	"backend-rms/model"
	"backend-rms/service"
	"backend-rms/util"
	"backend-rms/util/formatter"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type TradeController struct{}

func (tc *TradeController) GetTradeList(c *gin.Context) {
	req := report.GetTradeListRequest{}
	if ok := request.Validate(c, &req); !ok {
		return
	}

	req.Book = strings.ToUpper(req.Book)
	if req.Sort == "lots" {
		req.Sort = "volume"
	}

	var login []int64
	var ticket int64
	switch req.SearchType {
	case "login":
		login = []int64{req.SearchText}
	case "passport":
		tradeAccountIDs, err := service.TradeAccountService.GetTradeAccountIDs(req.SearchText)
		if err != nil {
			response.Error(c, err)
			return
		}
		if len(tradeAccountIDs) == 0 {
			response.SucceedWithData(c, "", gin.H{
				"trades":      nil,
				"total_count": 0,
			})
			return
		}
		login = tradeAccountIDs
	case "ticket":
		ticket = req.SearchText
	}

	paging := model.Paging{Limit: req.Limit, Offset: req.Offset, Sort: req.Sort, Order: req.Order}

	groups, err := service.GroupService.GetGroupsByBook(req.Book)
	if err != nil {
		response.Error(c, err)
		return
	}

	totalCount, err := service.TradeService.GetReportTradeCount(req.TradeStatus, req.Symbol, groups, req.StartTime, req.EndTime, login, ticket)
	if err != nil {
		response.Error(c, err)
		return
	}

	tradeList, err := service.TradeService.GetReportTradeList(req.TradeStatus, req.Symbol, groups, req.StartTime, req.EndTime, login, ticket, paging)
	if err != nil {
		response.Error(c, err)
		return
	}

	trades, err := tc.buildTrades(tradeList)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SucceedWithData(c, "", gin.H{
		"trades":      trades,
		"total_count": totalCount,
	})
}

func (tc *TradeController) buildTrades(tradeList []*model.ReportTrade) (trades []response.ReportTrade, err error) {
	if len(tradeList) == 0 {
		return
	}

	groupMap, err := service.GroupService.GetAllMap()
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}

	var tickets []int64
	loginMap := make(map[int64]struct{})
	for _, trade := range tradeList {
		tickets = append(tickets, trade.Ticket)
		loginMap[trade.Login] = struct{}{}
	}

	openTypeMap, err := service.EventService.GetOpenTypeMapByTickets(tickets)
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}

	var logins []int64
	for login, _ := range loginMap {
		logins = append(logins, login)
	}

	loginPassportMap, err := service.TradeAccountService.GetLoginPassportMap(logins)
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}

	floatingPLMap, err := service.CacheService.GetTradeFloatingPLMap(tickets)
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}

	for _, trade := range tradeList {
		closeTime := ""
		closePrice := ""
		closeTypeDesc := ""
		closedPL := decimal.Zero
		if !trade.CloseTime.IsZero() {
			closeTime = formatter.Datetime(trade.CloseTime)
			closePrice = trade.ClosePrice.String()
			closedPL = trade.Profit.Add(trade.Swaps).Add(trade.Commission)
			closeTypeDesc = service.TradeService.GetCloseTypeDesc(trade.Type)
		}

		trades = append(trades, response.ReportTrade{
			Ticket:        trade.Ticket,
			Login:         trade.Login,
			Passport:      loginPassportMap[trade.Login],
			Book:          groupMap[trade.Group],
			CmdDesc:       service.TradeService.GetCmdDesc(trade.Cmd),
			Symbol:        trade.Symbol,
			Lots:          decimal.NewFromFloat(float64(trade.Volume) / 100).StringFixed(2),
			FloatingPL:    formatter.Amount(floatingPLMap[trade.Ticket]),
			ClosedPL:      formatter.Amount(closedPL),
			OpenPrice:     trade.OpenPrice.String(),
			OpenTime:      formatter.Datetime(trade.OpenTime),
			ClosePrice:    closePrice,
			CloseTime:     closeTime,
			SL:            trade.SL.String(),
			TP:            trade.TP.String(),
			OpenTypeDesc:  service.TradeService.GetOpenTypeDesc(openTypeMap[trade.Ticket]),
			CloseTypeDesc: closeTypeDesc,
		})
	}
	return
}
