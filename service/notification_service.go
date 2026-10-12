package service

import (
	"backend-rms/model"
	"backend-rms/repo"
)

type notificationService struct{}

var NotificationService = &notificationService{}

func (s *notificationService) Add(params model.Notification) (id int64, err error) {
	return repo.NotificationRepo.Add(params)
}
