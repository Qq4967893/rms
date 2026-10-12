package service

import (
	"backend-rms/model"
	"backend-rms/repo"
	"backend-rms/util"

	"github.com/pkg/errors"
)

type symbolTypeService struct{}

var SymbolTypeService = &symbolTypeService{}

func (s *symbolTypeService) GetEnergySymbols(isValid bool) ([]string, error) {
	symbols, err := repo.SymbolTypeRepo.GetSymbols("energy", isValid)
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}
	return symbols, nil
}

func (s *symbolTypeService) GetForexSymbols(isValid bool) ([]string, error) {
	symbols, err := repo.SymbolTypeRepo.GetSymbols("forex", isValid)
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}
	return symbols, nil
}

func (s *symbolTypeService) GetIndexSymbols(isValid bool) ([]string, error) {
	symbols, err := repo.SymbolTypeRepo.GetSymbols("index", isValid)
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}
	return symbols, nil
}

func (s *symbolTypeService) GetMetalSymbols(isValid bool) ([]string, error) {
	symbols, err := repo.SymbolTypeRepo.GetSymbols("metal", isValid)
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}
	return symbols, nil
}

func (s *symbolTypeService) GetStockSymbols(isValid bool) ([]string, error) {
	symbols, err := repo.SymbolTypeRepo.GetSymbols("stock", isValid)
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}
	return symbols, nil
}

// 获取品种的标识
func (s *symbolTypeService) GetIdentifierBySymbol(symbol string) (string, error) {
	symbols, err := repo.SymbolTypeRepo.GetIdentifierByOriginal(symbol)
	if err != nil {
		return "", errors.WithMessage(err, util.ErrWrapMessage())
	}
	if symbols.Source == "" {
		return "", errors.New("Symbol is not found")
	}
	return model.SourceTypeMap[symbols.SourceType], nil
}

// 获取品种的货币
func (s *symbolTypeService) GetSymbolCurrencies(symbol, symbolIdentifier, currency, profitCurrency string) []string {
	currencies := make([]string, 0)
	switch symbolIdentifier {
	case model.SymbolIdentifierCurrency:
		currencies = append(currencies, symbol[:3], profitCurrency)
	case model.SymbolIdentifierMetal:
		currencies = append(currencies, currency, profitCurrency)
	case model.SymbolIdentifierEnergy, model.SymbolIdentifierCrypto:
		currencies = append(currencies, symbol[:3])
	case model.SymbolIdentifierIndices, model.SymbolIdentifierStock:
		currencies = append(currencies, symbol)
	}
	return currencies
}
