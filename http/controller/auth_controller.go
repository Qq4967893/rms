package controller

import "C"
import (
	"backend-rms/core/jwt"
	"backend-rms/core/redis"
	"backend-rms/http/request"
	"backend-rms/http/response"
	"backend-rms/service"
	"backend-rms/util/crypto"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mojocn/base64Captcha"
	"github.com/pkg/errors"
)

type AuthController struct{}

func (ac *AuthController) GetCode(c *gin.Context) {
	// 1分钟只能刷新5次
	ip := c.ClientIP()
	countCodeKey := "count-code-" + ip
	countCodeV, err := redis.RMSRedisCli.Get(countCodeKey)
	if err != nil {
		response.Error(c, err)
		return
	}

	var count int
	if countCodeV == "" {
		count = 0
	} else {
		count, err = strconv.Atoi(countCodeV)
		if err != nil {
			response.Error(c, err)
			return
		}
	}
	if count >= 5 {
		response.FailWithMessage(c, "Too Many Attempts.")
		return
	}

	if count == 0 {
		redis.RMSRedisCli.Set(countCodeKey, strconv.Itoa(count+1), time.Second*60)
	} else {
		countCodeTime, err := redis.RMSRedisCli.TTL(countCodeKey)
		if err != nil {
			response.Error(c, err)
			return
		}
		redis.RMSRedisCli.Set(countCodeKey, strconv.Itoa(count+1), countCodeTime)
	}

	driver := base64Captcha.NewDriverDigit(80, 240, 4, 0.5, 80)
	store := base64Captcha.DefaultMemStore
	captcha := base64Captcha.NewCaptcha(driver, store)
	// 生成验证码
	codeId, b64s, code, err := captcha.Generate()
	if err != nil {
		response.Error(c, err)
		return
	}

	//存储验证码
	err = redis.RMSRedisCli.Set("auth-code-id-"+codeId, code, time.Minute*5)
	if err != nil {
		response.Error(c, err)
		return
	}

	// 返回验证码图片
	response.SucceedWithData(c, "", gin.H{
		"code_id": codeId,
		"image":   b64s,
		"expires": time.Now().Add(time.Minute * 5).Unix(),
	})
}

func (ac *AuthController) Login(c *gin.Context) {
	req := request.LoginRequest{}
	if ok := request.Validate(c, &req); !ok {
		return
	}

	authCode, err := service.CacheService.ValidateAuthCode(req.CodeId, req.Code)
	if err != nil {
		response.Error(c, err)
		return
	}

	if !authCode {
		response.FailWithMessage(c, "Wrong code.")
		return
	}

	manager, err := service.AuthService.Login(req.Username, req.Password)
	if err != nil {
		response.FailWithMessage(c, "Wrong username or password.")
		return
	}

	c.Set("login_manager_id", manager.ID)
	accessToken, expireTime, err := jwt.Instance.CreateToken(manager.ID)
	if err != nil {
		response.Error(c, err)
		return
	}

	refreshToken := crypto.Md5(uuid.New().String())
	if err = service.AuthService.AddToken(manager.ID, accessToken, refreshToken); err != nil {
		response.Error(c, err)
		return
	}

	response.SucceedWithData(c, "Login success", gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"expire_time":   expireTime.UnixMilli(),
	})
}

func (ac *AuthController) RefreshToken(c *gin.Context) {
	req := request.RefreshTokenRequest{}
	if ok := request.Validate(c, &req); !ok {
		return
	}

	accessToken, expireTime, claims, err := jwt.Instance.RefreshToken(c.Request)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpiredMaxRefresh) {
			response.Abort(c, response.StatusUnauthorized, response.CodeRefreshTokenExpired,
				response.WithMessage("Token is expired."),
			)
			return
		} else {
			response.Error(c, err)
			return
		}
	}

	c.Set("login_manager_id", claims.LoginManagerID)
	isRefreshTokenExisted, err := service.AuthService.IsRefreshTokenExisted(claims.LoginManagerID, req.RefreshToken)
	if err != nil {
		response.Error(c, err)
		return
	}

	if !isRefreshTokenExisted {
		response.Abort(c, response.StatusUnauthorized, response.CodeRefreshTokenExpired,
			response.WithMessage("Token is expired."),
		)
		return
	}

	oldAccessToken, err := jwt.Instance.GetTokenFromHeader(c.Request)
	if err != nil {
		response.Error(c, err)
		return
	}

	if err = service.AuthService.UpdateToken(claims.LoginManagerID, accessToken, req.RefreshToken, oldAccessToken); err != nil {
		response.Error(c, err)
		return
	}

	response.SucceedWithData(c, "Refresh token success", gin.H{
		"access_token": accessToken,
		"expire_time":  expireTime.UnixMilli(),
	})
}

func (ac *AuthController) Logout(c *gin.Context) {
	accessToken, err := jwt.Instance.GetTokenFromHeader(c.Request)
	if err != nil {
		response.Error(c, err)
		return
	}

	managerID := c.GetInt64("login_manager_id")
	if err = service.AuthService.DeleteToken(managerID, accessToken); err != nil {
		response.Error(c, err)
		return
	}

	response.SucceedWithMessage(c, "Logout success")
}
