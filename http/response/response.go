package response

import (
	"backend-rms/core/logger"
	"backend-rms/util/email"
	"fmt"
	"net/http"

	"github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
)

type response struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Data      any    `json:"data,omitempty"`
	Errors    any    `json:"errors,omitempty"`
	RequestID string `json:"request_id"`
}

type Option func(*response)

func WithMessage(message string) Option {
	return func(res *response) {
		res.Message = message
	}
}

func WithData(data any) Option {
	return func(res *response) {
		res.Data = data
	}
}

func WithErrors(errors any) Option {
	return func(res *response) {
		res.Errors = errors
	}
}

func newResponse(c *gin.Context, errorCode ErrorCode, opts ...Option) *response {
	resp := &response{
		Code:      int(errorCode),
		RequestID: c.MustGet("request_id").(string),
	}

	for _, o := range opts {
		o(resp)
	}

	return resp
}

func Succeed(c *gin.Context) {
	c.JSON(http.StatusOK, newResponse(c, CodeSuccess))
}

func SucceedWithMessage(c *gin.Context, message string) {
	c.JSON(http.StatusOK, newResponse(c, CodeSuccess, WithMessage(message)))
}

func SucceedWithData(c *gin.Context, message string, data any) {
	c.JSON(http.StatusOK, newResponse(c, CodeSuccess, WithMessage(message), WithData(data)))
}

func FailWithMessage(c *gin.Context, message string) {
	c.JSON(http.StatusOK, newResponse(c, CodeFailure, WithMessage(message)))
}

func FailWithData(c *gin.Context, message string, data any) {
	c.JSON(http.StatusOK, newResponse(c, CodeFailure, WithMessage(message), WithData(data)))
}

func Abort(c *gin.Context, statusCode StatusCode, errorCode ErrorCode, opts ...Option) {
	c.AbortWithStatusJSON(int(statusCode), newResponse(c, errorCode, opts...))
}

func Error(c *gin.Context, err error) {
	logger.WithContext(c).ErrorIf(err)
	email.SendHttpException(fmt.Sprintf("%v", err), c)
	if hub := sentrygin.GetHubFromContext(c); hub != nil {
		hub.CaptureException(err)
	}

	c.AbortWithStatusJSON(int(StatusServerError),
		newResponse(c, CodeServerError, WithMessage("Server Error")),
	)
}
