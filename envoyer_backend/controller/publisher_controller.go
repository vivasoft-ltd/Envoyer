package controller

import (
	"envoyer/config/service_name"
	"envoyer/logger"
	"envoyer/model/serializers"
	"envoyer/model/service"
	"envoyer/notification"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
)

type PublisherController struct {
	*BaseController
	publisher      service.Publisher
	handlerService service.HandlerServiceInterface
}

func NewPublisherController(base *BaseController) *PublisherController {
	return &PublisherController{
		BaseController: base,
		publisher:      base.container.Get(service_name.PublisherService).(service.Publisher),
		handlerService: base.container.Get(service_name.HandlerService).(service.HandlerServiceInterface),
	}
}

func (c PublisherController) PublishInQueueV2(context *gin.Context) {
	notificationType := context.Params.ByName("type")
	messageType, err := notification.StringToMessageType(notificationType)
	if err != nil {
		c.logger.Error("invalid message type", logger.Extra("type", notificationType))
		c.ReplyError(context, err.Error(), http.StatusBadRequest)
		return
	}

	var payload serializers.PublishReq
	err = context.ShouldBindBodyWith(&payload, binding.JSON)
	if err != nil {
		c.logger.Error("failed to bind parameters", logger.Extra("error", err.Error()))
		c.ReplyError(context, "can not parse request parameters", http.StatusBadRequest)
		return
	}

	app, _, validateErr := c.handlerService.ValidateRequest(payload)
	if validateErr != nil {
		c.logger.Error("invalid payload", logger.Extra("validationError", validateErr))
		c.ReplyError(context, validateErr.Error(), validateErr.Status)
		return
	}

	context.Set("app_name", app.Name)
	context.Set("app_id", app.ID)
	//context.Set("client_id", client.ID)
	//context.Set("client_name", client.Name)
	//context.Set("client_key", client.ClientKey)

	handler := c.handlerService.GetHandler(messageType, "v2")
	c.publish(context, messageType, handler)
	return
}

func (c PublisherController) publish(context *gin.Context, messageType notification.MessageType, handler notification.Handler) {
	request, restErr := handler.GetRequest(context, messageType)
	if restErr != nil {
		c.ReplyError(context, restErr.Error(), restErr.Status)
		return
	}

	err := c.publisher.Publish(request)
	if err != nil {
		c.logger.Error("Failed to publish notification", logger.Extra("publishError", err.Error()))
		c.ReplyError(context, err.Error(), http.StatusInternalServerError)
		return
	}

	c.ReplySuccess(context, "Successfully published to the queue")
	return
}
