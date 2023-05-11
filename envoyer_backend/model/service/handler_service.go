package service

import (
	"envoyer/config/service_name"
	"envoyer/errors"
	"envoyer/logger"
	"envoyer/model/entity"
	"envoyer/model/repository"
	"envoyer/model/serializers"
	"envoyer/notification"
	"gorm.io/gorm"
)

type HandlerServiceInterface interface {
	GetHandler(messageType notification.MessageType, version string) notification.Handler
	ValidateRequest(publishReq serializers.PublishReq) (*entity.Application, *entity.Client, *errors.RestErr)
}

type handlerService struct {
	Db         *gorm.DB
	logger     logger.Logger
	appRepo    repository.AppRepositoryInterface
	clientRepo repository.ClientRepositoryInterface
}

func NewHandlerService(baseService *BaseService) HandlerServiceInterface {
	return &handlerService{
		Db:         baseService.container.Get(service_name.DbService).(*gorm.DB),
		logger:     baseService.logger,
		appRepo:    baseService.container.Get(service_name.AppRepository).(repository.AppRepositoryInterface),
		clientRepo: baseService.container.Get(service_name.ClientRepository).(repository.ClientRepositoryInterface),
	}
}

func (s handlerService) GetHandler(messageType notification.MessageType, version string) notification.Handler {
	if version == "v2" {
		return s.getHandlerV2(messageType)
	}
	return nil
}

func (s handlerService) getHandlerV2(messageType notification.MessageType) notification.Handler {
	switch messageType {
	case notification.Custom:
		return notification.GetCustomHandler(s.Db, s.logger)
	case notification.Email:
		return notification.GetEmailHandlerV2(s.Db, s.logger)
	case notification.Sms:
		return notification.GetSmsHandlerV2(s.Db, s.logger)
	case notification.Push:
		return notification.GetPushHandler(s.Db, s.logger)
	case notification.Webhook:
		return notification.GetWebhookHandler(s.Db, s.logger)
	}

	return notification.GetNoHandler(s.Db, s.logger)
}

func (s handlerService) ValidateRequest(publishReq serializers.PublishReq) (*entity.Application, *entity.Client, *errors.RestErr) {
	app, err := s.appRepo.GetByAppKey(publishReq.AppKey)
	if err != nil {
		return nil, nil, err
	}

	client, err := s.clientRepo.GetByClientKey(publishReq.ClientKey)
	if err != nil {
		return nil, nil, err
	}

	if client.AppId != app.ID {
		return nil, nil, errors.NewUnauthorizedError("invalid secret key", errors.NewError("invalid secret key"))
	}

	if !app.Active {
		s.logger.Info("app is not active", logger.Extra("appName", app.Name))
		return nil, nil, errors.NewForbiddenError("App is not active", nil)
	}

	return app, client, nil
}
