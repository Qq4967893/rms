package model

import (
	"time"

	"github.com/golang-module/carbon/v2"
	"github.com/shopspring/decimal"
)

type ReportDaily struct {
	BaseModel

	Login   int64           `json:"login"`
	Time    string          `json:"time"`
	Balance decimal.Decimal `json:"balance"`
	Equity  decimal.Decimal `json:"equity"`
	Group   string          `json:"group"`
	UsdRate decimal.Decimal `json:"usd_rate"`
	GbpRate decimal.Decimal `json:"gbp_rate"`
}

func GetDailyDateTime(dateTime time.Time) carbon.Carbon {
	return carbon.Time2Carbon(dateTime).SetHour(23).SetMinute(59).SetSecond(59)
}

func (ReportDaily) TableName() string {
	return "report_daily"
}
