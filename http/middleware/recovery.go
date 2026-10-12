package middleware

import (
	"backend-rms/core/logger"
	"backend-rms/http/response"
	"errors"
	"net"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					var se *os.SyscallError
					if errors.As(ne, &se) {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") ||
							strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				if brokenPipe {
					// If the connection is dead, we can't write a status to it.
					if e, ok := err.(error); ok {
						logger.WithContext(c).ErrorIf(e)
					} else {
						logger.WithContext(c).Error(err)
					}
					c.Abort()
					return
				}

				if e, ok := err.(error); ok {
					response.Error(c, e)
					return
				}

				logger.WithContext(c).Error(err)
				response.Abort(c, response.StatusServerError, response.CodeServerError,
					response.WithMessage("Server Error"),
				)
			}
		}()
		c.Next()
	}
}
