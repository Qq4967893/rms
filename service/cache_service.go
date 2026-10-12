package service

import (
	"backend-rms/core/logger"
	"backend-rms/core/redis"
	"backend-rms/grpc/tsp"
	"backend-rms/util"
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"github.com/elliotchance/phpserialize"
)

type SymbolInfo struct {
	Symbol         string  `php:"symbol"`
	Digits         int32   `php:"digits"`
	ContractSize   float64 `php:"contract_size"`
	LotMin         float64 `php:"lot_min"`
	LotMax         float64 `php:"lot_max"`
	LotStep        float64 `php:"lot_step"`
	SwapType       int32   `php:"swap_type"`
	SwapLong       float64 `php:"swap_long"`
	SwapShort      float64 `php:"swap_short"`
	MarginDivider  float64 `php:"margin_divider"`
	MarginInitial  float64 `php:"margin_initial"`
	MarginMode     int32   `php:"margin_mode"`
	StopLevel      int32   `php:"stops_level"`
	ProfitCurrency string  `php:"profit_currency"`
}

type Summary struct {
	Symbol     string          `json:"symbol"`
	BuyLots    decimal.Decimal `json:"buy_lots"`
	SellLots   decimal.Decimal `json:"sell_lots"`
	Et6NetLots decimal.Decimal `json:"et6_net_lots"`
	LPNetLots  decimal.Decimal `json:"lp_net_lots"`
	BuyPrice   decimal.Decimal `json:"buy_price"`
	SellPrice  decimal.Decimal `json:"sell_price"`
	Orders     int             `json:"orders"`
}

type Exposure struct {
	Currency    string          `json:"currency"`
	Exposure    decimal.Decimal `json:"exposure"`
	ExposureUSD decimal.Decimal `json:"exposure_usd"`
	Symbol      string          `json:"symbol"`
}

type cacheService struct{}

var CacheService = &cacheService{}

const (
	TradePLKey   = "et6:trade_pl"
	ClientPLKey  = "et6:client_pl"
	SymbolsLpKey = "source_lp"
	SummariesKey = "et6:summary"
	ExposuresKey = "et6:exposure"
)

var LpPriceKeys = map[int8]string{
	2:  "cfh_quote_prices",
	11: "iq_quote_prices",
	15: "mix_quote_prices",
}

func (cs *cacheService) GetSymbolInfo(symbol string) (*SymbolInfo, error) {
	key := "grpc:getsymbolinfos:" + symbol

	res, err := redis.BaseRedisCli.Get(key)
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}

	si := &SymbolInfo{}
	if res != "" {
		err = phpserialize.Unmarshal([]byte(res), si)
		if err != nil {
			return nil, errors.Wrap(err, "json: json unmarshal fail")
		}
		return si, nil
	}

	symbolInfos, err := tsp.GetSymbolInfos(context.Background(), "test request id", symbol, true)
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}

	if len(symbolInfos) > 0 {
		si.Symbol = symbolInfos[0].GetSymbol()
		si.Digits = symbolInfos[0].GetDigits()
		si.ContractSize = symbolInfos[0].GetContractSize()
		si.LotMin = symbolInfos[0].GetLotMin()
		si.LotMax = symbolInfos[0].GetLotMax()
		si.LotStep = symbolInfos[0].GetLotStep()
		si.SwapType = symbolInfos[0].GetSwapType()
		si.SwapLong = symbolInfos[0].GetSwapLong()
		si.SwapShort = symbolInfos[0].GetSwapShort()
		si.MarginDivider = symbolInfos[0].GetMarginDivider()
		si.MarginInitial = symbolInfos[0].GetMarginInitial()
		si.MarginMode = symbolInfos[0].GetMarginMode()
		si.StopLevel = symbolInfos[0].GetStopLevel()

		si.ProfitCurrency, err = SymbolOriginalService.ConvertCurrency(symbolInfos[0].GetProfitCurrency())
		if err != nil {
			return nil, errors.WithMessage(err, util.ErrWrapMessage())
		}
	}

	b, err := phpserialize.Marshal(si, &phpserialize.MarshalOptions{OnlyStdClass: true})
	if err != nil {
		return nil, errors.Wrap(err, "serialize: Marshal failed")
	}
	redis.BaseRedisCli.Set(key, string(b), time.Hour)
	return si, nil
}

