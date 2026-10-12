package controller

import (
	"backend-rms/http/request"
	"backend-rms/http/response"
	"backend-rms/model"
	"backend-rms/service"
	"backend-rms/util"
	"backend-rms/util/formatter"
	"context"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type AccountHedgingController struct{}

func (ac *AccountHedgingController) GetPreview(c *gin.Context) {
	req := request.GetAccountHedgingPreviewRequest{}
	if ok := request.Validate(c, &req); !ok {
		return
	}

	accountInfo, err := service.AccountService.GetInfoByLogin(req.Login)
	if err != nil {
		response.Error(c, err)
		return
	}
	if accountInfo == nil {
		response.FailWithMessage(c, "Login does not exist")
		return
	}

	groupPairs, err := service.GroupPairService.GetToGroups()
	if err != nil {
		response.Error(c, err)
		return
	}
	//判断准备转组的组名是否正确
	if groupPairs[accountInfo.Group] != req.ToGroup {
		response.FailWithMessage(c, "To group "+req.ToGroup+" is error")
		return
	}

	openingTrades, err := service.TradeService.GetOpeningTrades(nil, req.Login)
	if err != nil {
		response.Error(c, err)
		return
	}

	symbols := make([]string, 0)
	symbolNets := make(map[string]int64)
	for _, openingTrade := range openingTrades {
		symbolOrigin, err := service.SymbolLeverageService.GetSymbolOrigin(openingTrade.Symbol)
		if err != nil {
			response.Error(c, err)
			return
		}

		if !util.InSlice(symbolOrigin, symbols) {
			symbols = append(symbols, symbolOrigin)
		}

		volume := openingTrade.Volume
		if openingTrade.Cmd == 1 {
			volume = openingTrade.Volume * -1
		}
		if _, ok := symbolNets[symbolOrigin]; ok {
			symbolNets[symbolOrigin] = symbolNets[symbolOrigin] + volume
		} else {
			symbolNets[symbolOrigin] = volume
		}
	}

	//检测账号持仓单是否都是可交易品种
	for _, symbol := range symbols {
		isPossibleTrade, err := service.SymbolService.IsPossibleTrade(symbol)
		if err != nil {
			response.Error(c, err)
			return
		}

		if !isPossibleTrade {
			response.FailWithMessage(c, symbol+" is closing, unable to transfer group")
			return
		}
	}

	requestID := c.GetString("request_id")
	//获取交易账号对冲的预览数据
	result, err := ac.buildPreview(req.Login, accountInfo.Group, symbolNets, requestID)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.SucceedWithData(c, "", result)
}

func (ac *AccountHedgingController) GetConfirm(c *gin.Context) {
	req := request.GetAccountHedgingConfirmRequest{}
	if ok := request.Validate(c, &req); !ok {
		return
	}

	accountInfo, err := service.AccountService.GetInfoByLogin(req.Login)
	if err != nil {
		response.Error(c, err)
		return
	}
	if accountInfo == nil {
		response.FailWithMessage(c, "Login does not exist")
		return
	}

	groupPairs, err := service.GroupPairService.GetToGroups()
	if err != nil {
		response.Error(c, err)
		return
	}
	//判断准备转组的组名是否正确
	if groupPairs[accountInfo.Group] != req.ToGroup {
		response.FailWithMessage(c, "To group "+req.ToGroup+" is error")
		return
	}

	openingTrades, err := service.TradeService.GetOpeningTrades(nil, req.Login)
	if err != nil {
		response.Error(c, err)
		return
	}

	book, err := service.GroupService.GetBookByGroup(accountInfo.Group)
	toBook, err := service.GroupService.GetBookByGroup(req.ToGroup)

	exposures, err := ac.getConfirmExposures(openingTrades, book, toBook)
	if err != nil {
		response.Error(c, err)
		return
	}

	summaries, err := ac.getConfirmSummaries(openingTrades, book, toBook)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.SucceedWithData(c, "", response.Confirm{Summaries: summaries, Exposures: exposures})
}

func (ac *AccountHedgingController) Hedge(c *gin.Context) {
	req := request.HedgeAccountHedgingRequest{}
	if ok := request.Validate(c, &req); !ok {
		return
	}

	logins, err := service.AccountService.GetInfoByLogin(req.Login)
	if err != nil {
		response.Error(c, err)
		return
	}
	if logins == nil {
		response.FailWithMessage(c, "Login is not found")
		return
	}
	if logins.Group != req.Group {
		response.FailWithMessage(c, "Group is not login's")
		return
	}

	//判断关联关系是否正确
	groups, err := service.GroupPairService.GetToGroups()
	if err != nil {
		response.Error(c, err)
		return
	}
	if _, ok := groups[req.Group]; !ok || groups[req.Group] != req.ToGroup {
		response.FailWithMessage(c, "Parameter group or to group is error")
		return
	}

	//获取账号转组对冲所需的对冲账号
	hedgeAccount, err := service.HedgeAccountService.GetInfo()
	if err != nil {
		response.Error(c, err)
		return
	}
	if hedgeAccount == nil {
		response.FailWithMessage(c, "Hedge account is not found")
		return
	}

	//判断当前账号角色是否可以操作该对冲账号
	managerId := c.GetInt64("login_manager_id")
	manager, err := service.ManagerService.GetInfoById(managerId)
	if err != nil {
		response.Error(c, err)
		return
	}
	roleIds := strings.Split(hedgeAccount.RoleIds, ",")
	if !util.InSlice(strconv.FormatInt(manager.RoleId, 10), roleIds) {
		response.FailWithMessage(c, "No permissions")
		return
	}

	_, err = service.AccountHedgingService.ChangeGroupHedging(req.Login, req.Group, req.ToGroup, hedgeAccount.HedgeAccount, manager.ID)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Succeed(c)
}

func (ac *AccountHedgingController) buildPreview(login int64, group string, symbolNets map[string]int64, requestID string) (response.Preview, error) {
	book, err := service.GroupService.GetBookByGroup(group)
	if err != nil {
		return response.Preview{}, errors.WithMessage(err, util.ErrWrapMessage())
	}

	logins := make([]int64, 0)
	logins = append(logins, login)
	equityMap, err := service.AccountService.GetEquityMap(context.Background(), requestID, logins)
	if err != nil {
		return response.Preview{}, errors.WithMessage(err, util.ErrWrapMessage())
	}

	closePLs, err := service.TradeService.GetAccountClosedPLWithoutGroup(logins)
	if err != nil {
		return response.Preview{}, errors.WithMessage(err, util.ErrWrapMessage())
	}
	var closePL decimal.Decimal
	if len(closePLs) > 0 {
		closePL = closePLs[0].ClosedPL
	} else {
		closePL = decimal.Zero
	}

	grossPL := equityMap[login].Equity.Sub(equityMap[login].Balance).Add(closePL)
	net := make(map[string]string)
	for k, v := range symbolNets {
		net[k] = decimal.NewFromFloat(float64(v) / 100).StringFixed(2)
	}

	return response.Preview{
		Login:   login,
		Group:   group,
		Book:    book,
		Balance: formatter.Amount(equityMap[login].Balance),
		Equity:  formatter.Amount(equityMap[login].Equity),
		GrossPL: formatter.Amount(grossPL),
		Net:     net,
	}, nil
}

func (ac *AccountHedgingController) getConfirmSummaries(openingTrades []model.OpeningTrade, book, toBook string) (map[string][]response.ConfirmSummary, error) {
	books := []string{
		strings.ToLower(book),
	}
	toBooks := []string{
		strings.ToLower(toBook),
	}
	bookSummaries, err := service.CacheService.GetBookSummaries(books)
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}
	toBookSummaries, err := service.CacheService.GetBookSummaries(toBooks)
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}

	//获取账号的持仓单品种的头寸
	symbolNets := make(map[string]decimal.Decimal)
	for _, trade := range openingTrades {
		symbol, err := service.SymbolLeverageService.GetSymbolOrigin(trade.Symbol)
		if err != nil {
			return nil, errors.WithMessage(err, util.ErrWrapMessage())
		}
		volume := decimal.NewFromInt(trade.Volume)
		if trade.Cmd == 1 {
			volume = volume.Mul(decimal.NewFromInt(-1))
		}

		if _, ok := symbolNets[symbol]; ok {
			symbolNets[symbol] = symbolNets[symbol].Add(volume)
		} else {
			symbolNets[symbol] = volume
		}
	}

	confirmSummaries := map[string][]response.ConfirmSummary{}
	for k, v := range symbolNets {
		var bookAfterNet decimal.Decimal
		var toBookAfterNet decimal.Decimal
		//判断book转toBook是否需要反向下单（A转B反向下单，B转A正向下单）
		if book == "A" && toBook == "B" {
			bookAfterNet = bookSummaries[k].Et6NetLots.Sub(v)
			toBookAfterNet = toBookSummaries[k].Et6NetLots.Add(v)
		} else if book == "B" && toBook == "A" {
			bookAfterNet = bookSummaries[k].Et6NetLots.Sub(v)
			toBookAfterNet = toBookSummaries[k].Et6NetLots.Add(v)
		}

		confirmSummaries[book] = append(confirmSummaries[book], response.ConfirmSummary{
			Symbol:    k,
			BeforeNet: formatter.Amount(bookSummaries[k].Et6NetLots.Div(decimal.NewFromInt(100))),
			AfterNet:  formatter.Amount(bookAfterNet.Div(decimal.NewFromInt(100))),
		})

		confirmSummaries[toBook] = append(confirmSummaries[toBook], response.ConfirmSummary{
			Symbol:    k,
			BeforeNet: formatter.Amount(toBookSummaries[k].Et6NetLots.Div(decimal.NewFromInt(100))),
			AfterNet:  formatter.Amount(toBookAfterNet.Div(decimal.NewFromInt(100))),
		})
	}
	return confirmSummaries, nil
}

