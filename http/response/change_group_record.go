package response

type ChangeGroupRecordList struct {
	Login      int64  `json:"login"`
	Operate    string `json:"operate"`
	GroupOld   string `json:"group_old"`
	TransferTo string `json:"transfer_to"`
	Hedge      string `json:"hedge"`
	Time       string `json:"time"`
	Operator   string `json:"operator"`
}
