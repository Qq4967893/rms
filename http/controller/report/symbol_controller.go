package report

import (
	"backend-rms/http/request"
	"backend-rms/http/request/report"
	"backend-rms/http/response"
	"backend-rms/service"
	"backend-rms/util"
	"backend-rms/util/formatter"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"golang.org/x/exp/maps"
	"golang.org/x/sync/errgroup"
)

type SymbolController struct{}

type symbolOpeningTrades struct {
	FloatingPL, Profit, Swap, Commission decimal.Decimal
}

func (stc *SymbolController) GetList(c *gin.Context) {
	req := report.GetSymbolListRequest{}
	if ok := request.Validate(c, &req); !ok {
		return
	}

	groups, err := service.GroupService.GetGroupsByBook(req.Book)
	if err != nil {
		response.Error(c, err)
		return
	}

	amount, err := stc.buildAmountInfo(groups, req.StartTime, req.EndTime)
	if err != nil {
		response.Error(c, err)
		return
	}

	symbolTrades, err := stc.buildSymbolTradesInfo(groups, req.Symbol, req.StartTime, req.EndTime, c.GetString("request_id"))
	if err != nil {
		response.Error(c, err)
		return
	}

	symbolSummary := response.ReportSymbolSummary{
		Date:               req.StartTime.Format(formatter.DateTime) + "~" + req.EndTime.Format(formatter.DateTime),
		Symbol:             "Total",
		CashInflow:         amount["cashInflow"],
		CashOutflow:        amount["cashOutflow"],
		NetCashFlow:        amount["netCashFlow"],
		EndTimeBalance:     amount["endTimeBalance"],
		EndTimeEquity:      decimal.Zero,
		ClosedPL:           decimal.Zero,
		ClosedProfit:       decimal.Zero,
		ClosedSwap:         decimal.Zero,
		ClosedCommission:   decimal.Zero,
		FloatingPL:         decimal.Zero,
		FloatingProfit:     decimal.Zero,
		FloatingSwap:       decimal.Zero,
		FloatingCommission: decimal.Zero,
		PeriodPL:           decimal.Zero,
	}

	for _, v := range symbolTrades {
		symbolSummary.ClosedPL = symbolSummary.ClosedPL.Add(v.ClosedPL)
		symbolSummary.ClosedProfit = symbolSummary.ClosedProfit.Add(v.ClosedProfit)
		symbolSummary.ClosedSwap = symbolSummary.ClosedSwap.Add(v.ClosedSwap)
		symbolSummary.ClosedCommission = symbolSummary.ClosedCommission.Add(v.ClosedCommission)
		symbolSummary.FloatingPL = symbolSummary.FloatingPL.Add(v.FloatingPL)
		symbolSummary.FloatingProfit = symbolSummary.FloatingProfit.Add(v.FloatingProfit)
		symbolSummary.FloatingSwap = symbolSummary.FloatingSwap.Add(v.FloatingSwap)
		symbolSummary.FloatingCommission = symbolSummary.FloatingCommission.Add(v.FloatingCommission)
		symbolSummary.PeriodPL = symbolSummary.PeriodPL.Add(v.PeriodPL)
	}

	err = util.DataOrderBy(symbolTrades, req.Sort, req.Order)
	if err != nil {
		response.Error(c, err)
		return
	}

	symbolSummaryResult := response.ReportSymbolSummaryResult{
		Date:               req.StartTime.Format(formatter.DateTime) + "~" + req.EndTime.Format(formatter.DateTime),
		Symbol:             "Total",
		CashInflow:         formatter.Amount(amount["cashInflow"]),
		CashOutflow:        formatter.Amount(amount["cashOutflow"]),
		NetCashFlow:        formatter.Amount(amount["netCashFlow"]),
		EndTimeBalance:     formatter.Amount(amount["endTimeBalance"]),
		EndTimeEquity:      formatter.Amount(amount["endTimeBalance"].Add(symbolSummary.FloatingPL)),
		ClosedPL:           formatter.Amount(symbolSummary.ClosedPL),
		ClosedProfit:       formatter.Amount(symbolSummary.ClosedProfit),
		ClosedSwap:         formatter.Amount(symbolSummary.ClosedSwap),
		ClosedCommission:   formatter.Amount(symbolSummary.ClosedCommission),
		FloatingPL:         formatter.Amount(symbolSummary.FloatingPL),
		FloatingProfit:     formatter.Amount(symbolSummary.FloatingProfit),
		FloatingSwap:       formatter.Amount(symbolSummary.FloatingSwap),
		FloatingCommission: formatter.Amount(symbolSummary.FloatingCommission),
		PeriodPL:           formatter.Amount(symbolSummary.PeriodPL),
	}

	var symbolTradesResult []response.ReportSymbolGetListResult
	for _, v := range symbolTrades {
		symbolTrade := response.ReportSymbolGetListResult{
			Date:               req.StartTime.Format(formatter.DateTime) + "~" + req.EndTime.Format(formatter.DateTime),
			Symbol:             v.Symbol,
			ClosedPL:           formatter.Amount(v.ClosedPL),
			ClosedProfit:       formatter.Amount(v.ClosedProfit),
			ClosedSwap:         formatter.Amount(v.ClosedSwap),
			ClosedCommission:   formatter.Amount(v.ClosedCommission),
			FloatingPL:         formatter.Amount(v.FloatingPL),
			FloatingProfit:     formatter.Amount(v.FloatingProfit),
			FloatingSwap:       formatter.Amount(v.FloatingSwap),
			FloatingCommission: formatter.Amount(v.FloatingCommission),
			PeriodPL:           formatter.Amount(v.PeriodPL),
		}
		symbolTradesResult = append(symbolTradesResult, symbolTrade)
	}

	response.SucceedWithData(c, "", gin.H{
		"summary": symbolSummaryResult,
		"symbols": symbolTradesResult,
	})
}

