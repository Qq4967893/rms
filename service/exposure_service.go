package service

import (
	"backend-rms/core/logger"
	"backend-rms/grpc/tsp"
	"backend-rms/model"
	"backend-rms/util"
	"context"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type exposureService struct{}

var ExposureService = &exposureService{}

// 获取品种的币种Exposure值
func (s *exposureService) GetCurrencyExposureBySymbol(symbol string, netVolume decimal.Decimal) (map[string]decimal.Decimal, error) {
	symbolIdentifier, err := SymbolTypeService.GetIdentifierBySymbol(symbol)
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}

	symbolInfos, err := tsp.GetSymbolInfos(context.Background(), "test request id", symbol, true)
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}
	currencies := SymbolTypeService.GetSymbolCurrencies(symbol, symbolIdentifier, symbolInfos[0].Currency, symbolInfos[0].ProfitCurrency)
	if len(currencies) == 0 {
		return nil, errors.New("Get symbol currencies is failed")
	}

	currencyExposures := make(map[string]decimal.Decimal)
	//头寸volume / 100 * 合约量
	currencyExposures[currencies[0]] = netVolume.Mul(decimal.NewFromFloat(symbolInfos[0].ContractSize)).Div(decimal.NewFromInt(100))
	if len(currencies) > 1 {
		quotes, err := QuoteService.GetQuote(symbol)
		if err != nil {
			return nil, errors.WithMessage(err, util.ErrWrapMessage())
		}
		//如果netVolume是正的证明下的是买单，如果netVolume是负的证明下的是卖单
		var quote decimal.Decimal
		if netVolume.IsNegative() {
			quote = quotes.Bid
		} else {
			quote = quotes.Ask
		}
		//品种开仓价 * 合约量 * 头寸手数 * -1
		currencyExposures[currencies[1]] = quote.Mul(currencyExposures[currencies[0]]).Mul(decimal.NewFromInt(-1))
	}
	return currencyExposures, nil
}

// 获取币种的Exposure USD
func (s *exposureService) GetCurrencyExposureUSD(identifier, currency string, exposure decimal.Decimal, symbol, exposureSymbol string) (decimal.Decimal, error) {
	//USD的 exposure_USD = exposure
	if currency == "USD" {
		return exposure, nil
	}

	//获取转换exposure USD所需的品种报价
	quotes, err := QuoteService.GetQuote(symbol)
	if err != nil {
		return decimal.Zero, errors.WithMessage(err, util.ErrWrapMessage())
	}
	quote := quotes.Ask.Add(quotes.Bid).Div(decimal.NewFromInt(2))

	//非USD的 exposure USD需要exposure 通过品种报价做转换
	//如果是股票，虚拟货币，贵金属，能源，直接乘以自身的报价
	if util.InSlice(identifier, []string{model.SymbolIdentifierStock, model.SymbolIdentifierCrypto, model.SymbolIdentifierMetal, model.SymbolIdentifierEnergy}) {
		return exposure.Mul(quote), nil
	}

	//如果是股指，先乘以自身的报价，如果需要再转换为exposure USD（例：HK50乘以自身报价是HKD的价值，再去乘以USDHKD就是USD的价值），如果不需要直接返回（例子：USA50）
	//货币对
	//判断转换USD的品种报价名称，如果前三位是USD除以，如果是后三位乘以
	if identifier == model.SymbolIdentifierIndices {
		if exposureSymbol == "" {
			return exposure.Mul(quote), nil
		}

		exposureSymbolQuotes, err := QuoteService.GetQuote(exposureSymbol)
		if err != nil {
			return decimal.Zero, errors.WithMessage(err, util.ErrWrapMessage())
		}
		exposureSymbolQuote := exposureSymbolQuotes.Ask.Add(exposureSymbolQuotes.Bid).Div(decimal.NewFromInt(2))
		if exposureSymbol[0:3] == "USD" {
			return exposure.Mul(quote).Div(exposureSymbolQuote), nil
		} else {
			return exposure.Mul(quote).Mul(exposureSymbolQuote), nil
		}
	}

	if identifier == model.SymbolIdentifierCurrency {
		//只有USD的exposureSymbol=""
		if exposureSymbol == "" {
			return exposure, nil
		}

		exposureSymbolQuotes, err := QuoteService.GetQuote(exposureSymbol)
		if err != nil {
			return decimal.Zero, errors.WithMessage(err, util.ErrWrapMessage())
		}
		exposureSymbolQuote := exposureSymbolQuotes.Ask.Add(exposureSymbolQuotes.Bid).Div(decimal.NewFromInt(2))
		if exposureSymbol[0:3] == "USD" {
			return exposure.Div(exposureSymbolQuote), nil
		} else {
			return exposure.Mul(exposureSymbolQuote), nil
		}
	}
	return decimal.Zero, errors.New("The symbol is not in symbol identifier")
}

