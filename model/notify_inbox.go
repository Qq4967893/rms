package model

type NotifyInbox struct {
	BaseModel

	UserCode int64  `json:"user_code"`
	Mid      int64  `json:"mid"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	Format   int    `json:"format"`
	Type     int    `json:"type"`
	Resource string `json:"resource"`

	CommonTimestampsField
}