func (stc *SymbolController) GetSymbolAccountList(c *gin.Context) {
	req := report.GetSymbolAccountListRequest{}
	if ok := request.Validate(c, &req); !ok {
		return
	}

	groups, err := service.GroupService.GetGroupsByBook(req.Book)
	if err != nil {
		response.Error(c, err)
		return
	}

	accountClosedPL, err := service.TradeService.GetAccountClosedPL([]int64{}, groups, req.StartTime, req.EndTime, req.Symbol)
	if err != nil {
		response.Error(c, err)
		return
	}
	closedPLAccounts := maps.Keys(accountClosedPL)

	//获取历史时间点浮动盈亏数据（start_time_floating_pl,end_time_floating_pl,period_pl,period_floating_pl）
	startTimeOpeningTrades, err := service.TradeService.GetPointOpeningTrades(groups, req.StartTime, req.Symbol)
	if err != nil {
		response.Error(c, err)
		return
	}
	accountStartTimeFloatingPL, err := service.FloatingPLService.GetPointAccountFloatingPL(startTimeOpeningTrades, req.StartTime, c.GetString("request_id"))
	if err != nil {
		response.Error(c, err)
		return
	}
	startTimeFloatingPLAccounts := maps.Keys(accountStartTimeFloatingPL)

	endTimeOpeningTrades, err := service.TradeService.GetPointOpeningTrades(groups, req.EndTime, req.Symbol)
	if err != nil {
		response.Error(c, err)
		return
	}
	accountEndTimeFloatingPL, err := service.FloatingPLService.GetPointAccountFloatingPL(endTimeOpeningTrades, req.EndTime, c.GetString("request_id"))
	if err != nil {
		response.Error(c, err)
		return
	}
	endTimeFloatingPLAccounts := maps.Keys(accountEndTimeFloatingPL)

	//合并数组
	symbolAccounts := append(append(closedPLAccounts, startTimeFloatingPLAccounts...), endTimeFloatingPLAccounts...)

	//对symbolAccounts做去重操作
	var sliceLen int = len(symbolAccounts)
	for i := 0; i < sliceLen; i++ {
		for j := i + 1; j < sliceLen; j++ {
			if symbolAccounts[i] == symbolAccounts[j] {
				symbolAccounts = append(symbolAccounts[:i], symbolAccounts[i+1:]...)
				sliceLen--
				i--
				break
			}
		}
	}
	symbolAccounts = symbolAccounts[:sliceLen]

	var symbolAccountList []response.ReportSymbolGetSymbolAccountList
	for _, v := range symbolAccounts {
		symbolAccount := response.ReportSymbolGetSymbolAccountList{
			Account:             v,
			PeriodPL:            accountEndTimeFloatingPL[v].Sub(accountStartTimeFloatingPL[v]).Add(accountClosedPL[v]),
			StartTimeFloatingPL: accountStartTimeFloatingPL[v],
			EndTimeFloatingPL:   accountEndTimeFloatingPL[v],
			ClosedPL:            accountClosedPL[v],
			PeriodFloatingPL:    accountEndTimeFloatingPL[v].Sub(accountStartTimeFloatingPL[v]),
		}
		symbolAccountList = append(symbolAccountList, symbolAccount)
	}
	totalCount := len(symbolAccountList)

	err = util.DataOrderBy(symbolAccountList, req.Sort, req.Order)
	if err != nil {
		response.Error(c, err)
		return
	}
	resultKey, err := util.GetResultKey(req.Offset, req.Limit, len(symbolAccountList))
	if err != nil {
		response.Error(c, err)
		return
	}

	var symbolAccountListResult []response.ReportSymbolGetSymbolAccountListResult
	for _, v := range symbolAccountList[resultKey["startKey"]:resultKey["endKey"]] {
		symbolAccountInfo := response.ReportSymbolGetSymbolAccountListResult{
			Account:             v.Account,
			PeriodPL:            formatter.Amount(v.PeriodPL),
			StartTimeFloatingPL: formatter.Amount(v.StartTimeFloatingPL),
			EndTimeFloatingPL:   formatter.Amount(v.EndTimeFloatingPL),
			ClosedPL:            formatter.Amount(v.ClosedPL),
			PeriodFloatingPL:    formatter.Amount(v.PeriodFloatingPL),
		}
		symbolAccountListResult = append(symbolAccountListResult, symbolAccountInfo)
	}

	response.SucceedWithData(c, "", gin.H{
		"symbol_account": symbolAccountListResult,
		"total_count":    totalCount,
	})
}

