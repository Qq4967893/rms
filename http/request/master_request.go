package request

import "time"

type GetMasterListRequest struct {
	Book       string    `form:"book" validate:"required,oneof=all a b"`
	IsPresent  string    `form:"is_present" validate:"omitempty,oneof=yes no"`
	SearchType string    `form:"search_type" validate:"omitempty,oneof=login"`
	SearchText int64     `form:"search_text" validate:"required_with=SearchType,omitempty,min=1"`
	StartTime  time.Time `form:"start_time" time_format:"2006-01-02 15:04:05" validate:"required,ltefield=EndTime"`
	EndTime    time.Time `form:"end_time" time_format:"2006-01-02 15:04:05" validate:"required,gtefield=StartTime"`
	Limit      int       `form:"limit" validate:"required,numeric,min=1,max=100"`
	Offset     int       `form:"offset" validate:"numeric,min=0"`
	Sort       string    `form:"sort" validate:"required,oneof=period_pl closed_pl floating_pl follower_count copied_amount margin_level"`
	Order      string    `form:"order" validate:"required,oneof=asc desc"`
}

type DownloadMasterListRequest struct {
	Book       string    `form:"book" validate:"required,oneof=all a b"`
	IsPresent  string    `form:"is_present" validate:"omitempty,oneof=yes no"`
	SearchType string    `form:"search_type" validate:"omitempty,oneof=login"`
	SearchText int64     `form:"search_text" validate:"required_with=SearchType,omitempty,min=1"`
	StartTime  time.Time `form:"start_time" time_format:"2006-01-02 15:04:05" validate:"required,ltefield=EndTime"`
	EndTime    time.Time `form:"end_time" time_format:"2006-01-02 15:04:05" validate:"required,gtefield=StartTime"`
}

type GetMasterHistoryListRequest struct {
	MasterId int64  `uri:"masterId" validate:"required,numeric,min=1"`
	Limit    int    `form:"limit" validate:"required,numeric,min=1,max=100"`
	Offset   int    `form:"offset" validate:"numeric,min=0"`
	Sort     string `form:"sort" validate:"required,oneof=start_time"`
	Order    string `form:"order" validate:"required,oneof=asc desc"`
}
