package middleware

import (
	"github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
)

func Sentry() gin.HandlerFunc {
	return sentrygin.New(sentrygin.Options{
		Repanic: true,
	})
}