func (stc *SymbolController) buildSymbolTradesInfo(groups []string, symbol string, startTime, endTime time.Time, requestID string) ([]response.ReportSymbolGetList, error) {
	g := new(errgroup.Group)

	var symbolClosedPL map[string]map[string]decimal.Decimal
	var err error
	g.Go(func() error {
		symbolClosedPL, err = service.TradeService.GetSymbolClosedPL(groups, startTime, endTime)
		return errors.WithMessage(err, util.ErrWrapMessage())
	})

	var startTimeSymbolOpeningTrades map[string]symbolOpeningTrades
	g.Go(func() error {
		startTimeSymbolOpeningTrades, err = stc.buildSymbolOpeningTrades(groups, startTime, requestID)
		return errors.WithMessage(err, util.ErrWrapMessage())
	})

	var endTimeSymbolOpeningTrades map[string]symbolOpeningTrades
	g.Go(func() error {
		endTimeSymbolOpeningTrades, err = stc.buildSymbolOpeningTrades(groups, endTime, requestID)
		return errors.WithMessage(err, util.ErrWrapMessage())
	})

	if err := g.Wait(); err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}

	var result []response.ReportSymbolGetList
	symbols, _ := service.SymbolOriginalService.GetAllMap(false)
	for k, _ := range symbols {
		if symbol != "" && k != symbol {
			continue
		}

		symbolTrades := response.ReportSymbolGetList{
			Date:               startTime.Format(formatter.DateTime) + "~" + endTime.Format(formatter.DateTime),
			Symbol:             k,
			ClosedPL:           symbolClosedPL[k]["closedPL"],
			ClosedProfit:       symbolClosedPL[k]["profit"],
			ClosedSwap:         symbolClosedPL[k]["swaps"],
			ClosedCommission:   symbolClosedPL[k]["commission"],
			FloatingPL:         endTimeSymbolOpeningTrades[k].FloatingPL,
			FloatingProfit:     endTimeSymbolOpeningTrades[k].Profit,
			FloatingSwap:       endTimeSymbolOpeningTrades[k].Swap,
			FloatingCommission: endTimeSymbolOpeningTrades[k].Commission,
			PeriodPL:           endTimeSymbolOpeningTrades[k].FloatingPL.Sub(startTimeSymbolOpeningTrades[k].FloatingPL).Add(symbolClosedPL[k]["closedPL"]),
		}
		result = append(result, symbolTrades)
	}
	return result, nil
}

