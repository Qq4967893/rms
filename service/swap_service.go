package service

import (
	"backend-rms/repo"
	"backend-rms/util"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type swapService struct {
}

var SwapService = &swapService{}

func (s *swapService) GetPointSwaps(date string) (swapMap map[int64]decimal.Decimal, err error) {
	swapMap, err = repo.SwapRecordeRepo.GetAllByDate(date)
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}
	return
}
