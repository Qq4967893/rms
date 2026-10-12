package service

import (
	pb "backend-rms/grpc/et6api"
	"backend-rms/model"
	"backend-rms/util"
	"math"
	"strings"

	"github.com/golang-module/carbon/v2"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"golang.org/x/exp/slices"
	"golang.org/x/sync/errgroup"
)

type CalculatedProfit struct {
	ClosePrice  decimal.Decimal
	GrossProfit decimal.Decimal
}

type PipProfit struct {
	Symbol         string
	ProfitCurrency string
	Did            decimal.Decimal
	GrossProfit    decimal.Decimal
}

type profitService struct{}

var ProfitService = &profitService{}

func (qs *profitService) CalcProfit(trade model.OpeningTrade, spreadInfo *pb.GetSpreadInfoResponse_ResponseData, group string) (cProfit CalculatedProfit, err error) {
	var quote *Quote
	var symbolInfo *SymbolInfo

	errGroup := new(errgroup.Group)
	errGroup.Go(func() error {
		var err error
		quote, err = QuoteService.GetOffsetQuote(trade.Symbol, spreadInfo, group)
		if err != nil {
			err = errors.Errorf("failed to get quote: %v", err)
		}
		return err
	})

	errGroup.Go(func() error {
		var err error
		symbolInfo, err = CacheService.GetSymbolInfo(trade.Symbol)
		if err != nil {
			err = errors.Errorf("failed to get symbol info: %v", err)
		}
		return err
	})

	if err = errGroup.Wait(); err != nil {
		return CalculatedProfit{}, errors.WithMessage(err, util.ErrWrapMessage())
	}

	cProfit, err = qs.calcProfit(trade, *quote, *symbolInfo, nil)
	return cProfit, errors.WithMessage(err, util.ErrWrapMessage())
}

func (qs *profitService) CalcProfitBatch(trades []model.OpeningTrade, spreadInfo *pb.GetSpreadInfoResponse_ResponseData, allQuotes map[string]*Quote) (map[int64]CalculatedProfit, error) {
	groupQuotes, err := QuoteService.GetAllOffsetQuotes(allQuotes, spreadInfo)
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}

	symbolLeverageMap, _ := SymbolLeverageService.GetAllMap()
	tradeProfits := make(map[int64]CalculatedProfit)
	for _, trade := range trades {
		symbolInfo, err := CacheService.GetSymbolInfo(trade.Symbol)
		if err != nil {
			return nil, errors.WithMessage(err, util.ErrWrapMessage())
		}

		symbolOriginal := symbolLeverageMap[trade.Symbol]

		var quote Quote
		if len(trade.Group) == 0 {
			return nil, errors.Errorf("can't find the group field in a trade(%v)", trade.Ticket)
		}
		if offsetQuote, ok := groupQuotes[trade.Group]; ok {
			quote = *offsetQuote[symbolOriginal]
		} else {
			return nil, errors.Errorf("can't find a quote for group %v", trade.Group)
		}

		cProfit, err := qs.calcProfit(trade, quote, *symbolInfo, allQuotes)
		if err != nil {
			return nil, errors.WithMessage(err, util.ErrWrapMessage())
		}

		tradeProfits[trade.Ticket] = cProfit
	}
	return tradeProfits, nil
}

func (qs *profitService) CalcPipProfitByQuotes(symbolInfo SymbolInfo, cmd int8, quotes map[string]*Quote) (PipProfit, error) {
	lots := decimal.NewFromFloat(1)
	//品种去杠杆
	symbolOrigin, err := SymbolLeverageService.GetSymbolOrigin(symbolInfo.Symbol)
	if err != nil {
		return PipProfit{}, errors.WithMessage(err, util.ErrWrapMessage())
	}

	symbol := symbolOrigin
	digits := symbolInfo.Digits
	contractSize := decimal.NewFromFloat(symbolInfo.ContractSize)
	profitCurrency := symbolInfo.ProfitCurrency

	//只计算单个pip的单位值
	pipPrice := decimal.NewFromFloat(1 * math.Pow(10, float64(-digits)))
	profitCalc := pipPrice.Mul(contractSize).Mul(lots)

	var bid decimal.Decimal
	if _, ok := quotes[profitCurrency]; ok {
		bid = quotes[profitCurrency].Bid
	} else {
		bid = decimal.Zero
	}

	var grossProfit decimal.Decimal
	if profitCurrency == "USD" {
		grossProfit = profitCalc
	} else {
		if len(quotes) == 0 {
			cmd = 0
			start := carbon.Now(carbon.UTC).StartOfDay().TimestampMilli()
			bid, err = HistoryQuoteService.GetLastClosePrice(profitCurrency, start)
			if err != nil {
				return PipProfit{}, errors.WithMessage(err, util.ErrWrapMessage())
			}
			quotes = map[string]*Quote{
				profitCurrency: {
					Bid: bid,
				},
			}
		}

		grossProfit, err = qs.exchangeToDepositCurrency(profitCalc, profitCurrency, quotes)
		if err != nil {
			return PipProfit{}, errors.WithMessage(err, util.ErrWrapMessage())
		}
	}

	return PipProfit{
		Symbol:         symbol,
		ProfitCurrency: profitCurrency,
		Did:            bid,
		GrossProfit:    grossProfit,
	}, nil
}

