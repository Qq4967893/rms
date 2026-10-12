package repo

import (
	"backend-rms/core/database"
	"backend-rms/model"
	"backend-rms/util"

	"github.com/pkg/errors"
)

type symbolRepo struct {
}

var SymbolRepo = &symbolRepo{}

func (r *symbolRepo) GetAll() (symbols []*model.Symbol, err error) {
	result := database.ET6ConfigDB.Table("symbol").Find(&symbols)
	return symbols, errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
}
