package report

import (
	"backend-rms/http/request"
	"backend-rms/http/request/report"
	"backend-rms/http/response"
	"backend-rms/model"
	"backend-rms/service"
	"backend-rms/util"
	"backend-rms/util/download"
	"backend-rms/util/formatter"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

type IncomeController struct{}

func (ic *IncomeController) GetSymbolList(c *gin.Context) {
	req := report.GetIncomeSymbolListRequest{}
	if ok := request.Validate(c, &req); !ok {
		return
	}

	groups, err := service.GroupService.GetGroupsByBook(req.Book)
	if err != nil {
		response.Error(c, err)
		return
	}

	symbolList, err := service.DailyClosedTradeService.GetSymbolList(groups, req.Symbol, req.StartDate, req.EndDate)
	if err != nil {
		response.Error(c, err)
		return
	}

	if symbolList == nil {
		response.SucceedWithData(c, "", gin.H{
			"total_count": response.ReportIncomeSymbolResponse{},
			"symbol_list": []response.ReportIncomeSymbolResponse{},
		})
		return
	}

	err = util.DataOrderBy(symbolList, req.Sort, req.Order)
	if err != nil {
		response.Error(c, err)
		return
	}

	var totalProfit, totalSwaps, totalCommission, totalRawSpreadProfit, totalMarkUpSpreadProfit, lots decimal.Decimal
	result := make([]response.ReportIncomeSymbolResponse, 0)
	for _, v := range symbolList {
		middle := response.ReportIncomeSymbolResponse{
			Date:               req.StartDate.Format(formatter.DateOnly) + " ~ " + req.EndDate.Format(formatter.DateOnly),
			Symbol:             v.Symbol,
			Profit:             formatter.Amount(v.Profit),
			Swaps:              formatter.Amount(v.Swaps),
			Commission:         formatter.Amount(v.Commission),
			RawSpreadProfit:    formatter.Amount(v.RawSpreadProfit),
			MarkUpSpreadProfit: formatter.Amount(v.MarkUpSpreadProfit),
			Lots:               formatter.Amount(v.Lots),
		}
		result = append(result, middle)

		totalProfit = totalProfit.Add(v.Profit)
		totalSwaps = totalSwaps.Add(v.Swaps)
		totalCommission = totalCommission.Add(v.Commission)
		totalRawSpreadProfit = totalRawSpreadProfit.Add(v.RawSpreadProfit)
		totalMarkUpSpreadProfit = totalMarkUpSpreadProfit.Add(v.MarkUpSpreadProfit)
		lots = lots.Add(v.Lots)
	}

	totalCount := response.ReportIncomeSymbolResponse{
		Date:               req.StartDate.Format(formatter.DateOnly) + " ~ " + req.EndDate.Format(formatter.DateOnly),
		Symbol:             "Total",
		Profit:             formatter.Amount(totalProfit),
		Swaps:              formatter.Amount(totalSwaps),
		Commission:         formatter.Amount(totalCommission),
		RawSpreadProfit:    formatter.Amount(totalRawSpreadProfit),
		MarkUpSpreadProfit: formatter.Amount(totalMarkUpSpreadProfit),
		Lots:               formatter.Amount(lots),
	}

	if req.Export == "no" {
		response.SucceedWithData(c, "", gin.H{
			"total_count": totalCount,
			"symbol_list": result,
		})
		return
	}

	downloadData := make([][]string, 0)
	downloadTitle := []string{"Date", "Symbol", "Profit", "Swaps", "Commission", "Raw Spread Profit", "Mark Up Spread Profit", "Lots"}
	downloadData = append(downloadData, downloadTitle)
	downloadData = append(downloadData, []string{totalCount.Date, totalCount.Symbol, totalCount.Profit, totalCount.Swaps, totalCount.Commission, totalCount.RawSpreadProfit, totalCount.MarkUpSpreadProfit, totalCount.Lots})

	for _, v := range result {
		middle := []string{v.Date, v.Symbol, v.Profit, v.Swaps, v.Commission, v.RawSpreadProfit, v.MarkUpSpreadProfit, v.Lots}
		downloadData = append(downloadData, middle)
	}

	err = download.HandleHttpDownload(c, downloadData, "Income Breakdown Symbols")
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Succeed(c)
	return
}

func (ic *IncomeController) GetAccountList(c *gin.Context) {
	req := report.GetIncomeAccountListRequest{}
	if ok := request.Validate(c, &req); !ok {
		return
	}

	groups, err := service.GroupService.GetGroupsByBook(req.Book)
	if err != nil {
		response.Error(c, err)
		return
	}

	paging := &model.Paging{}
	if req.Export == "no" {
		paging = &model.Paging{Limit: req.Limit, Offset: req.Offset, Sort: req.Sort, Order: req.Order}
	} else {
		paging = nil
	}
	accountList, err := service.DailyClosedTradeService.GetAccountList(groups, req.Symbol, req.StartDate, req.EndDate, paging)
	if err != nil {
		response.Error(c, err)
		return
	}

	totalCount, err := service.DailyClosedTradeService.GetAccountCount(groups, req.Symbol, req.StartDate, req.EndDate)
	if err != nil {
		response.Error(c, err)
		return
	}

	result := make([]response.ReportIncomeAccountResponse, 0)
	for _, v := range accountList {
		middle := response.ReportIncomeAccountResponse{
			Login:              v.Login,
			Profit:             formatter.Amount(v.Profit),
			Swaps:              formatter.Amount(v.Swaps),
			Commission:         formatter.Amount(v.Commission),
			RawSpreadProfit:    formatter.Amount(v.RawSpreadProfit),
			MarkUpSpreadProfit: formatter.Amount(v.MarkUpSpreadProfit),
			Lots:               formatter.Amount(v.Lots),
		}
		result = append(result, middle)
	}

	if req.Export == "no" {
		response.SucceedWithData(c, "", gin.H{
			"total_count":  totalCount,
			"account_list": result,
		})
		return
	}

	downloadData := make([][]string, 0)
	downloadLite := []string{"Account", "Profit", "Swaps", "Commission", "Raw Spread Profit", "Mark Up Spread Profit", "Lots"}
	downloadData = append(downloadData, downloadLite)
	for _, v := range result {
		middle := []string{strconv.FormatInt(v.Login, 10), v.Profit, v.Swaps, v.Commission, v.RawSpreadProfit, v.MarkUpSpreadProfit, v.Lots}
		downloadData = append(downloadData, middle)
	}
	err = download.HandleHttpDownload(c, downloadData, "Income Breakdown Accounts")
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Succeed(c)
	return
}

func (ic *IncomeController) GetTradeList(c *gin.Context) {
	req := report.GetIncomeTradeListRequest{}
	if ok := request.Validate(c, &req); !ok {
		return
	}

	groups, err := service.GroupService.GetGroupsByBook(req.Book)
	if err != nil {
		response.Error(c, err)
		return
	}

	paging := &model.Paging{}
	if req.Export == "no" {
		paging = &model.Paging{Limit: req.Limit, Offset: req.Offset, Sort: req.Sort, Order: req.Order}
	} else {
		paging = nil
	}
	tradeList, err := service.DailyClosedTradeService.GetTradeList(groups, req.Symbol, req.StartDate, req.EndDate, paging)
	if err != nil {
		response.Error(c, err)
		return
	}

	totalCount, err := service.DailyClosedTradeService.GetTradeCount(groups, req.Symbol, req.StartDate, req.EndDate)
	if err != nil {
		response.Error(c, err)
		return
	}

	result := make([]response.ReportIncomeTradeResponse, 0)
	for _, v := range tradeList {
		middle := response.ReportIncomeTradeResponse{
			Login:              v.Login,
			Ticket:             v.Ticket,
			Profit:             formatter.Amount(v.Profit),
			Swaps:              formatter.Amount(v.Swaps),
			Commission:         formatter.Amount(v.Commission),
			RawSpreadProfit:    formatter.Amount(v.RawSpreadProfit),
			MarkUpSpreadProfit: formatter.Amount(v.MarkUpSpreadProfit),
			Lots:               formatter.Amount(v.Lots),
			OpenPrice:          v.OpenPrice,
			ClosePrice:         v.ClosePrice,
			OpenTime:           formatter.Datetime(v.OpenTime),
			CloseTime:          formatter.Datetime(v.CloseTime),
		}
		result = append(result, middle)
	}

	if req.Export == "no" {
		response.SucceedWithData(c, "", gin.H{
			"total_count": totalCount,
			"trade_list":  result,
		})
		return
	}
	downloadData := make([][]string, 0)
	downloadLite := []string{"Account", "Profit", "Swaps", "Commission", "Raw Spread Profit", "Mark Up Spread Profit", "Lots", "Open Price", "Close Price", "Open Time", "Close Time"}
	downloadData = append(downloadData, downloadLite)
	for _, v := range result {
		middle := []string{strconv.FormatInt(v.Login, 10), v.Profit, v.Swaps, v.Commission, v.RawSpreadProfit, v.MarkUpSpreadProfit, v.Lots, v.OpenPrice.String(), v.ClosePrice.String(), v.OpenTime, v.CloseTime}
		downloadData = append(downloadData, middle)
	}
	err = download.HandleHttpDownload(c, downloadData, "Income Breakdown Trades")
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Succeed(c)
	return
}