func (qs *profitService) exchangeToDepositCurrency(profitCalc decimal.Decimal, profitCalcCurrency string, allQuotes map[string]*Quote) (decimal.Decimal, error) {
	return qs.exchangeToDepositCurrencyIndirect(profitCalc, profitCalcCurrency, allQuotes)
}

func (qs *profitService) exchangeToDepositCurrencyIndirect(profitCalc decimal.Decimal, profitCalcCurrency string, allQuotes map[string]*Quote) (decimal.Decimal, error) {
	var err error
	var grossProfit decimal.Decimal
	if profitCalcCurrency == "USD" {
		grossProfit = profitCalc
	} else {
		grossProfit, err = qs.toCurrency(profitCalc, profitCalcCurrency, "USD", allQuotes)
		if err != nil {
			return decimal.Zero, errors.WithMessage(err, util.ErrWrapMessage())
		}
	}

	//aimCurrency := "GBP"
	//grossProfit, err = qs.usdToOtherCurrency(grossProfit, allQuotes, aimCurrency)
	//if err != nil {
	//	return decimal.Zero, errors.WithMessage(err, util.ErrWrapMessage())
	//}

	return grossProfit, nil
}

func (qs *profitService) toCurrency(profitCalc decimal.Decimal, profitCalcCurrency string, aimCurrency string, allQuotes map[string]*Quote) (decimal.Decimal, error) {
	profitQuote, ok := allQuotes[profitCalcCurrency]
	if !ok {
		var err error
		profitQuote, err = QuoteService.GetQuote(profitCalcCurrency)
		if err != nil {
			return decimal.Zero, errors.WithMessage(err, util.ErrWrapMessage())
		}
	}

	var priceTrans decimal.Decimal
	priceTrans = profitQuote.Bid.Add(profitQuote.Ask).Div(decimal.NewFromInt(2))

	var grossProfit decimal.Decimal
	if strings.HasPrefix(profitCalcCurrency, aimCurrency) {
		grossProfit = profitCalc.Div(priceTrans)
	} else {
		grossProfit = profitCalc.Mul(priceTrans)
	}
	return grossProfit, nil
}

func (qs *profitService) usdToOtherCurrency(grossProfit decimal.Decimal, allQuotes map[string]*Quote, aimCurrency string) (decimal.Decimal, error) {
	otherCurrencies := []string{"GBP", "EUR"}
	if !slices.Contains(otherCurrencies, aimCurrency) {
		return decimal.Zero, errors.Errorf("the target currency(%v) is not supported", aimCurrency)
	}

	var symbol string
	if aimCurrency == "GBP" {
		symbol = "GBPUSD"
	} else {
		symbol = "EURUSD"
	}

	profitQuote, ok := allQuotes[symbol]
	if !ok {
		var err error
		profitQuote, err = QuoteService.GetQuote(symbol)
		if err != nil {
			return decimal.Zero, errors.WithMessage(err, util.ErrWrapMessage())
		}
	}

	var priceTrans decimal.Decimal
	priceTrans = profitQuote.Bid.Add(profitQuote.Ask).Div(decimal.NewFromInt(2))

	grossProfit = grossProfit.Div(priceTrans)
	return grossProfit, nil
}

func (qs *profitService) calcProfit(trade model.OpeningTrade, quote Quote, symbolInfo SymbolInfo, allQuotes map[string]*Quote) (cProfit CalculatedProfit, err error) {
	contractSize := decimal.NewFromFloat(symbolInfo.ContractSize)
	lots := decimal.NewFromInt(trade.Volume).Mul(decimal.NewFromFloat(0.01))

	var closePrice decimal.Decimal
	var profitCalc decimal.Decimal
	if trade.Cmd == 0 {
		closePrice = quote.Bid
		profitCalc = closePrice.Sub(trade.OpenPrice).Mul(contractSize).Mul(lots)
	} else {
		closePrice = quote.Ask
		profitCalc = trade.OpenPrice.Sub(closePrice).Mul(contractSize).Mul(lots)
	}

	profitCalcCurrency := symbolInfo.ProfitCurrency
	grossProfit, err := qs.exchangeToDepositCurrency(profitCalc, profitCalcCurrency, allQuotes)
	if err != nil {
		return cProfit, errors.WithMessage(err, util.ErrWrapMessage())
	}

	cProfit.ClosePrice = closePrice
	cProfit.GrossProfit = grossProfit
	return
}
