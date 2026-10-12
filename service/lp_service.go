package service

import (
	"backend-rms/model"
	"backend-rms/repo"
	"backend-rms/util"

	"github.com/pkg/errors"
)

type lpService struct {
}

var LpService = &lpService{}

func (ls *lpService) GetSwitchLpList(symbol string) (switchLpRecord []*model.SwitchLpRecord, err error) {
	switchLpRecord, err = repo.SwitchLpRecordRepo.GetAll(symbol)
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}
	return switchLpRecord, nil
}
