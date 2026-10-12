package repo

import (
	"backend-rms/core/database"
	"backend-rms/model"
	"backend-rms/util"

	"github.com/pkg/errors"
)

type eventRepo struct{}

var EventRepo = &eventRepo{}

func (r *eventRepo) GetEventListByTickets(tickets []int64) (events []*model.Event, err error) {
	result := database.ET6DB.Where("ticket IN ?", tickets).
		Distinct("ticket, action, type").Find(&events)
	return events, errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
}
