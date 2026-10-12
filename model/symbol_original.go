package model

type SymbolOriginal struct {
	BaseModel

	Original string `json:"original"`
	Digits   int64  `json:"digits"`

	CommonTimestampsField
}
