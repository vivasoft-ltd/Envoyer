package notification

import (
	"envoyer/errors"
	"envoyer/logger"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type noHandler struct {
	logger logger.Logger
	Db     *gorm.DB
}

var singletonNoHandler *noHandler

func GetNoHandler(Db *gorm.DB, logger logger.Logger) Handler {
	if singletonNoHandler == nil {
		lock.Lock()
		defer lock.Unlock()
		if singletonNoHandler == nil {
			singletonNoHandler = &noHandler{logger: logger, Db: Db}
		}
	}
	return singletonNoHandler
}

func (m noHandler) CanHandle(messageType MessageType) bool {
	return false
}

func (m noHandler) Handle(message *Message) (error, bool) {
	return errors.NewError("handler can not handle the message"), false
}

func (m noHandler) GetRequest(context *gin.Context, messageType MessageType) (Request, *errors.RestErr) {
	m.logger.Error("no handler found for this message type", logger.Extra("messageType", messageType))
	return Request{}, errors.NewBadRequestError("notification type is not supported", nil)
}
