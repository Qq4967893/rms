package service

import (
	"backend-rms/model"
	"backend-rms/repo"
	"backend-rms/util"

	"github.com/pkg/errors"
)

type sourceService struct {
}

var SourceService = &sourceService{}

func (s *sourceService) GetList() (sources []*model.Source, err error) {
	sources, err = repo.SourceRepo.GetAll([]int{})
	return sources, errors.WithMessage(err, util.ErrWrapMessage())
}

func (s *sourceService) GetSourceCurrencies(source *model.Source) []string {
	currencies := make([]string, 0)
	switch source.SourceType {
	case model.SourceFX:
		currencies = append(currencies, source.MarginCurrency, source.ProfitCurrency)
	case model.SourceGold, model.SourceSliver:
		currencies = append(currencies, source.Currency, source.ProfitCurrency)
	case model.SourceGas, model.SourceOil, model.SourceCrypto:
		currencies = append(currencies, source.Source[:3])
	case model.SourceIndex, model.SourceEquityCFD:
		currencies = append(currencies, source.Source)
	}
	return currencies
}

func (s *sourceService) GetCurrencyExposureTypes() map[string]string {
	exposureType := make(map[string]string)
	for outerKey, innerMap := range s.CurrencySourceMap() {
		for innerKey := range innerMap {
			exposureType[innerKey] = outerKey
		}
	}
	return exposureType
}

func (s *sourceService) GetCurrencySymbols() map[string]string {
	currencySymbols := make(map[string]string)
	for _, inner := range s.CurrencySourceMap() {
		for k, v := range inner {
			currencySymbols[k] = v
		}
	}
	return currencySymbols
}

func (s *sourceService) CurrencySourceMap() map[string]map[string]string {
	currencyTypes := map[string][]int{
		"currency":    []int{model.SourceFX},
		"commodities": []int{model.SourceGold, model.SourceSliver, model.SourceGas, model.SourceOil},
		"crypto":      []int{model.SourceCrypto},
		"indices":     []int{model.SourceIndex},
		"stock":       []int{model.SourceEquityCFD},
	}

	result := make(map[string]map[string]string)
	var currencies, commodities, cryptos, indices, stocks = make(map[string]string), make(map[string]string), make(map[string]string), make(map[string]string), make(map[string]string)
	for currencyType, sourceTypes := range currencyTypes {
		sourceInfos, _ := repo.SourceRepo.GetAll(sourceTypes)

		// 处理货币对
		if currencyType == "currency" {
			// USD单独处理添加
			currencies["USD"] = ""
			for _, sourceInfo := range sourceInfos {
				if sourceInfo.Currency != "USD" && sourceInfo.ProfitCurrency != "USD" {
					continue
				}

				var currencyV string
				if sourceInfo.Currency != "USD" {
					currencyV = sourceInfo.Currency
				} else {
					currencyV = sourceInfo.ProfitCurrency
				}
				currencies[currencyV] = sourceInfo.Source
			}
			result[currencyType] = currencies
		}

		if currencyType == "commodities" {
			for _, sourceInfo := range sourceInfos {
				commodities[sourceInfo.Currency] = sourceInfo.Source
			}
			result[currencyType] = commodities
		}

		if currencyType == "crypto" {
			for _, sourceInfo := range sourceInfos {
				// 切割前3位
				cryptos[sourceInfo.Source[0:3]] = sourceInfo.Source
			}
			result[currencyType] = cryptos
		}

		if currencyType == "indices" {
			// 查询股指的盈利货币转换USD的货币对
			for _, sourceInfo := range sourceInfos {
				// 如果是USD赋值为空，如果非USD去查找股指的盈利货币转换USD的货币对
				if sourceInfo.Currency == "USD" {
					indices[sourceInfo.Source] = ""
				} else {
					indices[sourceInfo.Source] = currencies[sourceInfo.Currency]
				}
			}
			result[currencyType] = indices
		}

		if currencyType == "stock" {
			// 查询股指的盈利货币转换USD的货币对
			for _, sourceInfo := range sourceInfos {
				stocks[sourceInfo.Source] = ""
			}
			result[currencyType] = stocks
		}
	}
	return result
}
