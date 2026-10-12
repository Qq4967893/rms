package report

import (
	"backend-rms/service"
	"backend-rms/util/formatter"
	"backend-rms/ws"
	"backend-rms/ws/controller"
	"backend-rms/ws/request"
	"backend-rms/ws/request/report"
	"backend-rms/ws/response"
	"time"

	"github.com/shopspring/decimal"
)

type TradeController struct{}

func (tc *TradeController) GetTradeList(c *ws.Context) {
	controller.TickerHandle(c, 3*time.Second, func(c *ws.Context) (data any, err error) {
		req := report.GetTradeListRequest{}
		if ok := request.Validate(c, &req); !ok {
			return
		}

		return tc.buildTradeList(req.Params.Tickets)
	})
}

func (tc *TradeController) buildTradeList(tickets []int64) ([]*response.ReportTrade, error) {
	groups, err := service.GroupService.GetGroupsByBook("ALL")
	if err != nil {
		return nil, err
	}

	tradeList, err := service.TradeService.GetReportTradeListByTickets(tickets, groups)
	if err != nil {
		return nil, err
	}

	floatingPLMap, err := service.CacheService.GetTradeFloatingPLMap(tickets)
	if err != nil {
		return nil, err
	}

	var trades []*response.ReportTrade
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

		trades = append(trades, &response.ReportTrade{
			Ticket:        trade.Ticket,
			FloatingPL:    formatter.Amount(floatingPLMap[trade.Ticket]),
			ClosedPL:      formatter.Amount(closedPL),
			ClosePrice:    closePrice,
			CloseTime:     closeTime,
			CloseTypeDesc: closeTypeDesc,
		})
	}

	return trades, nil
}
