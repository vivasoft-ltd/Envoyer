package notification

import (
	"encoding/json"
	"envoyer/errors"
	"envoyer/logger"
	"envoyer/model/entity"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"time"
)

type Handler interface {
	CanHandle(messageType MessageType) bool
	Handle(message *Message) (error, bool)
	GetRequest(context *gin.Context, messageType MessageType) (Request, *errors.RestErr)
}

type SmsHandler interface {
	Handle(template entity.Template, smsReq SmsReqV2, event entity.Event) (error, *RequeueError)
}

type LogRequest interface {
	GetAppId() uint
	GetClientKey() string
	GetEventName() string
}

func ErrorLog(log logger.Logger, db *gorm.DB, request LogRequest, message interface{}, providerId uint, errMsg string, requeue bool) {
	requestByte, err := json.Marshal(request)
	if err != nil {
		log.Error("failed to marshal message", logger.Extra("jsonMarshalError", err.Error()))
		return
	}

	messageByte, err := json.Marshal(message)
	if err != nil {
		log.Error("failed to marshal message", logger.Extra("jsonMarshalError", err.Error()))
		return
	}

	now := time.Now()
	errorLog := &entity.ErrorLog{
		AppId:      request.GetAppId(),
		ClientKey:  request.GetClientKey(),
		EventName:  request.GetEventName(),
		ProviderId: providerId,
		Message:    errMsg,
		Data:       messageByte,
		Request:    requestByte,
		Type:       "Error",
		Date:       &now,
		IsRequeue:  requeue,
	}
	err = db.Create(errorLog).Error
	if err != nil {
		log.Error("Error logging in db failed.", logger.Extra("error", err.Error()))
	}
	return
}
