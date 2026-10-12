package service

import (
	"backend-rms/model"
	"backend-rms/repo"
)

type copyPayService struct{}

var CopyPayService = &copyPayService{}

func (s *copyPayService) GetAll(followerId, masterId int64) (copyPays []model.CopyPay, err error) {
	return repo.CopyPayRepo.GetAll(followerId, masterId)
}
