package calc_profit

import (
	"backend-rms/core/logger"
	"backend-rms/core/redis"
	"backend-rms/grpc/tsp"
	"backend-rms/model"
	"backend-rms/service"
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/shopspring/decimal"
)

func Start() {
	ticker1 := time.NewTicker(2 * time.Second)
	defer ticker1.Stop()

	spreadInfo, err := tsp.GetSpreadInfo(context.Background(), "test request id")
	if err != nil {
		logger.ErrorIf(err)
		return
	}

	groups, err := service.GroupService.GetGroupsByBook("ALL")
	if err != nil {
		logger.ErrorIf(err)
		return
	}

	for {
		<-ticker1.C
		fmt.Println("--------------------------Start-------------------------------")

		openingTrades, err := service.TradeService.GetOpeningTrades(groups, 0)
		if err != nil {
			logger.ErrorIf(err)
			return
		}

		//暂无持仓单数据
		if len(openingTrades) < 1 {
			return
		}

		allQuotes, err := service.QuoteService.GetAllQuotes()
		if err != nil {
			logger.ErrorIf(err)
			return
		}

		calculatedProfits, err := service.ProfitService.CalcProfitBatch(openingTrades, spreadInfo, allQuotes)
		if err != nil {
			logger.ErrorIf(err)
			return
		}
		cache(openingTrades, calculatedProfits)
		fmt.Println("--------------------------End-------------------------------")
	}
}

func cache(openingTrades []model.OpeningTrade, calculatedProfits map[int64]service.CalculatedProfit) {
	openingTradeMap := make(map[int64]model.OpeningTrade)
	for _, ot := range openingTrades {
		openingTradeMap[ot.Ticket] = ot
	}

	tradePLs := make(map[string]string, len(calculatedProfits))
	clientPLs := make(map[string]string)
	for ticket, cp := range calculatedProfits {
		swaps := openingTradeMap[ticket].Swaps
		commission := openingTradeMap[ticket].Commission
		profit := cp.GrossProfit.Add(swaps).Add(commission)

		ticketStr := strconv.FormatInt(ticket, 10)
		tradePLs[ticketStr] = profit.String()

		login := strconv.FormatInt(openingTradeMap[ticket].Login, 10)
		totalPL, _ := decimal.NewFromString(clientPLs[login])
		clientPLs[login] = profit.Add(totalPL).String()
	}

	redis.RMSRedisCli.Del(service.TradePLKey, service.ClientPLKey)
	res, err := redis.RMSRedisCli.HSet(service.TradePLKey, tradePLs)
	if err != nil {
		logger.ErrorIf(err)
	}
	fmt.Printf("%d records of trade P&L have Saved\n", res)

	res, err = redis.RMSRedisCli.HSet(service.ClientPLKey, clientPLs)
	if err != nil {
		logger.ErrorIf(err)
	}
	fmt.Printf("%d records of client P&L have Saved\n", res)

	redis.RMSRedisCli.Expire(service.TradePLKey, 10*time.Second)
	redis.RMSRedisCli.Expire(service.ClientPLKey, 10*time.Second)
}
