package utils

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"reflect"
	"regexp"
)

func init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("_alpha-num", CustomStringFormat)
		v.RegisterValidation("variable_format", CustomVariableFormat)
	}
}

var CustomStringFormat validator.Func = func(fl validator.FieldLevel) bool {
	value := fl.Field()
	ans := false
	switch value.Kind() {
	case reflect.String:
		ans = CheckString(value.String())
	case reflect.Slice, reflect.Array:
		ans = CheckArray(value.Interface().([]string))
	}
	return ans
}

func CheckString(str string) bool {
	match, err := regexp.MatchString("^[A-Za-z0-9_-]*$", str)
	if err != nil {
		return false
	}
	if !match {
		return false
	}
	return true
}

func CheckArray(str []string) bool {
	for _, v := range str {
		if !CheckString(v) {
			return false
		}
	}
	return true
}

var CustomVariableFormat validator.Func = func(fl validator.FieldLevel) bool {
	value := fl.Field()
	ans := false
	switch value.Kind() {
	case reflect.String:
		ans = CheckVariable(value.String())
	case reflect.Slice, reflect.Array:
		ans = CheckVariableArray(value.Interface().([]string))
	}
	return ans
}

func CheckVariable(str string) bool {
	match, err := regexp.MatchString("^{{\\.[A-Za-z0-9_-]+}}$", str)
	if err != nil {
		return false
	}
	if !match {
		return false
	}
	return true
}

func CheckVariableArray(str []string) bool {
	for _, v := range str {
		if !CheckVariable(v) {
			return false
		}
	}
	return true
}
