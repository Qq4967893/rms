package service

import (
	"backend-rms/model"
	"backend-rms/repo"
	"backend-rms/util"
	"time"

	"github.com/pkg/errors"
)

type tradingLimitService struct{}

var TradingLimitService = &tradingLimitService{}

func (s *tradingLimitService) GetList(passport int64, tradingLimit int8, paging *model.Paging) ([]*model.TradingLimitList, error) {
	return repo.TradingLimitRepo.GetList(passport, tradingLimit, paging)
}

func (s *tradingLimitService) GetInfoById(id int64) (*model.TradingLimitInfo, error) {
	return repo.TradingLimitRepo.GetInfoById(id)
}

func (s *tradingLimitService) RemoveRestricted(id int64, removeValue int, RemoveTime time.Time, userName string) error {
	data := model.TradingLimit{
		RemoveValue: removeValue,
		RemoveTime:  RemoveTime,
		Operator:    userName,
	}
	err := repo.TradingLimitRepo.Update(id, data)
	if err != nil {
		return errors.WithMessage(err, util.ErrWrapMessage())
	}
	return nil
}
