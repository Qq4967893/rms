package report

import "time"

type GetIncomeSymbolListRequest struct {
	Book      string    `form:"book" validate:"required,oneof=all a b"`
	Symbol    string    `form:"symbol"`
	StartDate time.Time `form:"start_date" time_format:"2006-01-02" validate:"required,ltefield=EndDate,tick_queryable_time"`
	EndDate   time.Time `form:"end_date" time_format:"2006-01-02" validate:"required,gtefield=StartDate,tick_queryable_time"`
	Export    string    `form:"export" validate:"required,oneof=yes no"`
	Sort      string    `form:"sort" validate:"required,oneof=symbol profit swaps commission mark_up_spread_profit raw_spread_profit lots"`
	Order     string    `form:"order" validate:"required,oneof=asc desc"`
}

type GetIncomeAccountListRequest struct {
	Book      string    `form:"book" validate:"required,oneof=all a b"`
	Symbol    string    `form:"symbol"`
	StartDate time.Time `form:"start_date" time_format:"2006-01-02" validate:"required,ltefield=EndDate,tick_queryable_time"`
	EndDate   time.Time `form:"end_date" time_format:"2006-01-02" validate:"required,gtefield=StartDate,tick_queryable_time"`
	Export    string    `form:"export" validate:"required,oneof=yes no"`
	Limit     int       `form:"limit" validate:"required,numeric,min=1,max=100"`
	Offset    int       `form:"offset" validate:"numeric,min=0"`
	Sort      string    `form:"sort" validate:"required,oneof=login profit swaps commission mark_up_spread_profit raw_spread_profit lots"`
	Order     string    `form:"order" validate:"required,oneof=asc desc"`
}

type GetIncomeTradeListRequest struct {
	Book      string    `form:"book" validate:"required,oneof=all a b"`
	Symbol    string    `form:"symbol"`
	StartDate time.Time `form:"start_date" time_format:"2006-01-02" validate:"required,ltefield=EndDate,tick_queryable_time"`
	EndDate   time.Time `form:"end_date" time_format:"2006-01-02" validate:"required,gtefield=StartDate,tick_queryable_time"`
	Export    string    `form:"export" validate:"required,oneof=yes no"`
	Limit     int       `form:"limit" validate:"required,numeric,min=1,max=100"`
	Offset    int       `form:"offset" validate:"numeric,min=0"`
	Sort      string    `form:"sort" validate:"required,oneof=login ticket profit swaps commission mark_up_spread_profit raw_spread_profit lots open_time close_time"`
	Order     string    `form:"order" validate:"required,oneof=asc desc"`
}
