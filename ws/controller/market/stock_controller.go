package market

import (
	"backend-rms/core/logger"
	"backend-rms/model"
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

type StockController struct{}

func (sc *StockController) GetStockList(c *ws.Context) {
	controller.TickerHandle(c, 2*time.Second, func(c *ws.Context) (data any, err error) {
		req := market.GetStockListRequest{}
		if ok := request.Validate(c, &req); !ok {
			return
		}

		return sc.buildStockList(req.Params.Books, req.Params.Sort, req.Params.Order)
	})
}

func (sc *StockController) buildStockList(books []string, sort, order string) ([]response.MarketStock, error) {
	g := new(errgroup.Group)

	var stockSymbols []string
	var err error
	g.Go(func() error {
		stockSymbols, err = service.SymbolTypeService.GetStockSymbols(true)
		return errors.WithMessage(err, util.ErrWrapMessage())
	})

	var symbolMaps map[string]*model.SymbolOriginal
	var allSymbols []string
	g.Go(func() error {
		symbolMaps, err = service.SymbolOriginalService.GetAllMap(false)
		allSymbols = maps.Keys(symbolMaps)
		return errors.WithMessage(err, util.ErrWrapMessage())
	})

	var summaries map[string]service.Summary
	g.Go(func() error {
		summaries, err = service.CacheService.GetBookSummaries(books)
		return errors.WithMessage(err, util.ErrWrapMessage())
	})

	var exposures map[string][]service.Exposure
	g.Go(func() error {
		exposures, err = service.CacheService.GetBookExposures(books)
		return errors.WithMessage(err, util.ErrWrapMessage())
	})

	var allQuotes map[string]*service.Quote
	g.Go(func() error {
		allQuotes, err = service.QuoteService.GetAllQuotes()
		return errors.WithMessage(err, util.ErrWrapMessage())
	})

	if err := g.Wait(); err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}

	invalidStockSymbol := sc.getDiffSymbol(stockSymbols, allSymbols)
	stockSummaries := make([]response.MarketStock, 0)
	for _, v := range exposures["stock"] {
		symbolSummary := summaries[v.Currency]
		if !util.InSlice(v.Currency, stockSymbols) {
			continue
		}

		if util.InSlice(v.Currency, invalidStockSymbol) && symbolSummary.Orders == 0 {
			continue
		}

		symbolQuote, ok := allQuotes[symbolSummary.Symbol]
		if !ok {
			logger.Errorf("the %s quote is not exist", symbolSummary.Symbol)
			continue
		}

		var middle = response.MarketStock{
			Symbol:      symbolSummary.Symbol,
			BuyLots:     symbolSummary.BuyLots,
			SellLots:    symbolSummary.SellLots,
			Et6NetLots:  symbolSummary.Et6NetLots,
			LPNetLots:   symbolSummary.LPNetLots,
			Orders:      symbolSummary.Orders,
			BuyPrice:    symbolQuote.Ask,
			SellPrice:   symbolQuote.Bid,
			ExposureUSD: v.ExposureUSD,
		}

		stockSummaries = append(stockSummaries, middle)
	}

	if sort == "et6_net_lots" {
		sc.SortStockET6NetLots(stockSummaries, order)
	} else {
		err = util.DataOrderBy(stockSummaries, sort, order)
		if err != nil {
			return nil, errors.WithMessage(err, util.ErrWrapMessage())
		}
	}
	return stockSummaries, nil
}

func (sc *StockController) SortStockET6NetLots(stockSummaries []response.MarketStock, order string) []response.MarketStock {
	sort.Slice(stockSummaries, func(i, j int) bool {
		after, _ := stockSummaries[i].Et6NetLots.Float64()
		after = math.Abs(after)

		before, _ := stockSummaries[j].Et6NetLots.Float64()
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
	return stockSummaries
}

func (sc *StockController) getDiffSymbol(symbolsOne, symbolsTwo []string) []string {
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
