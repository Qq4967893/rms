package service

import (
	"backend-rms/model"
	"backend-rms/repo"
	"backend-rms/util"
	"math"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type Spread struct {
	Symbol   string
	After    int64
	Before   int64
	DateTime time.Time
}

type spreadService struct{}

var SpreadList []model.SpreadRecord

var SpreadService = &spreadService{}

func (ss spreadService) GetSpreadByDate(symbol string, spreadOffset decimal.Decimal, dateTime time.Time) (decimal.Decimal, error) {
	commonSpread, err := ss.GetCommonSpread(symbol, dateTime)
	if err != nil {
		return decimal.Zero, errors.WithMessage(err, util.ErrWrapMessage())
	}

	spreadByDate := decimal.NewFromInt(commonSpread).Add(spreadOffset)
	// spreadOffset在etc上可以设置为负值，所以注释该判断
	//if spreadByDate.LessThan(decimal.NewFromInt(int64(0))) {
	//	return decimal.Zero, errors.Errorf("%v spread by date is error; datetime is %v", symbol, dateTime)
	//}
	return spreadByDate, nil
}

func (ss spreadService) GetSpreadOffset(cmd int8, openOrClose string, price, transPrice decimal.Decimal, digit int32) decimal.Decimal {
	var spreadOffset decimal.Decimal
	if cmd == 0 {
		if openOrClose == "open" {
			spreadOffset = price.Sub(transPrice)
		} else {
			spreadOffset = transPrice.Sub(price)
		}
	} else {
		if openOrClose == "open" {
			spreadOffset = transPrice.Sub(price)
		} else {
			spreadOffset = price.Sub(transPrice)
		}
	}
	return spreadOffset.Mul(decimal.NewFromFloat(math.Pow10(int(digit)))).Mul(decimal.NewFromInt(2))
}

func (ss spreadService) GetCommonSpread(symbol string, dateTime time.Time) (int64, error) {
	for _, spreadInfo := range SpreadList {
		if dateTime.After(spreadInfo.ChangeTime) || dateTime.Equal(spreadInfo.ChangeTime) {
			if symbol == spreadInfo.Symbol {
				return spreadInfo.AfterSpread, nil
			}
		}
	}

	//获取ETB spread
	etbSpread, err := CacheService.GetETBSpread(symbol)
	if err != nil {
		return 0, errors.WithMessage(err, util.ErrWrapMessage())
	}
	return etbSpread, nil
}

// 原始点差（也就是近14天的平均点差）
func (ss *spreadService) GetRawSpread(symbol string, dateTime time.Time) (int64, error) {
	lp, err := ss.getHistoricalSymbolLp(symbol, dateTime)
	if err != nil {
		return 0, errors.WithMessage(err, util.ErrWrapMessage())
	}

	dateTime.Date()
	avgSpreadOf14Day, err := ss.GetAvgSpreadLast14Days(symbol, lp, dateTime)
	if err != nil {
		return 0, errors.WithMessage(err, util.ErrWrapMessage())
	}
	//if avgSpreadOf14Day == 0 {
	//	return -1, errors.Errorf("The raw spread of the symbol %v was not found on %v", symbol, dateTime)
	//}
	return avgSpreadOf14Day, nil
}

func (ss *spreadService) GetAvgSpreadLast14Days(symbol string, lp int8, dateTime time.Time) (int64, error) {
	//股票品种 .替换为_
	newSymbol := strings.Replace(symbol, ".", "_", 1)
	lpSpreadRecord, err := repo.LpSpreadRecordRepo.GetAvgSpreadLast14Days(newSymbol, lp, dateTime)
	if err != nil {
		return 0, errors.WithMessage(err, util.ErrWrapMessage())
	}

	if len(lpSpreadRecord) == 0 {
		return 0, nil
	}

	totalAvgSpread := int64(0)
	for _, v := range lpSpreadRecord {
		totalAvgSpread = totalAvgSpread + v.LpAvg
	}
	//向上取整
	avgSpreadLast14Day := math.Ceil(float64(totalAvgSpread) / float64(len(lpSpreadRecord)))
	return int64(avgSpreadLast14Day), nil
}

// 初始化品种点差修改列表
func (ss spreadService) InitSpreadList() error {
	spreads, err := repo.SpreadRecordRepo.GetList()
	if err != nil {
		return errors.WithMessage(err, util.ErrWrapMessage())
	}

	SpreadList = spreads
	return nil
}

func (ss *spreadService) getHistoricalSymbolLp(symbol string, dateTime time.Time) (int8, error) {
	switchLpList, err := LpService.GetSwitchLpList(symbol)
	if err != nil {
		return 0, errors.WithMessage(err, util.ErrWrapMessage())
	}

	if len(switchLpList) == 0 {
		//如果品种无lp的更改记录，从缓存里面获取
		symbolLp, err := CacheService.GetSymbolLp(symbol)
		if err != nil {
			return 0, errors.WithMessage(err, util.ErrWrapMessage())
		}
		return symbolLp, nil
	}

	for _, v := range switchLpList {
		if dateTime.After(v.SwitchTime) {
			return v.NewLp, nil
		}
	}
	return switchLpList[len(switchLpList)-1].OldLp, nil
}
