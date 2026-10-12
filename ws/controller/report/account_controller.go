package report

import (
	"backend-rms/service"
	"backend-rms/util/formatter"
	"backend-rms/ws"
	"backend-rms/ws/controller"
	"backend-rms/ws/request"
	"backend-rms/ws/request/report"
	"backend-rms/ws/response"
	"time"
)

type AccountController struct{}

func (ac *AccountController) GetAccountList(c *ws.Context) {
	controller.TickerHandle(c, 3*time.Second, func(c *ws.Context) (data any, err error) {
		req := report.GetAccountListRequest{}
		if ok := request.Validate(c, &req); !ok {
			return
		}

		return ac.buildAccountList(req.Params.Logins)
	})
}

func (ac *AccountController) buildAccountList(logins []int64) ([]*response.ReportAccount, error) {
	groups, err := service.GroupService.GetGroupsByBook("ALL")
	if err != nil {
		return nil, err
	}

	accountList, err := service.AccountService.GetReportAccountList(groups, logins, nil)
	if err != nil {
		return nil, err
	}

	floatingPLMap, err := service.CacheService.GetAccountFloatingPLMap(logins)
	if err != nil {
		return nil, err
	}

	var accounts []*response.ReportAccount
	for _, account := range accountList {
		accounts = append(accounts, &response.ReportAccount{
			Login:   account.Login,
			Balance: formatter.Amount(account.Balance),
			Equity:  formatter.Amount(account.Balance.Add(floatingPLMap[account.Login])),
		})
	}

	return accounts, nil
}
