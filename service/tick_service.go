package service

import (
	"backend-rms/model"
	"backend-rms/repo"
	"backend-rms/util"
	"context"
	"strings"
	"sync"
	"time"

	"github.com/golang-module/carbon/v2"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

type tickService struct{}

var TickService = &tickService{}

func (s *tickService) GetHistoryTick(ctx context.Context, symbol string, dateTime time.Time) (tick *model.Tick, err error) {
	collName := s.getCollName(symbol)
	ctm := carbon.Time2Carbon(dateTime).EndOfSecond().Carbon2Time().UnixMilli()
	tick, err = repo.TickRepo.GetHistoryTick(ctx, collName, ctm)
	if err != nil {
		return nil, errors.WithMessagef(err, "mongo: (%s) tick before %d was not found", collName, ctm)
	}
	if tick == nil {
		return nil, errors.Errorf("mongo: (%s) tick before %d was not found", collName, ctm)
	}
	return
}

func (s *tickService) GetAllHistoryTicks(ctx context.Context, dateTime time.Time) (ticks map[string]*model.Tick, err error) {
	symbols, err := SymbolOriginalService.GetAllList(false)
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}

	m := new(sync.Map)
	g, ctx := errgroup.WithContext(ctx)
	for _, symbol := range symbols {
		originalSymbol := symbol.Original
		g.Go(func() error {
			collName := s.getCollName(originalSymbol)
			ctm := carbon.Time2Carbon(dateTime).EndOfSecond().Carbon2Time().UnixMilli()
			var tick *model.Tick
			tick, err = repo.TickRepo.GetHistoryTick(ctx, collName, ctm)
			if err != nil {
				return errors.WithMessage(err, util.ErrWrapMessage())
			}
			if tick != nil {
				m.Store(originalSymbol, tick)
			}
			return nil
		})
	}

	if err = g.Wait(); err != nil {
		return nil, errors.Wrap(err, "g.Wait() is failed")
	}

	ticks = make(map[string]*model.Tick, len(symbols))
	m.Range(func(originalSymbol, tick any) bool {
		ticks[originalSymbol.(string)] = tick.(*model.Tick)
		return true
	})

	return
}

func (s *tickService) getCollName(symbol string) string {
	if strings.Contains(symbol, ".") {
		symbol = strings.Replace(symbol, ".", "_", -1)
	}
	return symbol
}
