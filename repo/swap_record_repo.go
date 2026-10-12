package repo

import (
	"backend-rms/core/database"
	"backend-rms/model"
	"backend-rms/util"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type swapRecordeRepo struct {
}

var SwapRecordeRepo = &swapRecordeRepo{}

func (r *swapRecordeRepo) GetAllByDate(date string) (swapMap map[int64]decimal.Decimal, err error) {
	var swaps []*model.SwapRecord
	result := database.APIDB.Table("swap_records").Where("`date` = ?", date).Find(&swaps)
	if result.Error != nil {
		return nil, errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
	}

	swapMap = make(map[int64]decimal.Decimal, len(swaps))
	for _, swap := range swaps {
		swapMap[swap.Ticket] = swap.Swaps
	}

	return swapMap, nil
}
