package request

import (
	"backend-rms/core/logger"
	"backend-rms/ws"
	"backend-rms/ws/response"
	"encoding/json"

	"github.com/go-playground/validator/v10"
)

var (
	validate *validator.Validate
)

func init() {
	validate = validator.New()
}

func Validate(c *ws.Context, obj any) (isValid bool) {
	params, exists := c.Get("params")
	if !exists || params == nil {
		return
	}

	if err := json.Unmarshal(params.([]byte), obj); err != nil {
		write(c)
		return
	}

	if err := validate.Struct(obj); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			logger.Debugw("WS", "validationErrors", errs.Error())
		}
		write(c)
		return
	}

	return true
}

func write(c *ws.Context) {
	c.Set("params", nil)
	response.Write(c,
		response.WithCode(response.StatusInvalidParams),
		response.WithMessage("invalid params"),
	)
}
