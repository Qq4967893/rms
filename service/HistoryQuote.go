package service

import (
	"backend-rms/core/database"
	"backend-rms/util"
	"strings"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type historyQuoteService struct {
}

type SymbolHistoryQuote struct {
	Close decimal.Decimal
}

var HistoryQuoteService = &historyQuoteService{}

//start参数要注意，tb的单位是毫秒
func (hqs *historyQuoteService) GetLastClosePrice(symbol string, start int64) (decimal.Decimal, error) {
	return hqs.GetClosingPrice(symbol, start, 1440)
}

func (hqs *historyQuoteService) GetClosingPrice(symbol string, start int64, priceTimeType int) (closingPrice decimal.Decimal, err error) {
	table := hqs.getTable(symbol)
	result := database.CHARTDB.Table(table).Where("ctm < ?", start).
		Where("type", priceTimeType).Select("close").Order("ctm desc").Limit(1).Scan(&closingPrice)
	if result.Error != nil {
		return decimal.Zero, errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
	}
	return closingPrice, nil
}

func (hqs *historyQuoteService) getTable(symbol string) string {
	table := strings.Replace(symbol, ".", "_", -1)
	return strings.ToLower(table)
}