func (stc *SymbolController) buildSymbolOpeningTrades(groups []string, dateTime time.Time, requestID string) (map[string]symbolOpeningTrades, error) {
	openingTrades, err := service.TradeService.GetPointOpeningTrades(groups, dateTime, "")
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}

	//获取历史时间点持仓盈亏
	floatingPLInfo, err := service.FloatingPLService.GetPointFloatingPLInfo(openingTrades, dateTime, requestID)
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}

	result := map[string]symbolOpeningTrades{}

	symbolLeverage, err := service.SymbolLeverageService.GetAllMap()
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}

	symbols, err := service.SymbolOriginalService.GetAllMap(false)
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}
	for k, _ := range symbols {
		symbolOpeningTrade := symbolOpeningTrades{
			FloatingPL: decimal.Zero,
			Profit:     decimal.Zero,
			Swap:       decimal.Zero,
			Commission: decimal.Zero,
		}

		for _, v := range openingTrades {
			symbol := symbolLeverage[v.Symbol]
			if symbol == k {
				swap := floatingPLInfo[v.Ticket].FloatingSwap
				profit := floatingPLInfo[v.Ticket].FloatingProfit
				floatingPL := floatingPLInfo[v.Ticket].FloatingPL

				symbolOpeningTrade.FloatingPL = symbolOpeningTrade.FloatingPL.Add(floatingPL)
				symbolOpeningTrade.Profit = symbolOpeningTrade.Profit.Add(profit)
				symbolOpeningTrade.Swap = symbolOpeningTrade.Swap.Add(swap)
				symbolOpeningTrade.Commission = symbolOpeningTrade.Commission.Add(v.Commission)
			}
		}
		result[k] = symbolOpeningTrade
	}
	return result, nil
}

func (stc *SymbolController) buildAmountInfo(groups []string, startTime, endTime time.Time) (map[string]decimal.Decimal, error) {
	//时间段内总入金
	amountIn, err := service.CashFlowService.GetCashInflow(groups, startTime, endTime)
	if err != nil {
		return map[string]decimal.Decimal{}, errors.WithMessage(err, util.ErrWrapMessage())
	}
	//时间段内总出金
	amountOut, err := service.CashFlowService.GetCashOutflow(groups, startTime, endTime)
	if err != nil {
		return map[string]decimal.Decimal{}, errors.WithMessage(err, util.ErrWrapMessage())
	}
	//查询结束时间点余额
	pointBalance, err := service.BalanceService.GetPointBalance(groups, endTime)
	if err != nil {
		return map[string]decimal.Decimal{}, errors.WithMessage(err, util.ErrWrapMessage())
	}
	result := map[string]decimal.Decimal{
		"cashInflow":     amountIn,
		"cashOutflow":    amountOut,
		"netCashFlow":    amountIn.Sub(amountOut),
		"endTimeBalance": pointBalance,
	}
	return result, nil
}
