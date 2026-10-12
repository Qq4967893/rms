package route

import (
	"backend-rms/ws"
	"backend-rms/ws/controller/market"
	"backend-rms/ws/controller/report"
)

func RegisterRoutes(r *ws.Engine) {
	reportAccountController := new(report.AccountController)
	{
		r.GET("/api/ws/report/accounts", reportAccountController.GetAccountList)
	}

	reportTradeController := new(report.TradeController)
	{
		r.GET("/api/ws/report/trades", reportTradeController.GetTradeList)
	}

	marketSummaryController := new(market.SummaryController)
	{
		r.GET("/api/ws/market/summaries", marketSummaryController.GetSummaryList)
	}

	marketExposureController := new(market.ExposureController)
	{
		r.GET("/api/ws/market/exposures", marketExposureController.GetExposureList)
	}

	marketStockController := new(market.StockController)
	{
		r.GET("/api/ws/market/stocks/summaries", marketStockController.GetStockList)
	}
}
