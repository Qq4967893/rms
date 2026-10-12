package middleware

import (
	"backend-rms/core/jwt"
	"backend-rms/http/response"
	"backend-rms/service"
	"strconv"

	"github.com/getsentry/sentry-go"
	"github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

var exclusionPaths = []string{
	"/api/auth/code",
	"/api/auth/login",
	"/api/auth/refresh_token",
}

func AuthJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, path := range exclusionPaths {
			if path == c.Request.URL.Path {
				return
			}
		}

		accessToken, err := jwt.Instance.GetTokenFromHeader(c.Request)
		if err != nil {
			if errors.Is(err, jwt.ErrHeaderEmpty) {
				response.Abort(c, response.StatusUnauthorized, response.CodeUnauthorized,
					response.WithMessage("Unauthenticated."),
				)
				return
			}

			response.Error(c, err)
			return
		}

		claims, err := jwt.Instance.ValidateToken(accessToken)
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				response.Abort(c, response.StatusUnauthorized, response.CodeAccessTokenExpired,
					response.WithMessage("Token is expired."),
				)
				return
			}

			response.Abort(c, response.StatusUnauthorized, response.CodeUnauthorized,
				response.WithMessage("Token is invalid."),
			)
			return
		}

		isAccessTokenExisted, err := service.AuthService.IsAccessTokenExisted(claims.LoginManagerID, accessToken)
		if err != nil {
			response.Error(c, err)
			return
		}

		if !isAccessTokenExisted {
			response.Abort(c, response.StatusUnauthorized, response.CodeAccessTokenExpired,
				response.WithMessage("Token is expired."),
			)
			return
		}

		c.Set("login_manager_id", claims.LoginManagerID)
		if hub := sentrygin.GetHubFromContext(c); hub != nil {
			hub.Scope().SetUser(sentry.User{
				ID: strconv.FormatInt(claims.LoginManagerID, 10),
			})
		}

		c.Next()
	}
}
