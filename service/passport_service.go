package service

import (
	"backend-rms/util"

	"github.com/pkg/errors"
)

type passportService struct{}

var PassportService = &passportService{}

func (s *passportService) GetAccounts(passport, login int64) (accounts map[string]any, err error) {
	if passport <= 0 {
		passport, err = TradeAccountService.GetPassportByLogin(login)
		if err != nil {
			return nil, errors.WithMessage(err, util.ErrWrapMessage())
		}
	}

	tradeAccountIDs, err := TradeAccountService.GetTradeAccountIDs(passport)
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}

	return map[string]any{
		"passport": passport,
		"accounts": tradeAccountIDs,
	}, nil
}
