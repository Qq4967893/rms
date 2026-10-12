package response

type Passport struct {
	Passport    int64  `json:"passport"`
	Username    string `json:"username"`
	RegTime     string `json:"reg_time"`
	ManagerName string `json:"manager_name"`
	Deposit     string `json:"deposit"`
	Withdrawal  string `json:"withdrawal"`
	NetDeposit  string `json:"net_deposit"`
	Balance     string `json:"balance"`
}
