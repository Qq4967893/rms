package report

import "time"

type GetTradeListRequest struct {
	Book        string    `form:"book" validate:"required,oneof=all a b"`
	TradeStatus string    `form:"trade_status" validate:"required,oneof=all open closed holding"`
	Symbol      string    `form:"symbol" validate:"omitempty,min=4"`
	SearchType  string    `form:"search_type" validate:"omitempty,oneof=login passport ticket"`
	SearchText  int64     `form:"search_text" validate:"required_with=SearchType,omitempty,min=1"`
	StartTime   time.Time `form:"start_time" time_format:"2006-01-02 15:04:05" validate:"required_without_all=SearchType SearchText,ltefield=EndTime,tick_queryable_time"`
	EndTime     time.Time `form:"end_time" time_format:"2006-01-02 15:04:05" validate:"required_without_all=SearchType SearchText,gtefield=StartTime,tick_queryable_time"`
	Limit       int       `form:"limit" validate:"required,numeric,min=1,max=100"`
	Offset      int       `form:"offset" validate:"numeric,min=0"`
	Sort        string    `form:"sort" validate:"required,oneof=symbol lots open_time close_time"`
	Order       string    `form:"order" validate:"required,oneof=asc desc"`
}
