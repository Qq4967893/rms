package service

import (
	"backend-rms/core/jwt"
	"backend-rms/core/redis"
	"backend-rms/model"
	"backend-rms/util"
	"backend-rms/util/crypto"
	"fmt"

	"github.com/pkg/errors"
)

type authService struct{}

var AuthService = &authService{}

func (as *authService) IsAccessTokenExisted(managerID int64, accessToken string) (bool, error) {
	n, err := redis.RMSRedisCli.Exists(
		as.getAccessTokenKey(managerID, crypto.Md5(accessToken)),
	)
	return n == 1, errors.WithMessage(err, util.ErrWrapMessage())
}

func (as *authService) IsRefreshTokenExisted(managerID int64, refreshToken string) (bool, error) {
	n, err := redis.RMSRedisCli.Exists(
		as.getRefreshTokenKey(managerID, crypto.Md5(refreshToken)),
	)
	return n == 1, errors.WithMessage(err, util.ErrWrapMessage())
}

func (as *authService) Login(username, password string) (manager *model.Manager, err error) {
	manager, err = ManagerService.GetInfoByUsername(username)
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}

	if manager == nil || manager.ID == 0 {
		err = errors.New("manager not found")
		return
	}

	if !crypto.PasswordVerify(password, manager.Password) {
		return manager, errors.New("wrong password")
	}

	return
}

func (as *authService) AddToken(managerID int64, accessToken, refreshToken string) error {
	if err := redis.RMSRedisCli.Set(
		as.getAccessTokenKey(managerID, crypto.Md5(accessToken)),
		crypto.Md5(refreshToken),
		jwt.Instance.ExpireTime,
	); err != nil {
		return errors.WithMessage(err, util.ErrWrapMessage())
	}

	if err := redis.RMSRedisCli.Set(
		as.getRefreshTokenKey(managerID, crypto.Md5(refreshToken)),
		"",
		jwt.Instance.MaxRefresh,
	); err != nil {
		return errors.WithMessage(err, util.ErrWrapMessage())
	}

	return nil
}

func (as *authService) UpdateToken(managerID int64, accessToken, refreshToken, oldAccessToken string) error {
	if err := redis.RMSRedisCli.Set(
		as.getAccessTokenKey(managerID, crypto.Md5(accessToken)),
		crypto.Md5(refreshToken),
		jwt.Instance.ExpireTime,
	); err != nil {
		return errors.WithMessage(err, util.ErrWrapMessage())
	}

	if _, err := redis.RMSRedisCli.Del(
		as.getAccessTokenKey(managerID, crypto.Md5(oldAccessToken)),
	); err != nil {
		return errors.WithMessage(err, util.ErrWrapMessage())
	}

	return nil
}

func (as *authService) DeleteToken(managerID int64, accessToken string) error {
	atKey := as.getAccessTokenKey(managerID, crypto.Md5(accessToken))
	rtMd5, err := redis.RMSRedisCli.Get(atKey)
	if err != nil {
		return errors.WithMessage(err, util.ErrWrapMessage())
	}

	rtKey := as.getRefreshTokenKey(managerID, rtMd5)
	_, err = redis.RMSRedisCli.Del(atKey, rtKey)
	if err != nil {
		return errors.WithMessage(err, util.ErrWrapMessage())
	}

	return nil
}

func (as *authService) getAccessTokenKey(managerID int64, accessTokenMd5 string) string {
	return fmt.Sprintf("access_token:%d:%s", managerID, accessTokenMd5)
}

func (as *authService) getRefreshTokenKey(managerID int64, refreshTokenMd5 string) string {
	return fmt.Sprintf("refresh_token:%d:%s", managerID, refreshTokenMd5)
}
