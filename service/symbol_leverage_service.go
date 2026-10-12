package service

import (
	"backend-rms/repo"
	"backend-rms/util"

	"github.com/pkg/errors"
)

type symbolLeverageService struct{}

var SymbolLeverageService = &symbolLeverageService{}

func (s *symbolLeverageService) GetAllMap() (map[string]string, error) {
	symbols, err := repo.SymbolLeverageRepo.GetAll()
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}

	var symbolMap = make(map[string]string)
	for _, symbol := range symbols {
		symbolMap[symbol.Symbol] = symbol.Original
	}
	return symbolMap, nil
}

func (s *symbolLeverageService) GetSymbolOrigin(symbol string) (symbolOrigin string, err error) {
	symbolInfo, err := repo.SymbolLeverageRepo.GetInfoBySymbol(symbol)
	if err != nil {
		return "", errors.WithMessage(err, util.ErrWrapMessage())
	}

	if symbolInfo == nil {
		return "", errors.Errorf("can't find the symbol(%v) in symbol_leverage table", symbol)
	}

	return symbolInfo.Original, nil
}
