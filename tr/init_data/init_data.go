package init_data

import (
	"backend-rms/core/logger"
	"backend-rms/model"
	"backend-rms/service"
)

func InitData() {
	sourceList, err := service.SourceService.GetList()
	if err != nil {
		logger.Error(err.Error())
		return
	}
	sources := make(map[string]*model.Source)
	for _, v := range sourceList {
		sources[v.Source] = v
	}

	symbolList, err := service.SymbolService.GetList()
	if err != nil {
		logger.Error(err.Error())
		return
	}
	symbols := make(map[string]*model.Symbol)
	for _, v := range symbolList {
		symbols[v.Symbol] = v
	}

	groupList, err := service.GroupService.GetAllList()
	if err != nil {
		logger.Error(err.Error())
		return
	}
	groupMaps := make(map[string]*model.Group)
	groups := make([]string, 0)
	for _, v := range groupList {
		groups = append(groups, v.Name)
		groupMaps[v.Name] = v
	}

	accountList, err := service.AccountService.GetList(groups)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	accounts := make(map[int64]*model.Account)
	for _, v := range accountList {
		accounts[v.Login] = v
	}

	tradeList, err := service.TradeService.GetOpeningTrades(groups, 0)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	trades := make(map[int64]*model.Trade)
	for _, v := range tradeList {
		trades[v.Ticket] = &model.Trade{
			Ticket:     v.Ticket,
			Login:      v.Login,
			Symbol:     v.Symbol,
			Cmd:        v.Cmd,
			Volume:     v.Volume,
			OpenTime:   v.OpenTime,
			OpenPrice:  v.OpenPrice,
			Commission: v.Commission,
			Swaps:      v.Swaps,
			Profit:     v.Profit,
		}
	}

	service.InitializedData = service.Data{
		Sources:  sources,
		Symbols:  symbols,
		Groups:   groupMaps,
		Accounts: accounts,
		Trades:   trades,
	}
}

var Books = []string{"A", "B"}

func InitGroupData() {
	currencySymbol := service.SourceService.GetCurrencySymbols()
	for _, v := range Books {
		if _, exists := service.SummaryBook[v]; !exists {
			service.SummaryBook[v] = make(map[string]*service.InitSummary)
		}
		if _, exists := service.ExposureBook[v]; !exists {
			service.ExposureBook[v] = make(map[string]*service.InitExposure)
		}
		for _, soucreV := range service.InitializedData.Sources {
			service.SummaryBook[v][soucreV.Source] = &service.InitSummary{Symbol: soucreV.Source}

			sourceCurrencies := service.SourceService.GetSourceCurrencies(soucreV)
			for _, currencyV := range sourceCurrencies {
				if _, ok := currencySymbol[currencyV]; !ok {
					logger.Errorf("Currency Symbol 内不存在该 %v", currencyV)
					continue
				}

				if _, ok := service.ExposureBook[v][currencyV]; !ok {
					service.ExposureBook[v][currencyV] = &service.InitExposure{
						Currency: currencyV,
						Symbol:   currencySymbol[currencyV],
					}
				}
			}
		}
	}

	for _, v := range service.InitializedData.Trades {
		service.ExposureService.HandleForOrder(v, "open")
		service.SummaryService.HandleForOrder(v, "open")
	}

	// 把初始化的summary更新到redis
	for _, v := range Books {
		service.CacheService.SetSummaryBook(v)
	}
}
