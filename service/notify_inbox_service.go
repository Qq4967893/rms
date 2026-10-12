package service

import (
	"backend-rms/model"
	"backend-rms/repo"
	"backend-rms/util"

	"github.com/pkg/errors"
)

type notifyInboxService struct{}

var NotifyInboxService = &notifyInboxService{}

func (s *notifyInboxService) Add(params model.NotifyInbox) error {
	err := repo.NotifyInboxRepo.Add(params)
	if err != nil {
		return errors.WithMessage(err, util.ErrWrapMessage())
	}
	return nil
}
