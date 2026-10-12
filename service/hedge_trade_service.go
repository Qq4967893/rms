package service

import (
	"backend-rms/model"
	"backend-rms/repo"
)

type hedgeTradeService struct{}

var HedgeTradeService = &hedgeTradeService{}

func (s *hedgeTradeService) GetList(login, managerId int64, paging *model.Paging) ([]*model.HedgeTrade, error) {
	return repo.HedgeTradeRepo.GetAll(login, managerId, paging)
}

func (s *hedgeTradeService) GetCount(login, managerId int64) (int64, error) {
	return repo.HedgeTradeRepo.GetCount(login, managerId)
}

func (s *hedgeTradeService) BatchesAdd(hedgeTrades []model.HedgeTrade) error {
	return repo.HedgeTradeRepo.Add(hedgeTrades)
}
