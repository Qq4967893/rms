package model

import "time"

type TradingLimit struct {
	BaseModel

	Passport     int64     `json:"passport"`
	Symbol       string    `json:"symbol"`
	TriggerValue int       `json:"trigger_value"`
	LimitValue   int       `json:"limit_value"`
	CurrentValue int       `json:"current_value"`
	RemoveValue  int       `json:"remove_value"`
	TriggerTime  time.Time `json:"trigger_time"`
	RemoveTime   time.Time `json:"remove_time"`
	Operator     string    `json:"operator"`

	CommonTimestampsField
}

type TradingLimitList struct {
	Id           int64     `json:"id"`
	Passport     int64     `json:"passport"`
	Symbol       string    `json:"symbol"`
	TriggerValue int       `json:"trigger_value"`
	LimitValue   int       `json:"limit_value"`
	CurrentValue int       `json:"current_value"`
	RemoveValue  int       `json:"remove_value"`
	TriggerTime  time.Time `json:"trigger_time"`
	RemoveTime   time.Time `json:"remove_time"`
	Operator     string    `json:"operator"`
}

type TradingLimitInfo struct {
	Passport     int64     `json:"passport"`
	Symbol       string    `json:"symbol"`
	LimitValue   int       `json:"limit_value"`
	CurrentValue int       `json:"current_value"`
	RemoveTime   time.Time `json:"remove_time"`
	Email        string    `json:"email"`
	Username     string    `json:"username"`
	RealName     string    `json:"real_name"`
	WorldCode    string    `json:"world_code"`
}

const TRADING_LIMIT = 1                           //限制交易
const REMOVE_TRADING_LIMIT = 0                    //解除限制交易
const TRADING_REMOVE_TIME = "0000-00-00 00:00:00" //解除限制的默认时间