func (cs *cacheService) GetSymbolLp(symbol string) (int8, error) {
	//symbolLp, err := redis.PositionRedisCli.HGet(SymbolsLpKey, symbol)
	//if err != nil {
	//	return 0, errors.WithMessage(err, util.ErrWrapMessage())
	//}
	//
	//symbolLpValue, err := strconv.Atoi(symbolLp)
	//if err != nil {
	//	return 0, errors.WithMessage(err, util.ErrWrapMessage())
	//}
	//return int8(symbolLpValue), nil

	//交易端现在报价只有一个LP mix=15
	return 15, nil
}

func (cs *cacheService) GetAllAccountFloatingPLMap() (map[int64]decimal.Decimal, error) {
	result, err := redis.RMSRedisCli.HGetAll(ClientPLKey)
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}
	floatingPLMap := make(map[int64]decimal.Decimal, len(result))
	for field, value := range result {
		login, _ := strconv.ParseInt(field, 10, 64)
		floatingPL, _ := decimal.NewFromString(value)
		floatingPLMap[login] = floatingPL
	}
	return floatingPLMap, nil
}

func (cs *cacheService) GetAccountFloatingPLMap(logins []int64) (map[int64]decimal.Decimal, error) {
	var fields []string
	for _, login := range logins {
		fields = append(fields, strconv.FormatInt(login, 10))
	}
	result, err := redis.RMSRedisCli.HMGet(ClientPLKey, fields...)
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}
	floatingPLMap := make(map[int64]decimal.Decimal, len(logins))
	for i, login := range logins {
		if result[i] == nil {
			continue
		}
		floatingPL, _ := decimal.NewFromString(result[i].(string))
		floatingPLMap[login] = floatingPL
	}
	return floatingPLMap, nil
}

func (cs *cacheService) GetTradeFloatingPLMap(tickets []int64) (map[int64]decimal.Decimal, error) {
	var fields []string
	for _, login := range tickets {
		fields = append(fields, strconv.FormatInt(login, 10))
	}
	result, err := redis.RMSRedisCli.HMGet(TradePLKey, fields...)
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}
	floatingPLMap := make(map[int64]decimal.Decimal, len(tickets))
	for i, login := range tickets {
		if result[i] == nil {
			continue
		}
		floatingPL, _ := decimal.NewFromString(result[i].(string))
		floatingPLMap[login] = floatingPL
	}
	return floatingPLMap, nil
}

func (cs *cacheService) GetBookSummaries(books []string) (map[string]Summary, error) {
	summaries := make(map[string]Summary)
	for _, v := range books {
		jsonSummaries, err := redis.RMSRedisCli.Get(SummariesKey + "_" + v)
		if err != nil {
			return nil, errors.WithMessage(err, util.ErrWrapMessage())
		}

		var bookSummaries []Summary
		err = json.Unmarshal([]byte(jsonSummaries), &bookSummaries)
		if err != nil {
			return nil, errors.Wrap(err, "json: unmarshal json failed")
		}

		for _, bv := range bookSummaries {
			if _, ok := summaries[bv.Symbol]; ok {
				var symbolSummary = summaries[bv.Symbol]
				summaries[bv.Symbol] = Summary{
					Symbol:     bv.Symbol,
					BuyLots:    symbolSummary.BuyLots.Add(bv.BuyLots),
					SellLots:   symbolSummary.SellLots.Add(bv.SellLots),
					Et6NetLots: symbolSummary.Et6NetLots.Add(bv.Et6NetLots),
					LPNetLots:  symbolSummary.LPNetLots.Add(bv.LPNetLots),
					Orders:     symbolSummary.Orders + bv.Orders,
				}
				continue
			}
			summaries[bv.Symbol] = Summary{
				Symbol:     bv.Symbol,
				BuyLots:    bv.BuyLots,
				SellLots:   bv.SellLots,
				Et6NetLots: bv.Et6NetLots,
				LPNetLots:  bv.LPNetLots,
				Orders:     bv.Orders,
			}
		}
	}
	return summaries, nil
}

