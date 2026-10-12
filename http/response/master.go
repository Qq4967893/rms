package response

type Master struct {
	Login         int64  `json:"login"`
	Passport      int64  `json:"passport"`
	PeriodPL      string `json:"period_pl"`
	ClosedPL      string `json:"closed_pl"`
	FloatingPL    string `json:"floating_pl"`
	Balance       string `json:"balance"`
	FollowerCount int64  `json:"follower_count"`
	CopiedAmount  string `json:"copied_amount"`
	MarginLevel   string `json:"margin_level"`
	Book          string `json:"book"`
	IsPresent     string `json:"is_present"`
}

type MasterHistoricalRecord struct {
	Login         int64  `json:"login"`
	Commission    string `json:"commission"`
	ClosedPL      string `json:"closed_pl"`
	FollowerCount int64  `json:"follower_count"`
	FollowedTimes int64  `json:"followed_times"`
	StartTime     string `json:"start_time"`
	EndTime       string `json:"end_time"`
}
