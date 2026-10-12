package request

import (
	"backend-rms/http/request/validation"
	"backend-rms/http/response"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/validator/v10"

	ut "github.com/go-playground/universal-translator"
	en_translation "github.com/go-playground/validator/v10/translations/en"
)

var (
	validate *validator.Validate
	trans    ut.Translator
)

func init() {
	validate = validator.New()

	// register translator
	en := en.New()
	uni := ut.New(en, en)
	trans, _ = uni.GetTranslator("en")
	en_translation.RegisterDefaultTranslations(validate, trans)

	// register custom validations
	validate.RegisterValidation("exists", validation.ExistsValidation)
	validate.RegisterValidation("not_exists", validation.NotExistsValidation)
	validate.RegisterValidation("tick_queryable_time", validation.TickQueryableTime)

	validate.RegisterAlias("password_len", "len=32")
	validate.RegisterAlias("refresh_token_len", "len=32")

	registerTranslations(trans)
}

func Validate(c *gin.Context, data interface{}) bool {
	bindErr := c.ShouldBind(data)
	bindUriErr := c.ShouldBindUri(data)

	if bindErr != nil || bindUriErr != nil {
		response.Abort(c, response.StatusBadRequest, response.CodeBadRequest,
			response.WithMessage("Bad request."),
		)
		return false
	}

	err := validate.Struct(data)
	if err == nil {
		return true
	}

	if invalidErr, ok := err.(*validator.InvalidValidationError); ok {
		response.Error(c, invalidErr)
		return false
	}

	result := translate(err)
	response.Abort(c, response.StatusBadRequest, response.CodeBadRequest,
		response.WithMessage("The validation failed, please check errors for details."),
		response.WithErrors(result),
	)
	return false

}

func translate(err error) map[string][]string {
	var result = make(map[string][]string)

	errs := err.(validator.ValidationErrors)
	for _, e := range errs {
		result[e.Field()] = append(result[e.Field()], e.Translate(trans))
	}
	return result
}

func registerTranslations(trans ut.Translator) {
	validate.RegisterTranslation(
		"not_exists",
		trans,
		func(ut ut.Translator) error {
			return ut.Add("not_exists", "the value of {0} field exists already", true) // see universal-translator for details
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("not_exists", fe.Field())
			return t
		})

	validate.RegisterTranslation(
		"exists",
		trans,
		func(ut ut.Translator) error {
			return ut.Add("exists", "the value of {0} field doesn't exist", true) // see universal-translator for details
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("exists", fe.Field())
			return t
		})

	validate.RegisterTranslation(
		"password_len",
		trans,
		func(ut ut.Translator) error {
			return ut.Add("password_len", "Wrong password", true) // see universal-translator for details
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("password_len", fe.Field())
			return t
		})

	validate.RegisterTranslation(
		"refresh_token_len",
		trans,
		func(ut ut.Translator) error {
			return ut.Add("refresh_token_len", "Wrong refresh token", true) // see universal-translator for details
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("refresh_token_len", fe.Field())
			return t
		})

	validate.RegisterTranslation(
		"required_with",
		trans,
		func(ut ut.Translator) error {
			return ut.Add("required_with", "{0} is a required field", true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("required_with", fe.Field())
			return t
		},
	)

	validate.RegisterTranslation(
		"required_without_all",
		trans,
		func(ut ut.Translator) error {
			return ut.Add("required_without_all", "{0} is a required field", true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("required_without_all", fe.Field())
			return t
		},
	)

	validate.RegisterTranslation(
		"tick_queryable_time",
		trans,
		func(ut ut.Translator) error {
			return ut.Add("tick_queryable_time", "the value of {0} field is invalid", true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("tick_queryable_time", fe.Field())
			return t
		})
}