// 获取币种类型的Exposure Map
func (s *exposureService) GetIdentifierExposureMap(books []string) (map[string]map[string]Exposure, error) {
	bookExposures, err := CacheService.GetBookExposures(books)
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}

	bookExposureMap := make(map[string]map[string]Exposure)
	for identifier, _ := range bookExposures {
		exposureMap := make(map[string]Exposure)
		for _, exposure := range bookExposures[identifier] {
			exposureMap[exposure.Currency] = Exposure{
				Currency:    exposure.Currency,
				Exposure:    exposure.Exposure,
				ExposureUSD: exposure.ExposureUSD,
				Symbol:      exposure.Symbol,
			}
		}
		bookExposureMap[identifier] = exposureMap
	}
	return bookExposureMap, nil
}

// 新增订单对exposure的处理
func (s *exposureService) HandleForOrder(trade *model.Trade, orderType string) {
	var source string
	if symbolInfo, ok := InitializedData.Symbols[trade.Symbol]; ok {
		source = symbolInfo.Source
	} else {
		// 新增订单品种不在初始化的品种内
		logger.Errorf("The order symbol %v is not within the initialized symbol", trade.Symbol)
		return
	}

	sourceInfo := InitializedData.Sources[source]
	sourceCurrencies := SourceService.GetSourceCurrencies(sourceInfo)
	if len(sourceCurrencies) == 0 {
		logger.Errorf("get source currencies error, found 0 currency; source:%v", source)
		return
	}

	book := InitializedData.Groups[InitializedData.Accounts[trade.Login].Group].Book
	sourceCurrenciesOne := sourceCurrencies[0]
	if _, ok := ExposureBook[book][sourceCurrenciesOne]; !ok {
		logger.Errorf("source currencies one not found in exposure map, source:%v, currency:%v", source, sourceCurrenciesOne)
		return
	}

	exposure := decimal.NewFromInt(trade.Volume).Mul(sourceInfo.ContractSize)
	switch trade.Cmd {
	case 0:
		exposure = exposure.Div(decimal.NewFromInt(100))
		break
	case 1:
		exposure = exposure.Div(decimal.NewFromInt(-100))
		break
	}

	//判断是开仓订单还是平仓订单
	if orderType == "open" {
		ExposureBook[book][sourceCurrenciesOne].Exposure = ExposureBook[book][sourceCurrenciesOne].Exposure.Add(exposure)
	} else {
		ExposureBook[book][sourceCurrenciesOne].Exposure = ExposureBook[book][sourceCurrenciesOne].Exposure.Sub(exposure)
	}

	if len(sourceCurrencies) > 1 {
		sourceCurrenciesTwo := sourceCurrencies[1]
		if _, ok := ExposureBook[book][sourceCurrenciesTwo]; !ok {
			logger.Errorf("source currencies two not found in exposure map, source:%v, currency:%v", source, sourceCurrenciesTwo)
			return
		}
		exposureTwo := trade.OpenPrice.Mul(exposure).Mul(decimal.NewFromInt(-1))
		if orderType == "open" {
			ExposureBook[book][sourceCurrenciesTwo].Exposure = ExposureBook[book][sourceCurrenciesTwo].Exposure.Add(exposureTwo)
		} else {
			ExposureBook[book][sourceCurrenciesTwo].Exposure = ExposureBook[book][sourceCurrenciesTwo].Exposure.Sub(exposureTwo)
		}
	}
}
