package response

type GetRestrictedPassportList struct {
	Id           int64  `json:"id"`
	Passport     int64  `json:"passport"`
	Symbol       string `json:"symbol"`
	TriggerValue int    `json:"trigger_value"`
	LimitValue   int    `json:"limit_value"`
	CurrentValue int    `json:"current_value"`
	TriggerTime  string `json:"trigger_time"`
}

type GetRestrictedPassportRecord struct {
	Passport     int64  `json:"passport"`
	Symbol       string `json:"symbol"`
	TriggerValue int    `json:"trigger_value"`
	LimitValue   int    `json:"limit_value"`
	RemoveValue  int    `json:"remove_value"`
	TriggerTime  string `json:"trigger_time"`
	RemoveTime   string `json:"remove_time"`
	Operator     string `json:"operator"`
}
