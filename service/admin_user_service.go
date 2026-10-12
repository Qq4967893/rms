package service

import (
	"backend-rms/repo"
	"backend-rms/util"

	"github.com/pkg/errors"
)

type adminUserService struct{}

var AdminUserService = &adminUserService{}

func (s *adminUserService) GetUsernameByID(id int64) (username string, err error) {
	username, err = repo.AdminUserRepo.GetUsernameByID(id)
	return username, errors.WithMessage(err, util.ErrWrapMessage())
}
