package model

type SymbolType struct {
	BaseModel

	Original   string `json:"original"`
	Identifier string `json:"identifier"`

	CommonTimestampsField
}

const (
	SymbolIdentifierCurrency = "currency" // 外汇
	SymbolIdentifierMetal    = "metal"    // 贵金属（包含黄金，白银）
	SymbolIdentifierEnergy   = "energy"   // 能源（包含原油，天然气）
	SymbolIdentifierCrypto   = "crypto"   // 虚拟货币
	SymbolIdentifierStock    = "stock"    // 股票
	SymbolIdentifierIndices  = "indices"  // 股指
)

// source表内source_type字段对应的redis内exposure hash key的field
var SymbolIdentifierRedisExposure = map[string]string{
	"currency": "currency",
	"metal":    "commodities",
	"energy":   "commodities",
	"indices":  "indices",
	"crypto":   "crypto",
	"stock":    "stock",
}
