package errors

import (
	"envoyer/logger"
	"envoyer/utils"
	"errors"
	"net/http"
)

type RestErr struct {
	Message string `json:"message"`
	Detail  string `json:"detail"`
	Status  int    `json:"status"`
}

func (err *RestErr) Error() string {
	return err.Message
}

func NewError(msg string) error {
	return errors.New(msg)
}

func NewRestErr(message string, status int, err error) *RestErr {
	restErr := &RestErr{
		Message: utils.Trans(message, nil),
		Status:  status,
	}

	if err != nil {
		logger.GetLogger().Error(err.Error())
		restErr.Detail = err.Error()
	} else {
		logger.GetLogger().Error(message)
	}

	return restErr
}

func NewInternalServerError(message string, err error) *RestErr {
	return NewRestErr(message, http.StatusInternalServerError, err)
}

func NewBadRequestError(message string, err error) *RestErr {
	return NewRestErr(message, http.StatusBadRequest, err)
}

func NewNotFoundError(message string, err error) *RestErr {
	return NewRestErr(message, http.StatusNotFound, err)
}

func NewAlreadyExistError(message string, err error) *RestErr {
	return NewRestErr(message, http.StatusConflict, err)
}

func NewUnauthorizedError(message string, err error) *RestErr {
	return NewRestErr(message, http.StatusUnauthorized, err)
}

func NewForbiddenError(message string, err error) *RestErr {
	return NewRestErr(message, http.StatusForbidden, err)
}
