package market

import (
	"backend-rms/core/logger"
	"backend-rms/service"
	"backend-rms/util"
	"backend-rms/ws"
	"backend-rms/ws/controller"
	"backend-rms/ws/request"
	"backend-rms/ws/request/market"
	"backend-rms/ws/response"
	"math"
	"sort"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/exp/maps"
	"golang.org/x/sync/errgroup"
)

type SummarySymbol struct {
	NotStockSymbols       []string //非股票的所有品种
	notStockInvalidSymbol []string //非股票的非上架品种
}

type SummaryController struct{}

func (sc *SummaryController) GetSummaryList(c *ws.Context) {
	controller.TickerHandle(c, 2*time.Second, func(c *ws.Context) (data any, err error) {
		req := market.GetSummaryListRequest{}
		if ok := request.Validate(c, &req); !ok {
			return
		}

		return sc.buildSummaryList(req.Params.Books, req.Params.Sort, req.Params.Order)
	})
}

func (sc *SummaryController) buildSummaryList(books []string, sortField, order string) ([]response.MarketSummary, error) {
	//获取品种买卖手数，et6净头寸，持仓单订单数，mt4净头寸，买卖的实时报价
	g := new(errgroup.Group)

	var symbols SummarySymbol
	var err error
	g.Go(func() error {
		symbols, err = sc.buildSymbol()
		return errors.WithMessage(err, util.ErrWrapMessage())
	})

	var allQuotes map[string]*service.Quote
	g.Go(func() error {
		allQuotes, err = service.QuoteService.GetAllQuotes()
		return errors.WithMessage(err, util.ErrWrapMessage())
	})

	bookSummaries := make(map[string]service.Summary)
	g.Go(func() error {
		bookSummaries, err = service.CacheService.GetBookSummaries(books)
		return errors.WithMessage(err, util.ErrWrapMessage())
	})

	if err := g.Wait(); err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}

	//将Map中的键复制到一个新的切片中
	bookSummaryKeys := make([]string, 0, len(bookSummaries))
	for key := range bookSummaries {
		bookSummaryKeys = append(bookSummaryKeys, key)
	}

	//对切片进行排序
	sort.Slice(bookSummaryKeys, func(i, j int) bool {
		return bookSummaryKeys[i] < bookSummaryKeys[j]
	})

	summaries := make([]response.MarketSummary, 0)
	for _, key := range bookSummaryKeys {
		summary := bookSummaries[key]
		//判断是否在切片内
		if !util.InSlice(summary.Symbol, symbols.NotStockSymbols) {
			continue
		}

		if util.InSlice(summary.Symbol, symbols.notStockInvalidSymbol) && summary.Orders == 0 {
			continue
		}

		symbolQuote, ok := allQuotes[summary.Symbol]
		if !ok {
			logger.Errorf("the %s quote is not exist", summary.Symbol)
			continue
		}

		var symbolSummary = response.MarketSummary{
			Symbol:     summary.Symbol,
			BuyLots:    summary.BuyLots,
			SellLots:   summary.SellLots,
			Et6NetLots: summary.Et6NetLots,
			BuyPrice:   symbolQuote.Bid,
			SellPrice:  symbolQuote.Ask,
			Orders:     summary.Orders,
		}
		summaries = append(summaries, symbolSummary)
	}

	if sortField == "et6_net_lots" {
		sc.SortET6NetLots(summaries, order)
	} else {
		err = util.DataOrderBy(summaries, sortField, order)
		if err != nil {
			return nil, errors.WithMessage(err, util.ErrWrapMessage())
		}
	}
	return summaries, nil
}

func (sc *SummaryController) buildSymbol() (SummarySymbol, error) {
	//获取所有品种
	symbolMaps, err := service.SymbolOriginalService.GetAllMap(false)
	if err != nil {
		return SummarySymbol{}, errors.WithMessage(err, util.ErrWrapMessage())
	}
	allSymbols := maps.Keys(symbolMaps)

	//获取股票品种
	stockSymbols, err := service.SymbolTypeService.GetStockSymbols(false)
	if err != nil {
		return SummarySymbol{}, errors.WithMessage(err, util.ErrWrapMessage())
	}

	//获取非股票的所有品种
	notStockSymbols := sc.getDiffSymbol(allSymbols, stockSymbols)

	validSymbolMaps, err := service.SymbolOriginalService.GetAllMap(true)
	if err != nil {
		return SummarySymbol{}, errors.WithMessage(err, util.ErrWrapMessage())
	}
	var validSymbols []string
	for k := range validSymbolMaps {
		validSymbols = append(validSymbols, k)
	}
	//获取排除股票的非上架品种
	notStockInvalidSymbol := sc.getDiffSymbol(notStockSymbols, validSymbols)
	return SummarySymbol{notStockSymbols, notStockInvalidSymbol}, nil
}

func (sc *SummaryController) SortET6NetLots(summaries []response.MarketSummary, order string) {
	sort.Slice(summaries, func(i, j int) bool {
		after, _ := summaries[i].Et6NetLots.Float64()
		after = math.Abs(after)

		before, _ := summaries[j].Et6NetLots.Float64()
		before = math.Abs(before)

		if order == "desc" {
			if before <= after {
				return true
			}
		} else {
			if after <= before {
				return true
			}
		}
		return false
	})
}

func (sc *SummaryController) getDiffSymbol(symbolsOne, symbolsTwo []string) []string {
	temp := map[string]struct{}{}
	for _, v := range symbolsTwo {
		if _, ok := temp[v]; !ok {
			temp[v] = struct{}{}
		}
	}

	var diffSymbols []string
	for _, v := range symbolsOne {
		if _, ok := temp[v]; !ok {
			diffSymbols = append(diffSymbols, v)
		}
	}
	return diffSymbols
}
