package controller

import (
	"backend-rms/http/request"
	"backend-rms/http/response"
	"backend-rms/service"
	"backend-rms/util"
	"backend-rms/util/formatter"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

type SymbolController struct{}

func (sc *SymbolController) GetOriginalList(c *gin.Context) {
	originals, err := service.SymbolOriginalService.GetAllList(false)
	if err != nil {
		response.Error(c, err)
		return
	}

	var ret []string
	for _, o := range originals {
		ret = append(ret, o.Original)
	}
	response.SucceedWithData(c, "", ret)
}

func (sc *SymbolController) GetHoldingAccountList(c *gin.Context) {
	req := request.GetSymbolHoldingAccountListRequest{}
	if ok := request.Validate(c, &req); !ok {
		return
	}

	groups, err := service.GroupService.GetGroupsByBook(req.Book)
	if err != nil {
		response.Error(c, err)
		return
	}

	list, err := service.SymbolService.GetHoldingAccountList(req.Symbol, groups)
	if err != nil {
		response.Error(c, err)
		return
	}

	var totalBuyLots, totalSellLots, totalNetLots, totalSymbolFloatingPL, totalFloatingPL decimal.Decimal
	holdingAccountList := make([]response.SymbolHoldingAccount, 0)
	for _, v := range list {
		middle := response.SymbolHoldingAccount{
			Login:            v.Login,
			BuyLots:          formatter.Amount(v.BuyLots),
			SellLots:         formatter.Amount(v.SellLots),
			NetLots:          formatter.Amount(v.NetLots),
			SymbolFloatingPL: formatter.Amount(v.SymbolFloatingPL),
			FloatingPL:       formatter.Amount(v.FloatingPL),
		}
		holdingAccountList = append(holdingAccountList, middle)

		totalBuyLots = totalBuyLots.Add(v.BuyLots)
		totalSellLots = totalSellLots.Add(v.SellLots)
		totalNetLots = totalNetLots.Add(v.NetLots)
		totalSymbolFloatingPL = totalSymbolFloatingPL.Add(v.SymbolFloatingPL)
		totalFloatingPL = totalFloatingPL.Add(v.FloatingPL)
	}

	totalHoldingAccounts := response.TotalSymbolHoldingAccount{
		Total:            len(holdingAccountList),
		BuyLots:          formatter.Amount(totalBuyLots),
		SellLots:         formatter.Amount(totalSellLots),
		NetLots:          formatter.Amount(totalNetLots),
		SymbolFloatingPL: formatter.Amount(totalSymbolFloatingPL),
		FloatingPL:       formatter.Amount(totalFloatingPL),
	}

	err = util.DataOrderBy(holdingAccountList, "net_lots", "desc")
	if err != nil {
		response.Error(c, err)
		return
	}
	resultKey, err := util.GetResultKey(req.Offset, req.Limit, totalHoldingAccounts.Total)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.SucceedWithData(c, "", gin.H{
		"total_count": totalHoldingAccounts,
		"accounts":    holdingAccountList[resultKey["startKey"]:resultKey["endKey"]],
	})
}

func (sc *SymbolController) GetTradeList(c *gin.Context) {
	req := request.GetSymbolTradeListRequest{}
	if ok := request.Validate(c, &req); !ok {
		return
	}

	groups, err := service.GroupService.GetGroupsByBook(req.Book)
	if err != nil {
		response.Error(c, err)
		return
	}

	list, err := service.SymbolService.GetTradeList(req.Symbol, groups, req.Cmd)
	if err != nil {
		response.Error(c, err)
		return
	}

	tradeList := make([]response.SymbolTrade, 0)
	var lots, floatingPL decimal.Decimal
	for _, v := range list {
		var cmd string
		if v.Cmd == 0 {
			cmd = "Buy"
		} else if v.Cmd == 1 {
			cmd = "Sell"
		}

		lots = lots.Add(v.Lots)
		floatingPL = floatingPL.Add(v.FloatingPL)
		middle := response.SymbolTrade{
			Ticket:      v.Ticket,
			Login:       v.Login,
			OpenPrice:   formatter.Amount(v.OpenPrice),
			SymbolPrice: formatter.Amount(v.SymbolPrice),
			OpenTime:    formatter.Datetime(v.OpenTime),
			Cmd:         cmd,
			Lots:        formatter.Amount(v.Lots),
			FloatingPL:  formatter.Amount(v.FloatingPL),
			Profit:      formatter.Amount(v.Profit),
			Swap:        formatter.Amount(v.Swap),
			Commission:  formatter.Amount(v.Commission),
			SL:          formatter.Amount(v.SL),
			TP:          formatter.Amount(v.TP),
		}
		tradeList = append(tradeList, middle)
	}

	totalTrades := response.SymbolTotalTrade{
		Total:      len(list),
		Lots:       formatter.Amount(lots),
		FloatingPL: formatter.Amount(floatingPL),
	}

	resultKey, err := util.GetResultKey(req.Offset, req.Limit, len(list))
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SucceedWithData(c, "", gin.H{
		"total_count": totalTrades,
		"accounts":    tradeList[resultKey["startKey"]:resultKey["endKey"]],
	})
}

func (sc *SymbolController) GetPendingTradeList(c *gin.Context) {
	req := request.GetSymbolPendingTradeListRequest{}
	if ok := request.Validate(c, &req); !ok {
		return
	}

	groups, err := service.GroupService.GetGroupsByBook(req.Book)
	if err != nil {
		response.Error(c, err)
		return
	}

	list, err := service.SymbolService.GetPendingTradeList(req.Symbol, groups, req.Cmd)
	if err != nil {
		response.Error(c, err)
		return
	}

	quote, err := service.QuoteService.GetQuote(req.Symbol)
	if err != nil {
		response.Error(c, err)
		return
	}

	pendingTradeList := make([]response.SymbolPendingTrade, 0)
	var lots decimal.Decimal
	for _, v := range list {
		var cmd string
		if v.Cmd == 2 || v.Cmd == 4 {
			cmd = "Buy"
		} else {
			cmd = "Sell"
		}

		var symbolPrice decimal.Decimal
		if req.Cmd == 0 {
			symbolPrice = quote.Bid
		} else if req.Cmd == 1 {
			symbolPrice = quote.Ask
		}

		lots = lots.Add(decimal.NewFromInt(v.Volume).Div(decimal.NewFromInt(100)))
		middle := response.SymbolPendingTrade{
			Ticket:       v.Ticket,
			Login:        v.Login,
			PendingTime:  formatter.Datetime(v.OpenTime),
			PendingPrice: formatter.Amount(v.OpenPrice),
			SymbolPrice:  formatter.Amount(symbolPrice),
			Cmd:          cmd,
			Lots:         formatter.Amount(decimal.NewFromInt(v.Volume).Div(decimal.NewFromInt(100))),
		}
		pendingTradeList = append(pendingTradeList, middle)
	}

	totalPendingTrades := response.SymbolTotalPendingTrade{
		Total: len(list),
		Lots:  formatter.Amount(lots),
	}

	resultKey, err := util.GetResultKey(req.Offset, req.Limit, len(list))
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SucceedWithData(c, "", gin.H{
		"total_count": totalPendingTrades,
		"accounts":    pendingTradeList[resultKey["startKey"]:resultKey["endKey"]],
	})
}
