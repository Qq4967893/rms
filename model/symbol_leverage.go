package model

type SymbolLeverages struct {
	BaseModel

	Symbol   string `json:"symbol"`
	Original string `json:"original"`
	Leverage int32  `json:"leverage"`

	CommonTimestampsField
}
