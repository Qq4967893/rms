package service

import (
	"backend-rms/core/redis"
	pb "backend-rms/grpc/et6api"
	"backend-rms/model"
	"backend-rms/util"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type Quote = model.Tick

type SymbolQuotes map[string]*Quote

type quoteService struct{}

var QuoteService = &quoteService{}

const quoteRedisKey = "et6:tick"

func (qs *quoteService) GetQuote(symbolOriginal string) (*Quote, error) {
	v, err := redis.QuoteRedisCli.HGet(quoteRedisKey, symbolOriginal)
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}

	quote := strings.Split(v, ",")
	if len(quote) != 3 {
		return nil, errors.Errorf("acquire quote failed: symbol=%s, quote=%s", symbolOriginal, quote)
	}

	ask, _ := decimal.NewFromString(quote[0])
	bid, _ := decimal.NewFromString(quote[1])
	ctm, _ := strconv.ParseInt(quote[2], 10, 64)

	return &Quote{Ask: ask, Bid: bid, Ctm: ctm}, nil
}

func (qs *quoteService) GetAllQuotes() (map[string]*Quote, error) {
	v, err := redis.QuoteRedisCli.HGetAll(quoteRedisKey)
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}

	quotes := make(map[string]*Quote)
	for symbolOriginal, q := range v {
		quote := strings.Split(q, ",")
		if len(quote) != 3 {
			return nil, errors.Errorf("acquire quote failed: symbol=%s, quote=%s", symbolOriginal, quote)
		}

		ask, _ := decimal.NewFromString(quote[0])
		bid, _ := decimal.NewFromString(quote[1])
		ctm, _ := strconv.ParseInt(quote[2], 10, 64)

		quotes[symbolOriginal] = &Quote{Ask: ask, Bid: bid, Ctm: ctm}
	}
	return quotes, nil
}

func (qs *quoteService) GetOffsetQuote(symbolOriginal string, spreadInfo *pb.GetSpreadInfoResponse_ResponseData, group string) (quote *Quote, err error) {
	quote, err = qs.GetQuote(symbolOriginal)
	if err != nil {
		return quote, errors.WithMessage(err, util.ErrWrapMessage())
	}

	securityMap := make(map[string]string)
	for security, symbols := range spreadInfo.GetSecuritySymbol() {
		for _, s := range symbols.GetSymbol() {
			securityMap[s] = security
		}
	}
	security, ok := securityMap[symbolOriginal]
	if !ok {
		return
	}

	groupSecuritySpreadOffset := spreadInfo.GetGroupSecuritySpreadOffset()
	if groupSecuritySpreadOffset == nil {
		err = errors.Errorf("GroupSecuritySpreadOffset of SpreadInfo from GRPC doesn't exist")
		return
	}

	securityOffset, ok := groupSecuritySpreadOffset[group]
	if !ok {
		return
	}

	var spread = decimal.Zero
	if offset, okk := securityOffset.GetData()[security]; okk {
		digits, _ := SymbolOriginalService.GetDigitsByOriginal(symbolOriginal)
		spread = decimal.NewFromInt(10).Pow(decimal.NewFromInt(-digits)).Mul(decimal.NewFromInt32(offset))
	}

	quote.Ask = quote.Ask.Add(spread.Div(decimal.NewFromInt(2)))
	quote.Bid = quote.Bid.Sub(spread.Div(decimal.NewFromInt(2)))
	return
}

func (qs *quoteService) GetAllOffsetQuotes(allQuotes map[string]*Quote, spreadInfo *pb.GetSpreadInfoResponse_ResponseData) (map[string]SymbolQuotes, error) {
	if len(allQuotes) == 0 {
		var err error
		allQuotes, err = qs.GetAllQuotes()
		if err != nil {
			return nil, errors.WithMessage(err, util.ErrWrapMessage())
		}
	}

	securityMap := make(map[string]string)
	for security, symbols := range spreadInfo.GetSecuritySymbol() {
		for _, s := range symbols.GetSymbol() {
			securityMap[s] = security
		}
	}

	groupSecuritySpreadOffset := spreadInfo.GetGroupSecuritySpreadOffset()
	if groupSecuritySpreadOffset == nil {
		return nil, errors.Errorf("GroupSecuritySpreadOffset of SpreadInfo from GRPC doesn't exist")
	}

	symbolOriginals, err := SymbolLeverageService.GetAllMap()
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}

	symbolOriginalMap, err := SymbolOriginalService.GetAllMap(false)
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}

	var groupQuotes = make(map[string]SymbolQuotes)
	for group, securityOffset := range groupSecuritySpreadOffset {
		var symbolQuotes = make(SymbolQuotes)

		for symbol, quote := range allQuotes {
			if symbol == "AUSD" {
				continue
			}
			security, ok := securityMap[symbol]
			if !ok {
				continue
			}

			//特殊组品种点差去除杠杆（原因：allQuotes内的品种是没杠杆的品种，而spreadInfo内记录的是带杠杆的（如：spreadInfo内只有GOOG.nd5的点差值，而没有GOOG.nd的。这是由于该股票GOOG.nd5只能交易5倍杠杆，不能交易100倍杠杆所以没有写入100倍杠杆的点差值））
			groupSymbolSpread := make(map[string]int32, 0)
			for k, v := range securityOffset.GetData() {
				symbolOriginal := symbolOriginals[k]
				if _, ok := groupSymbolSpread[symbolOriginal]; ok {
					continue
				}
				groupSymbolSpread[symbolOriginal] = v
			}

			var spread = decimal.Zero
			if offset, okk := groupSymbolSpread[security]; okk {
				digits := symbolOriginalMap[symbol].Digits
				spread = decimal.NewFromInt(10).Pow(decimal.NewFromInt(-digits)).Mul(decimal.NewFromInt32(offset))
			}

			symbolQuotes[symbol] = &Quote{
				Ask: quote.Ask.Add(spread.Div(decimal.NewFromInt(2))),
				Bid: quote.Bid.Sub(spread.Div(decimal.NewFromInt(2))),
				Ctm: quote.Ctm}
		}
		groupQuotes[group] = symbolQuotes
	}
	return groupQuotes, nil
}
