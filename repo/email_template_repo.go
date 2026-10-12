package repo

import (
	"backend-rms/core/database"
	"backend-rms/model"
	"backend-rms/util"

	"github.com/pkg/errors"
)

type emailTemplateRepo struct{}

var EmailTemplateRepo = &emailTemplateRepo{}

func (r *emailTemplateRepo) GetInfosByEvent(event string, userType int, worldCode string) (emailTemplate []model.EmailTemplate, err error) {
	tx := database.APIDB.Where("event", event).
		Where("user_type IN ?", []int{0, userType})
	if worldCode != "" {
		tx.Where("only_world = ? or only_world like ?", "", "%"+worldCode+"%")
	}
	result := tx.Find(&emailTemplate)
	return emailTemplate, errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
}
