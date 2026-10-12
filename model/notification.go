package model

type Notification struct {
	BaseModel

	Title    string `json:"title"`
	Content  string `json:"content"`
	Format   int    `json:"format"`
	Type     int    `json:"type"`
	Resource string `json:"resource"`
	Partner  string `json:"partner"`
}
