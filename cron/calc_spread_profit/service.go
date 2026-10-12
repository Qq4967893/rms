package calc_spread_profit

import (
	"backend-rms/core/logger"
	"backend-rms/core/redis"
	"backend-rms/model"
	"backend-rms/service"
	"backend-rms/util/formatter"
	"time"

	"github.com/golang-module/carbon/v2"
	cron "github.com/robfig/cron/v3"
)

func Start() {
	c := cron.New()

	c.AddFunc("01 * * * *", func() {

		lockName := "cron:calc_spread_profit"
		//添加分布式锁，保证多台机器的cron任务只有一台机器是执行的。
		//设置的过期事件小于cron执行的间隔时间
		setNX, err := redis.RMSRedisCli.SetNX(lockName, true, 3590*time.Second)
		if err != nil {
			logger.Errorf("cron redis lock setNX %v is error: %v\n", lockName, err)
		}

		if !setNX {
			return
		}

		startTimeStamp := time.Now().Unix()

		startTimeOfCarbon := carbon.Now(carbon.UTC).SubHours(2).EndOfHour().ToDateTimeString()
		endTimeOfCarbon := carbon.Now(carbon.UTC).SubHours(1).EndOfHour().ToDateTimeString()

		startTime, _ := time.Parse(formatter.DateTime, startTimeOfCarbon)
		endTime, _ := time.Parse(formatter.DateTime, endTimeOfCarbon)

		groups, err := service.GroupService.GetGroupsByBook("all")
		if err != nil {
			logger.Errorf("cron calc spread profit get groups is error: %v\n", err)
			return
		}

		Handle(startTime, endTime, groups)

		handleTime := time.Now().Unix() - startTimeStamp
		logger.Infof("cal spread profit time is %v s\n", handleTime)
	})

	c.Start()
}

func Handle(startTime, endTime time.Time, groups []string) {
	trades, err := service.TradeService.GetClosedTrades(groups, startTime, endTime)
	if err != nil {
		logger.Errorf("cron calc spread profit get closed trades is error: %v\n", err)
		return
	}

	if len(trades) == 0 {
		return
	}

	tradeSpreadProfits, err := service.SpreadProfitService.CalcSpreadProfit(trades)
	if err != nil {
		logger.Errorf("cron calc spread profit get spread profit is error: %v\n", err)
		return
	}

	var data []model.DailyClosedTrade
	for ticket, spreadProfit := range tradeSpreadProfits {
		closeDate, err := time.Parse(formatter.DateOnly, spreadProfit.CloseDate.Format(formatter.DateOnly))
		if err != nil {
			logger.Errorf("cron calc spread profit change closeDate style is error: %v\n", err)
			return
		}
		middle := model.DailyClosedTrade{
			Ticket:                 ticket,
			CloseDate:              closeDate,
			SpreadProfit:           spreadProfit.SpreadProfit.Round(2),
			WithoutRawSpreadProfit: spreadProfit.WithoutRawSpreadProfit.Round(2),
		}
		data = append(data, middle)
	}

	err = service.DailyClosedTradeService.Add(data)
	if err != nil {
		logger.Errorf("cron calc spread profit add data to table is error: %v\n", err)
		return
	}
	return
}
