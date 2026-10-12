package service

import (
	"backend-rms/model"
	"backend-rms/repo"
	"backend-rms/util"

	"github.com/pkg/errors"
)

type userService struct{}

var UserService = &userService{}

func (s *userService) GetInfoByPassport(passport int64) (user *model.User, err error) {
	user, err = repo.UserRepo.GetInfoByPassport(passport)
	return user, errors.WithMessage(err, util.ErrWrapMessage())
}