func (ac *AccountHedgingController) getConfirmExposures(openingTrades []model.OpeningTrade, book, toBook string) (map[string][]response.ConfirmExposure, error) {
	books := []string{
		strings.ToLower(book),
	}
	toBooks := []string{
		strings.ToLower(toBook),
	}
	bookIdentifierExposures, err := service.ExposureService.GetIdentifierExposureMap(books)
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}
	toBookIdentifierExposures, err := service.ExposureService.GetIdentifierExposureMap(toBooks)
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}

	//获取账号的持仓单品种的头寸
	symbolNets := make(map[string]decimal.Decimal)
	for _, trade := range openingTrades {
		symbol, err := service.SymbolLeverageService.GetSymbolOrigin(trade.Symbol)
		if err != nil {
			return nil, errors.WithMessage(err, util.ErrWrapMessage())
		}
		volume := decimal.NewFromInt(trade.Volume)
		if trade.Cmd == 1 {
			volume = volume.Mul(decimal.NewFromInt(-1))
		}

		if _, ok := symbolNets[symbol]; ok {
			symbolNets[symbol] = symbolNets[symbol].Add(volume)
		} else {
			symbolNets[symbol] = volume
		}
	}

	currencyExposureForABook := make(map[string]map[string]decimal.Decimal)
	currencyExposureForBBook := make(map[string]map[string]decimal.Decimal)
	for k, v := range symbolNets {
		symbolIdentifier, err := service.SymbolTypeService.GetIdentifierBySymbol(k)
		if err != nil {
			return nil, errors.WithMessage(err, util.ErrWrapMessage())
		}

		//redis内存储的exposure
		//特殊情况：例：如果品种的表示 = energy或 = metal需要转换成commodities（原因：因为system_summary_rms系统存储的是commodities）
		redisExposureIdentifier := model.SymbolIdentifierRedisExposure[symbolIdentifier]

		bookExposure := make(map[string]service.Exposure)
		toBookExposure := make(map[string]service.Exposure)

		currencyExposures, err := service.ExposureService.GetCurrencyExposureBySymbol(k, v)
		if err != nil {
			return nil, errors.WithMessage(err, util.ErrWrapMessage())
		}
		for currency, exposure := range currencyExposures {
			//如果currency == USD，直接取redis内品种标识currency的USD(原因：如果是XAUUSD品种会分为XAU,USD，XAUUSD的品种标识属于贵金属，但是redis标识贵金属内没有USD的值)
			if currency == "USD" {
				bookExposure = bookIdentifierExposures["currency"]
				toBookExposure = toBookIdentifierExposures["currency"]
			} else {
				bookExposure = bookIdentifierExposures[redisExposureIdentifier]
				toBookExposure = toBookIdentifierExposures[redisExposureIdentifier]
			}

			//获取exposure转换成exposure_usd的值
			exposureUsd, err := service.ExposureService.GetCurrencyExposureUSD(symbolIdentifier, currency, exposure, k, bookExposure[currency].Symbol)
			if err != nil {
				return nil, errors.WithMessage(err, util.ErrWrapMessage())
			}

			if currencyExposureForABook[currency] == nil {
				currencyExposureForABook[currency] = make(map[string]decimal.Decimal)
			}
			if currencyExposureForBBook[currency] == nil {
				currencyExposureForBBook[currency] = make(map[string]decimal.Decimal)
			}
			//判断book转toBook（A转B：A的exposure减少，B的exposure增加, B转A：A的exposure增加，B的exposure减少）
			if book == "A" && toBook == "B" {
				//操作A组
				//redis内取出的初始值
				currencyExposureForABook[currency]["beforeNetUsd"] = bookExposure[currency].ExposureUSD
				//判断currency是否存在
				if _, ok := currencyExposureForABook[currency]["afterNetUsd"]; ok {
					currencyExposureForABook[currency]["afterNetUsd"] = currencyExposureForABook[currency]["afterNetUsd"].Sub(exposureUsd)
				} else {
					//redis内取出的初始值+计算后的值
					currencyExposureForABook[currency]["afterNetUsd"] = currencyExposureForABook[currency]["beforeNetUsd"].Sub(exposureUsd)
				}

				//操作B组
				//redis内取出的初始值
				currencyExposureForBBook[currency]["beforeNetUsd"] = toBookExposure[currency].ExposureUSD
				//判断currency是否存在
				if _, ok := currencyExposureForBBook[currency]["afterNetUsd"]; ok {
					currencyExposureForBBook[currency]["afterNetUsd"] = currencyExposureForBBook[currency]["afterNetUsd"].Add(exposureUsd)
				} else {
					//redis内取出的初始值+计算后的值
					currencyExposureForBBook[currency]["afterNetUsd"] = currencyExposureForBBook[currency]["beforeNetUsd"].Add(exposureUsd)
				}
			} else if book == "B" && toBook == "A" {
				//操作A组
				//redis内取出的初始值
				currencyExposureForABook[currency]["beforeNetUsd"] = toBookExposure[currency].ExposureUSD
				//判断currency是否存在
				if _, ok := currencyExposureForABook[currency]["afterNetUsd"]; ok {
					currencyExposureForABook[currency]["afterNetUsd"] = currencyExposureForABook[currency]["afterNetUsd"].Add(exposureUsd)
				} else {
					//redis内取出的初始值+计算后的值
					currencyExposureForABook[currency]["afterNetUsd"] = currencyExposureForABook[currency]["beforeNetUsd"].Add(exposureUsd)
				}

				//操作B组
				//redis内取出的初始值
				currencyExposureForBBook[currency]["beforeNetUsd"] = bookExposure[currency].ExposureUSD
				//判断currency是否存在
				if _, ok := currencyExposureForBBook[currency]["afterNetUsd"]; ok {
					currencyExposureForBBook[currency]["afterNetUsd"] = currencyExposureForBBook[currency]["afterNetUsd"].Sub(exposureUsd)
				} else {
					//redis内取出的初始值+计算后的值
					currencyExposureForBBook[currency]["afterNetUsd"] = currencyExposureForBBook[currency]["beforeNetUsd"].Sub(exposureUsd)
				}
			}
		}
	}
	confirmExposures := make(map[string][]response.ConfirmExposure)
	sliceConfirmExposureForABook := make([]response.ConfirmExposure, 0)
	for currency, exposureInfo := range currencyExposureForABook {
		sliceConfirmExposureForABook = append(sliceConfirmExposureForABook, response.ConfirmExposure{
			Currency:     currency,
			BeforeNetUSD: formatter.Amount(exposureInfo["beforeNetUsd"]),
			AfterNetUSD:  formatter.Amount(exposureInfo["afterNetUsd"]),
		})
	}
	confirmExposures["A"] = sliceConfirmExposureForABook

	sliceConfirmExposureForBBook := make([]response.ConfirmExposure, 0)
	for currency, exposureInfo := range currencyExposureForBBook {
		sliceConfirmExposureForBBook = append(sliceConfirmExposureForBBook, response.ConfirmExposure{
			Currency:     currency,
			BeforeNetUSD: formatter.Amount(exposureInfo["beforeNetUsd"]),
			AfterNetUSD:  formatter.Amount(exposureInfo["afterNetUsd"]),
		})
	}
	confirmExposures["B"] = sliceConfirmExposureForBBook
	return confirmExposures, nil
}
