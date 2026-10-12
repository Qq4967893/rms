package repo

import (
	"backend-rms/core/database"
	"backend-rms/model"
	"backend-rms/util"

	"github.com/pkg/errors"
)

type sourceRepo struct {
}

var SourceRepo = &sourceRepo{}

func (r *sourceRepo) GetAll(sourceTypes []int) (sources []*model.Source, err error) {
	query := database.ET6ConfigDB.Table("source").Where("status != ?", 3)
	if len(sourceTypes) > 0 {
		query.Where("source_type IN ?", sourceTypes)
	}
	result := query.Find(&sources)
	return sources, errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
}
