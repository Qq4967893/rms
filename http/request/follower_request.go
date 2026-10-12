package request

import "time"

type GetFollowerListRequest struct {
	MasterId      int64     `uri:"masterId" validate:"required,numeric,min=1"`
	StartTime     time.Time `form:"start_time" time_format:"2006-01-02 15:04:05" validate:"required,ltefield=EndTime"`
	EndTime       time.Time `form:"end_time" time_format:"2006-01-02 15:04:05" validate:"required,gtefield=StartTime"`
	Book          string    `form:"book" validate:"required,oneof=all a b"`
	FollowThanOne int       `form:"follow_than_one" validate:"numeric"` //是否跟随超过一个高手（-1:全部，0:否，1:超过一个高手）
	IsFollowing   int       `form:"is_following" validate:"numeric"`    //是否正在跟随(-1:全部，0:否，1:正在跟随)
	Limit         int       `form:"limit" validate:"required,numeric,min=1,max=100"`
	Offset        int       `form:"offset" validate:"numeric,min=0"`
	Sort          string    `form:"sort" validate:"required,oneof=period_pl start_time closed_pl floating_pl balance margin_level follow_master_num copied_amount copied_all_amount"`
	Order         string    `form:"order" validate:"required,oneof=asc desc"`
}

type GetFollowerHistoryListRequest struct {
	MasterId   int64  `uri:"masterId" validate:"required,numeric,min=1"`
	FollowerId int64  `uri:"followerId" validate:"required,numeric,min=1"`
	Limit      int    `form:"limit" validate:"required,numeric,min=1,max=100"`
	Offset     int    `form:"offset" validate:"numeric,min=0"`
	Sort       string `form:"sort" validate:"required,oneof=start_time"`
	Order      string `form:"order" validate:"required,oneof=asc desc"`
}
