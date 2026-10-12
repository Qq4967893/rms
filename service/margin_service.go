package service

import (
	"backend-rms/grpc/tsp"
	"context"

	"github.com/shopspring/decimal"
)

type marginService struct{}

var MarginService = &marginService{}

func (s *marginService) GetMarginLevelBatch(requestID string, logins []int64) (map[int64]decimal.Decimal, error) {
	userEquities, err := tsp.GetUserEquity(context.Background(), requestID, logins)
	if err != nil {
		return nil, err
	}

	marginLevels := make(map[int64]decimal.Decimal, 0)
	for _, v := range userEquities {
		var marginLevel decimal.Decimal
		if v.MarginLevel != 0 {
			marginLevel = decimal.NewFromFloat(v.MarginLevel)
		}
		marginLevels[int64(v.Login)] = marginLevel
	}
	return marginLevels, nil
}
