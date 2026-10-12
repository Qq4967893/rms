package repo

import (
	"backend-rms/core/database"
	"backend-rms/model"
	"backend-rms/util"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type symbolTypeRepo struct{}

var SymbolTypeRepo = &symbolTypeRepo{}

func (r *symbolTypeRepo) GetSymbols(identifier string, isValid bool) (symbols []string, err error) {
	var list []*model.SymbolType
	tx := database.APIDB.Table("symbols as t").
		Where("t.symbol_type like ?", "%"+identifier+"%")

	if isValid {
		tx.Where("t.status = ?", 1)
	}

	result := tx.Select("DISTINCT t.symbol_origin as original").Find(&list)

	for _, v := range list {
		symbols = append(symbols, v.Original)
	}
	return symbols, errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
}

// 获取品种的类型
func (r *symbolTypeRepo) GetIdentifierByOriginal(symbol string) (symbolType model.SourceType, err error) {
	err = database.ET6ConfigDB.Table("source").Where("source", symbol).First(&symbolType).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return model.SourceType{}, nil
	}
	return symbolType, errors.Wrapf(err, "%s failed", util.GetCurrentFuncName())
}
