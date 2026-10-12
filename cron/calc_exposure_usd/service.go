package calc_exposure_usd

import (
	"backend-rms/core/logger"
	"backend-rms/core/redis"
	"backend-rms/service"
	"backend-rms/tr/init_data"
	"strconv"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/shopspring/decimal"
	"golang.org/x/sync/errgroup"
)

func Start() {
	c := cron.New()
	_, err := c.AddFunc("@every 1s", handleExposureUsd)
	if err != nil {
		logger.Errorf("Calc exposure usd cron is err:%v", err)
		return
	}
	c.Start()
}

func handleExposureUsd() {
	handlerExposureStartTime := time.Now().UnixMilli()
	g := new(errgroup.Group)
	for _, v := range init_data.Books {
		book := v
		g.Go(func() error {
			err := CalculateExposureUsd(book)
			return err
		})
	}

	if err := g.Wait(); err != nil {
		logger.Error("The cron is error:", err)
	}
	logger.Info("cron执行所需时间：" + strconv.Itoa(int(time.Now().UnixMilli()-handlerExposureStartTime)) + "ms")
}

func CalculateExposureUsd(book string) (err error) {
	exposureBook := make(map[string][]*service.InitExposure)

	//获取品种的实时报价
	quotes, err := redis.QuoteRedisCli.HGetAll("et6:tick")
	if err != nil {
		logger.Error("Redis hgetall quotes is error:", err)
		return
	}

	exposureTypes := service.SourceService.GetCurrencyExposureTypes()
	for currency, exposureInfo := range service.ExposureBook[book] {
		exposureType := exposureTypes[currency]

		//根据exposure计算exposure USD
		// USD的 exposure_USD = exposure
		if currency == "USD" {
			exposureInfo.ExposureUsd = exposureInfo.Exposure
		}

		// 非USD的 exposure_USD需要exposure 通过品种报价做转换
		// 判断如果是股指或股票需要先乘以自身的报价，再转换为exposure_USD
		var processedExposure decimal.Decimal
		if exposureType == "indices" || exposureType == "stock" {
			if _, ok := quotes[exposureInfo.Currency]; !ok {
				logger.Errorf("the %s quote is not exist", exposureInfo.Currency)
				continue
			}
			quoteInfo := strings.Split(quotes[exposureInfo.Currency], ",")
			ask, _ := decimal.NewFromString(quoteInfo[0])
			bid, _ := decimal.NewFromString(quoteInfo[1])
			middleQuote := (ask.Add(bid)).Div(decimal.NewFromInt(2))
			processedExposure = exposureInfo.Exposure.Mul(middleQuote)
		} else {
			processedExposure = exposureInfo.Exposure
		}

		var exposureUsd decimal.Decimal
		if exposureInfo.Symbol == "" {
			exposureUsd = processedExposure.Round(2)
		} else {
			quoteInfo := strings.Split(quotes[exposureInfo.Symbol], ",")
			ask, _ := decimal.NewFromString(quoteInfo[0])
			bid, _ := decimal.NewFromString(quoteInfo[1])
			middleQuote := (ask.Add(bid)).Div(decimal.NewFromInt(2))

			// 判断转换USD的品种报价名称，如果前三位是USD除以，如果是后三位乘以
			if exposureInfo.Symbol[0:3] == "USD" {
				exposureUsd = processedExposure.Div(middleQuote).Round(2)
			} else {
				exposureUsd = processedExposure.Mul(middleQuote).Round(2)
			}
		}
		exposureInfo.ExposureUsd = exposureUsd

		//判断该exposure的类型是否存在
		if _, ok := exposureBook[exposureType]; ok {
			exposureBook[exposureType] = append(exposureBook[exposureType], exposureInfo)
			continue
		}
		exposureBook[exposureType] = []*service.InitExposure{
			exposureInfo,
		}
	}
	service.CacheService.SetExposureBook(book, exposureBook)
	return
}
