package repo

import (
	"backend-rms/core/database"
	"backend-rms/model"
	"backend-rms/util"

	"github.com/pkg/errors"
)

type symbolOriginalRepo struct{}

var SymbolOriginalRepo = &symbolOriginalRepo{}

func (r *symbolOriginalRepo) GetAll(isValid bool) (list []*model.SymbolOriginal, err error) {
	tx := database.APIDB.Table("symbols as o")
	if isValid {
		tx.Where("o.status = ?", 1)
	}
	result := tx.Select("DISTINCT o.symbol_origin as original, o.digits").Find(&list)
	return list, errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
}

func (r *symbolOriginalRepo) GetDigitsByOriginal(original string) (digits int64, err error) {
	err = database.APIDB.Table("symbols").
		Select("digits").
		Where("symbol_origin", original).
		Scan(&digits).Error
	return digits, errors.Wrapf(err, "%s failed", util.GetCurrentFuncName())
}
