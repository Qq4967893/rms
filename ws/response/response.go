package response

import (
	"backend-rms/core/logger"
	"backend-rms/util/email"
	"backend-rms/ws"
	"fmt"
)

type response struct {
	Code    StatusCode `json:"code"`
	Message string     `json:"message,omitempty"`
	Data    any        `json:"data,omitempty"`
}

type Option func(*response)

func WithCode(code StatusCode) Option {
	return func(res *response) {
		res.Code = code
	}
}

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

func newResponse(opts ...Option) *response {
	resp := &response{}
	for _, o := range opts {
		o(resp)
	}
	return resp
}

func Write(c *ws.Context, opts ...Option) {
	c.JSON(newResponse(opts...))
}

func Abort(c *ws.Context, opts ...Option) {
	c.Abort()
	Write(c, opts...)
}

func Error(c *ws.Context, err error) {
	logger.WithContext(c).ErrorIf(err)
	email.SendException(fmt.Sprintf("%v", err))
	Abort(c, WithCode(StatusError), WithMessage("server error"))
}
