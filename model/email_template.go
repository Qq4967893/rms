package model

type EmailTemplate struct {
	BaseModel

	Title     string `json:"title"`
	Event     string `json:"event"`
	OnlyWorld string `json:"only_world"`
	Content   string `json:"content"`
	FrameId   int64  `json:"frame_id"`

	CommonTimestampsField
}
