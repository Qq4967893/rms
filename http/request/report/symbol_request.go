package report

import "time"

type GetSymbolListRequest struct {
	Book      string    `form:"book" validate:"required,oneof=a b all"`
	Symbol    string    `form:"symbol"`
	StartTime time.Time `form:"start_time" time_format:"2006-01-02 15:04:05" validate:"required,ltefield=EndTime,tick_queryable_time"`
	EndTime   time.Time `form:"end_time" time_format:"2006-01-02 15:04:05" validate:"required,gtefield=StartTime,tick_queryable_time"`
	Sort      string    `form:"sort" validate:"required,oneof=closed_pl closed_profit closed_swap closed_commission floating_pl floating_profit floating_swap floating_commission period_pl"`
	Order     string    `form:"order" validate:"required,oneof=asc desc"`
}

type GetSymbolAccountListRequest struct {
	Symbol    string    `form:"symbol" validate:"required,min=4"`
	Book      string    `form:"book" validate:"required,oneof=a b all"`
	StartTime time.Time `form:"start_time" time_format:"2006-01-02 15:04:05" validate:"required,ltefield=EndTime,tick_queryable_time"`
	EndTime   time.Time `form:"end_time" time_format:"2006-01-02 15:04:05" validate:"required,gtefield=StartTime,tick_queryable_time"`
	Sort      string    `form:"sort" validate:"required,oneof=closed_pl start_time_floating_pl end_time_floating_pl period_pl period_floating_pl"`
	Order     string    `form:"order" validate:"required,oneof=asc desc"`
	Limit     int       `form:"limit" validate:"required,numeric,min=1,max=100"`
	Offset    int       `form:"offset" validate:"numeric,min=0"`
}
