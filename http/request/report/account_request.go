package report

import "time"

type GetAccountListRequest struct {
	Book       string    `form:"book" validate:"required,oneof=all a b"`
	SearchType string    `form:"search_type" validate:"omitempty,oneof=login passport"`
	SearchText int64     `form:"search_text" validate:"required_with=SearchType,omitempty,min=1"`
	StartTime  time.Time `form:"start_time" time_format:"2006-01-02 15:04:05" validate:"required_without_all=SearchType SearchText,ltefield=EndTime,tick_queryable_time"`
	EndTime    time.Time `form:"end_time" time_format:"2006-01-02 15:04:05" validate:"required_without_all=SearchType SearchText,gtefield=StartTime,tick_queryable_time"`
	Limit      int       `form:"limit" validate:"required,numeric,min=1,max=100"`
	Offset     int       `form:"offset" validate:"numeric,min=0"`
	Sort       string    `form:"sort" validate:"required,oneof=reg_time cash_inflow cash_outflow net_cash_flow start_time_balance end_time_balance start_time_equity end_time_equity start_time_floating_pl end_time_floating_pl floating_pl closed_pl gross_pl period_closed_pl net_period_pl"`
	Order      string    `form:"order" validate:"required,oneof=asc desc"`
}
