package service

import (
	"backend-rms/model"
	"backend-rms/repo"
	"backend-rms/util"

	"github.com/pkg/errors"
)

type hedgeAccountService struct {
}

var HedgeAccountService = &hedgeAccountService{}

func (hs *hedgeAccountService) GetInfo() (*model.HedgeAccount, error) {
	hedgeAccount, err := repo.HedgeAccountRepo.GetInfo()
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}
	return hedgeAccount, nil
}

func (hs *hedgeAccountService) GetList() ([]*model.HedgeAccount, error) {
	return repo.HedgeAccountRepo.GetList()
}

func (hs *hedgeAccountService) Update(id int64, params model.HedgeAccount) error {
	return repo.HedgeAccountRepo.Update(id, params)
}