func (cs *cacheService) GetBookExposures(books []string) (map[string][]Exposure, error) {
	exposures := make(map[string][]Exposure)
	for _, v := range books {
		jsonExposures, err := redis.RMSRedisCli.HGetAll(ExposuresKey + "_" + v)
		if err != nil {
			return nil, errors.WithMessage(err, util.ErrWrapMessage())
		}

		for jk, jv := range jsonExposures {
			var bookExposure []Exposure
			err = json.Unmarshal([]byte(jv), &bookExposure)
			if err != nil {
				return nil, err
			}

			//判断exposure的类型是否存在
			if _, ok := exposures[jk]; ok {
				for ek := range exposures[jk] {
					for _, bv := range bookExposure {
						//判断该类型的exposure品种一致
						if exposures[jk][ek].Currency == bv.Currency {
							exposures[jk][ek].Exposure = exposures[jk][ek].Exposure.Add(bv.Exposure)
							exposures[jk][ek].ExposureUSD = exposures[jk][ek].ExposureUSD.Add(bv.ExposureUSD)
							continue
						}
					}
				}
				continue
			}
			exposures[jk] = bookExposure
		}
	}
	return exposures, nil
}

func (cs *cacheService) GetETBSpread(symbol string) (int64, error) {
	symbolLp, err := cs.GetSymbolLp(symbol)
	if err != nil {
		return 0, errors.WithMessage(err, util.ErrWrapMessage())
	}

	if _, ok := LpPriceKeys[symbolLp]; !ok {
		return 0, errors.Errorf("lp %v is not found for %v", symbolLp, symbol)
	}

	key := LpPriceKeys[symbolLp]
	priceJson, err := redis.PositionRedisCli.HGet(key, symbol)
	if err != nil {
		return 0, errors.WithMessage(err, util.ErrWrapMessage())
	}

	var priceMap struct {
		Source    string
		Bid       decimal.Decimal
		Ask       decimal.Decimal
		Digit     int8
		Etbspread int64
		Timestamp int64
		Lp        int8
	}
	err = json.Unmarshal([]byte(priceJson), &priceMap)
	if err != nil {
		return 0, errors.Wrap(err, "json: unmarshal json failed")
	}
	return priceMap.Etbspread, nil
}

func (cs *cacheService) ValidateAuthCode(codeId string, code string) (bool, error) {
	authCodeKey := "auth-code-id-" + codeId
	redisV, err := redis.RMSRedisCli.Get(authCodeKey)
	if err != nil {
		return false, errors.WithMessage(err, util.ErrWrapMessage())
	}

	//删除已通过的验证码
	redis.RMSRedisCli.Del(authCodeKey)

	if redisV != code {
		return false, nil
	}
	return true, nil
}

func (cs *cacheService) SetSummaryBook(book string) {
	bookSummaryRedisValue := make([]*InitSummary, 0)
	for _, symbolInfo := range SummaryBook[book] {
		bookSummaryRedisValue = append(bookSummaryRedisValue, symbolInfo)
	}

	jsonBookSummaryRedisValue, err := json.Marshal(bookSummaryRedisValue)
	if err != nil {
		logger.Errorf("JSON marshal bookSummaryRedisValue failed: %v", err)
		return
	}

	err = redis.RMSRedisCli.Set(SummariesKey+"_"+strings.ToLower(book), jsonBookSummaryRedisValue, 0)
	if err != nil {
		logger.Errorf("Redis set et6:summary_"+strings.ToLower(book)+" is err:%v", err)
		return
	}
}

func (cs *cacheService) SetExposureBook(book string, exposureBook map[string][]*InitExposure) {
	for key, value := range exposureBook {
		jsonBookExposureRedisValue, err := json.Marshal(value)
		if err != nil {
			logger.Errorf("JSON marshal bookExposureRedisValue failed: %v", err)
			return
		}

		_, err = redis.RMSRedisCli.HSet(ExposuresKey+"_"+strings.ToLower(book), key, jsonBookExposureRedisValue)
		if err != nil {
			logger.Error("Redis hset et6:exposure_"+book+" is error:", err)
			return
		}
	}
}
