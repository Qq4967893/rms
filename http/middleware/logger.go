package middleware

import (
	"backend-rms/core/logger"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		uri := c.Request.RequestURI
		method := c.Request.Method

		c.Next()

		fields := []interface{}{
			"ip", c.ClientIP(),
			"method", method,
			"uri", uri,
			"status", c.Writer.Status(),
			"elapsed_time", time.Since(start).String(),
		}

		logger.AccessLogger().WithContext(c).Infow("access", fields...)
	}
}
