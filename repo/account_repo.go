package repo

import (
	"backend-rms/core/database"
	"backend-rms/model"
	"backend-rms/util"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type accountRepo struct{}

var AccountRepo = &accountRepo{}

func (r *accountRepo) GetAll(groups []string) (accounts []*model.Account, err error) {
	tx := database.ET6DB.Model(model.Account{})
	if len(groups) > 0 {
		tx.Where("`group` IN ?", groups)
	}

	err = errors.Wrapf(tx.Find(&accounts).Error, "%s failed", util.GetCurrentFuncName())
	return
}

func (r *accountRepo) GetInfoByLogin(login int64) (account *model.Account, err error) {
	err = database.ET6DB.Table("account as a").
		Joins("left join rms.groups as g on a.group = g.name").
		Where("a.login", login).
		Select("a.login, a.group, a.currency, a.balance, a.regdate").First(&account).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return account, errors.Wrapf(err, "%s failed", util.GetCurrentFuncName())
}

func (r *accountRepo) GetListByLogins(logins []int64, groups []string) (accounts []*model.Account, err error) {
	tx := database.ET6DB.Where("login IN ?", logins)
	if len(groups) > 0 {
		tx.Where("`group` IN ?", groups)
	}

	err = errors.Wrapf(tx.Find(&accounts).Error, "%s failed", util.GetCurrentFuncName())
	return
}

func (r *accountRepo) GetReportAccountList(groups []string, logins []int64, paging *model.Paging) (accounts []*model.ReportAccount, err error) {
	tx := database.ET6DB.Model(model.Account{})
	if len(groups) > 0 {
		tx.Where("`group` IN ?", groups)
	}

	if len(logins) > 0 {
		tx.Where("login IN ?", logins)
	}

	if paging != nil {
		tx.Order(paging.Sort + " " + paging.Order)
		if paging.Sort != "login" {
			tx.Order("login desc")
		}
		tx.Limit(paging.Limit).Offset(paging.Offset)
	}

	err = errors.Wrapf(tx.Find(&accounts).Error, "%s failed", util.GetCurrentFuncName())
	return
}

func (r *accountRepo) GetReportAccountCount(groups []string, login []int64) (totalCount int64, err error) {
	tx := database.ET6DB.Model(model.Account{})
	if len(groups) > 0 {
		tx.Where("`group` IN ?", groups)
	}

	if len(login) > 0 {
		if len(login) == 1 {
			tx.Where("login = ?", login)
		} else {
			tx.Where("`login` IN ?", login)
		}
	}

	result := tx.Count(&totalCount)
	return totalCount, errors.Wrapf(result.Error, "%s failed", util.GetCurrentFuncName())
}
