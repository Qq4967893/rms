package service

import (
	"backend-rms/repo"
	"backend-rms/util"

	"github.com/pkg/errors"
)

type eventService struct{}

var EventService = &eventService{}

func (s *eventService) GetOpenTypeMapByTickets(tickets []int64) (map[int64]int8, error) {
	events, err := repo.EventRepo.GetEventListByTickets(tickets)
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}

	recordTickets := make(map[int64]struct{})
	pendingTickets := make(map[int64]struct{})
	for _, evt := range events {
		recordTickets[evt.Ticket] = struct{}{}
		if evt.Action == 2 && (evt.Type == 10 || evt.Type == 11 || evt.Type == 12 || evt.Type == 13) {
			pendingTickets[evt.Ticket] = struct{}{}
		}
	}

	openTypeMap := make(map[int64]int8, len(tickets))
	for _, ticket := range tickets {
		if _, ok := recordTickets[ticket]; ok {
			if _, ok2 := pendingTickets[ticket]; ok2 {
				openTypeMap[ticket] = 2
			} else {
				openTypeMap[ticket] = 1
			}
		} else {
			openTypeMap[ticket] = 0
		}
	}

	return openTypeMap, nil
}
