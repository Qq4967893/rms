package report

import (
	"backend-rms/http/request"
	"backend-rms/http/request/report"
	"backend-rms/http/response"
	"backend-rms/model"
	"backend-rms/service"
	"backend-rms/util"
	"backend-rms/util/formatter"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

type AccountController struct{}

func (ac *AccountController) GetAccountList(c *gin.Context) {
	req := report.GetAccountListRequest{}
	if ok := request.Validate(c, &req); !ok {
		return
	}

	req.Book = strings.ToUpper(req.Book)
	groups, err := service.GroupService.GetGroupsByBook(req.Book)
	if err != nil {
		response.Error(c, err)
		return
	}

	var login []int64
	switch req.SearchType {
	case "login":
		login = []int64{req.SearchText}
	case "passport":
		tradeAccountIDs, err := service.TradeAccountService.GetTradeAccountIDs(req.SearchText)
		if err != nil {
			response.Error(c, err)
			return
		}
		if len(tradeAccountIDs) == 0 {
			response.SucceedWithData(c, "", gin.H{
				"total_count": 0,
				"accounts":    nil,
			})
			return
		}
		login = tradeAccountIDs
	}

	g := new(errgroup.Group)
	var totalCount int64
	g.Go(func() error {
		totalCount, err = service.AccountService.GetReportAccountCount(groups, login)
		return errors.WithMessage(err, util.ErrWrapMessage())
	})

	var accountList []*service.ReportAccount
	g.Go(func() error {
		requestID := c.GetString("request_id")
		paging := model.Paging{Limit: req.Limit, Offset: req.Offset, Sort: req.Sort, Order: req.Order}
		accountList, err = service.ReportAccountService.GetReportAccountList(requestID, groups, login, req.StartTime, req.EndTime, paging)
		return errors.WithMessage(err, util.ErrWrapMessage())
	})

	if err = g.Wait(); err != nil {
		response.Error(c, err)
		return
	}

	accounts, err := ac.buildAccounts(accountList)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SucceedWithData(c, "", gin.H{
		"total_count": totalCount,
		"accounts":    accounts,
	})
}

func (ac *AccountController) buildAccounts(accountList []*service.ReportAccount) (accounts []response.ReportAccount, err error) {
	groupMap, err := service.GroupService.GetAllMap()
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}

	for _, account := range accountList {
		accounts = append(accounts, response.ReportAccount{
			Login:               account.Login,
			Book:                groupMap[account.Group],
			Passport:            account.Passport,
			CashInflow:          formatter.Amount(account.CashInflow),
			CashOutflow:         formatter.Amount(account.CashOutflow),
			NetCashFlow:         formatter.Amount(account.GetNetCashFlow()),
			Balance:             formatter.Amount(account.Balance),
			StartTimeBalance:    formatter.Amount(account.StartTimeBalance),
			EndTimeBalance:      formatter.Amount(account.EndTimeBalance),
			Equity:              formatter.Amount(account.GetEquity()),
			StartTimeEquity:     formatter.Amount(account.GetStartTimeEquity()),
			EndTimeEquity:       formatter.Amount(account.GetEndTimeEquity()),
			StartTimeFloatingPL: formatter.Amount(account.StartTimeFloatingPL),
			EndTimeFloatingPL:   formatter.Amount(account.EndTimeFloatingPL),
			FloatingPL:          formatter.Amount(account.FloatingPL),
			ClosedPL:            formatter.Amount(account.ClosedPL),
			GrossPL:             formatter.Amount(account.GetGrossPL()),
			PeriodClosedPL:      formatter.Amount(account.PeriodClosedPL),
			NetPeriodPL:         formatter.Amount(account.GetNetPeriodPL()),
		})
	}
	return
}
