package middleware

import (
	"backend-rms/core/jwt"
	"backend-rms/service"
	"backend-rms/ws"
	"backend-rms/ws/response"
	"net/http"
	"strings"
)

func Authenticate() ws.HandlerFunc {
	return func(c *ws.Context) {
		tokens := getTokens(c.Request.Header)
		if len(tokens) != 2 || tokens[0] != "rms" || tokens[1] == "" {
			write(c)
			return
		}

		accessToken := tokens[1]
		claims, err := jwt.Instance.ValidateToken(accessToken)
		if err != nil {
			write(c)
			return
		}

		isAccessTokenExisted, err := service.AuthService.IsAccessTokenExisted(claims.LoginManagerID, accessToken)
		if err != nil {
			write(c)
			return
		}

		if !isAccessTokenExisted {
			write(c)
			return
		}

		c.Set("login_manager_id", claims.LoginManagerID)
		c.Next()
	}
}

func getTokens(h http.Header) []string {
	var tokens []string
	protocols := strings.TrimSpace(h.Get("Sec-WebSocket-Protocol"))
	for _, t := range strings.Split(protocols, ",") {
		t = strings.TrimSpace(t)
		tokens = append(tokens, t)
	}
	return tokens
}

func write(c *ws.Context) {
	response.Abort(c,
		response.WithCode(response.StatusUnauthorized),
		response.WithMessage("unauthorized"),
	)
}
