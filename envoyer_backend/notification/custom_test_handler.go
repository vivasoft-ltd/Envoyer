package notification

import (
	"encoding/json"
	"envoyer/errors"
	"envoyer/logger"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gorm.io/gorm"
	"time"
)

type custom struct {
	AppName string
	Data    string `json:"data"`
}

type customHandler struct {
	logger logger.Logger
	Db     *gorm.DB
}

var singletonCustomHandler *customHandler

// GetCustomHandler is a test handler that just prints the message in terminal.
func GetCustomHandler(Db *gorm.DB, logger logger.Logger) Handler {
	if singletonCustomHandler == nil {
		lock.Lock()
		defer lock.Unlock()
		if singletonCustomHandler == nil {
			singletonCustomHandler = &customHandler{logger: logger, Db: Db}
		}
	}
	return singletonCustomHandler
}

func (m customHandler) GetRequest(context *gin.Context, messageType MessageType) (Request, *errors.RestErr) {
	var payload custom

	err := context.ShouldBindBodyWith(&payload, binding.JSON)
	if err != nil {
		m.logger.Error("failed to bind parameters")
		return Request{}, errors.NewBadRequestError("can not parse request parameters", err)
	}
	payload.AppName = context.GetString("app_name")

	customContentJson, jsonErr := json.Marshal(payload)
	if jsonErr != nil {
		m.logger.Error("failed to marshal request", logger.Extra("jsonMarshalError", jsonErr.Error()))
		return Request{}, errors.NewBadRequestError("failed to marshal custom request", jsonErr)
	}

	request := Request{
		Message: Message{
			MessageType:  messageType,
			Body:         customContentJson,
			RequeueCount: 0,
		},
		DeliveryTime: time.Now(),
		Queue:        payload.AppName,
	}

	return request, nil
}

func (m customHandler) Handle(message *Message) (error, bool) {
	var data custom
	err := json.Unmarshal(message.Body, &data)
	if err != nil {
		m.logger.Error("failed to unmarshal message", logger.Extra("jsonUnmarshalError", err.Error()))
		return err, false
	}

	fmt.Println(data)
	return nil, false
}

func (m customHandler) CanHandle(messageType MessageType) bool {
	return messageType == Custom
}
