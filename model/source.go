package model

import (
	"github.com/shopspring/decimal"
)

type Source struct {
	Source         string          `json:"source"`
	SourceType     int             `json:"source_type"`
	Digits         int             `json:"digits"`
	ContractSize   decimal.Decimal `json:"contract_size"`
	ProfitCurrency string          `json:"profit_currency"`
	Currency       string          `json:"currency"`
	MarginCurrency string          `json:"margin_currency"`
	Status         int             `json:"status"`
}

type SourceType struct {
	Source     string `json:"source"`
	SourceType int    `json:"source_type"`
}

var SourceTypeMap = map[int]string{
	0: "currency",
	1: "metal",
	2: "metal",
	3: "energy",
	4: "energy",
	5: "indices",
	6: "crypto",
	7: "stock",
}

const (
	SourceFX        = iota // 外汇
	SourceGold             // 黄金
	SourceSliver           // 白银
	SourceGas              // 天然气
	SourceOil              // 原油
	SourceIndex            // 股指
	SourceCrypto           // 虚拟货币
	SourceEquityCFD        // 股票
)

var Currency = map[string]map[string]string{
	//外汇货币
	"currency": {
		"AUD": "AUDUSD",
		"CAD": "USDCAD",
		"CHF": "USDCHF",
		"CNH": "USDCNH",
		"EUR": "EURUSD",
		"GBP": "GBPUSD",
		"HKD": "USDHKD",
		"JPY": "USDJPY",
		"MXN": "USDMXN",
		"NZD": "NZDUSD",
		"PLN": "USDPLN",
		"SGD": "USDSGD",
		"TRY": "USDTRY",
		"NOK": "USDNOK",
		"SEK": "USDSEK",
		"USD": "", //单独处理美元
	},
	//贵金属
	"commodities": {
		"XAG": "XAGUSD",
		"XAU": "XAUUSD",
		"XBR": "XBRUSD",
		"XNG": "XNGUSD",
		"XTI": "XTIUSD",
	},
	//虚拟货币
	"crypto": {
		"BTC": "BTCUSD",
		"ETH": "ETHUSD",
		"DOT": "DOTUSD",
		"XRP": "XRPUSD",
		"LTC": "LTCUSD",
		"BCH": "BCHUSD",
		"XLM": "XLMUSD",
		"EOS": "EOSUSD",
		"LNK": "LNKUSD",
		"AUS": "AUSD", //EOSUSD品种改名字为AUSD
	},
	//股指
	"indices": {
		"AUS200":  "AUDUSD",
		"EUSTX50": "EURUSD",
		"GER30":   "EURUSD",
		"HK50":    "USDHKD",
		"JPN225":  "USDJPY",
		"NAS100":  "",
		"SPA35":   "EURUSD",
		"UK100":   "GBPUSD",
		"US30":    "",
		"USA500":  "",
	},
	//股票
	"stock": {
		"BABA.ny": "",
		"DIS.ny":  "",
		"VISA.ny": "",
		"AAPL.nd": "",
		"AMD.nd":  "",
		"AMZN.nd": "",
		"EBAY.nd": "",
		"META.nd": "",
		"GOOG.nd": "",
		"NFLX.nd": "",
		"TSLA.nd": "",
		"PYPL.nd": "",
		"MRNA.nd": "",
		"MSFT.nd": "",
		"KO.ny":   "",
		"C.ny":    "",
		"PFE.ny":  "",
		"JPM.ny":  "",
		"COIN.nd": "",
		"FCX.ny":  "",
	},
}
