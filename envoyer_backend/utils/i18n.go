package utils

import (
	ginI18n "github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func Trans(key string, data map[string]interface{}) string {
	param := &i18n.LocalizeConfig{
		MessageID:    key,
		TemplateData: data,
	}
	message := ginI18n.MustGetMessage(param)
	if message == "" {
		message = key
	}
	return message
}

func TransValidationErrors(errors validator.ValidationErrors) map[string]interface{} {
	errs := make(map[string]interface{}, 0)

	for _, err := range errors {
		str := TransValidationMessageKey(err)
		errs[err.Field()] = Trans(str, gin.H{
			"value": err.Value(),
			"param": err.Param(),
			"tag":   err.Tag(),
			"field": err.Field(),
		})
	}

	return errs
}

func TransValidationMessageKey(err validator.FieldError) string {

	switch err.Field() {
	case "PasswordConfirm":
		return "validationPassConfirm"
	}

	switch err.Tag() {
	case "email":
		return "validationEmail"
	case "required":
		return "validationRequired"
	case "unique":
		return "validationUnique"
	case "mobile-validation":
		return "invalidMobileNumber"
	}

	return err.Error()
}
