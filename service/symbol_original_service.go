package service

import (
	"backend-rms/model"
	"backend-rms/repo"
	"backend-rms/util"

	"github.com/pkg/errors"
	"golang.org/x/exp/slices"
)

type symbolOriginalService struct {
}

var SymbolOriginalService = &symbolOriginalService{}

func (s *symbolOriginalService) GetAllList(isValid bool) (list []*model.SymbolOriginal, err error) {
	list, err = repo.SymbolOriginalRepo.GetAll(isValid)
	return list, errors.WithMessage(err, util.ErrWrapMessage())
}

func (s *symbolOriginalService) GetAllMap(isValid bool) (symbolOriginalMap map[string]*model.SymbolOriginal, err error) {
	list, err := repo.SymbolOriginalRepo.GetAll(isValid)

	symbolOriginalMap = make(map[string]*model.SymbolOriginal)
	for _, symbolOriginal := range list {
		symbolOriginalMap[symbolOriginal.Original] = symbolOriginal
	}
	return symbolOriginalMap, errors.WithMessage(err, util.ErrWrapMessage())
}

func (s *symbolOriginalService) GetDigitsByOriginal(original string) (digits int64, err error) {
	digits, err = repo.SymbolOriginalRepo.GetDigitsByOriginal(original)
	return
}

func (s *symbolOriginalService) ConvertCurrency(currency string) (string, error) {
	if currency == "USD" {
		return currency, nil
	}

	list, err := s.GetAllList(false)
	if err != nil {
		return "", err
	}

	symbolOriginals := make([]string, 0, len(list))
	for _, s := range list {
		symbolOriginals = append(symbolOriginals, s.Original)
	}

	currencyUSD := currency + "USD"
	if slices.Contains(symbolOriginals, currencyUSD) {
		return currencyUSD, nil
	}

	usdCurrency := "USD" + currency
	if slices.Contains(symbolOriginals, usdCurrency) {
		return usdCurrency, nil
	}
	return "", nil
}
