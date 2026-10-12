package repo

import (
	"backend-rms/core/database"
	"backend-rms/model"
	"backend-rms/util"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type symbolLeverageRepo struct{}

var SymbolLeverageRepo = &symbolLeverageRepo{}

func (r *symbolLeverageRepo) GetAll() (symbols []*model.SymbolLeverages, err error) {
	result := database.APIDB.Table("symbols").
		Select("symbol, symbol_origin as original").Find(&symbols)
	return symbols, errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
}

func (r *symbolLeverageRepo) GetInfoBySymbol(symbol string) (symbols *model.SymbolLeverages, err error) {
	err = database.APIDB.Table("symbols").
		Select("symbol_origin as original, leverage").
		Where("symbol", symbol).First(&symbols).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return symbols, errors.Wrapf(err, "%s failed", util.GetCurrentFuncName())
}
