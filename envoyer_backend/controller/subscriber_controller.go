package controller

import (
	"envoyer/config/consts"
	"envoyer/config/service_name"
	"envoyer/logger"
	"envoyer/model/service"
	"github.com/sarulabs/di/v2"
)

type SubscriberController struct {
	logger     logger.Logger
	dispatcher *service.Dispatcher
	appService service.AppServiceInterface
	container  di.Container
}

func NewSubscriberController(base *BaseController) *SubscriberController {
	return &SubscriberController{
		container:  base.container,
		logger:     base.logger,
		dispatcher: base.container.Get(service_name.DispatcherService).(*service.Dispatcher),
		appService: base.container.Get(service_name.AppService).(service.AppServiceInterface),
	}
}

func (c *SubscriberController) StartDefaultSubscriber() {
	//may be removed in future
	c.StartSubscriber(consts.DefaultQueue, "v2")
}

// StartSubscriber will start the subscriber
func (c *SubscriberController) StartSubscriber(queueName string, version string) {
	if c.dispatcher.ISQueueRunning(queueName) {
		c.logger.Debug("queue is already active", logger.Extra("queueName", queueName))
		return
	}

	err := c.dispatcher.Start(queueName, version)
	if err != nil {
		c.logger.Error("subscriber start failed", logger.Extra("error", err.Error()))
		return
	}
	c.logger.Info("consumer started", logger.Extra("queueName", queueName))
}

func (c AppController) StopSubscriber(queueName string) {
	if c.dispatcher.ISQueueRunning(queueName) {
		c.dispatcher.StopQueue(queueName)
	}
}

// StartAllSubscribers starts all subscribers for all active app as appName as queueName
func (c *SubscriberController) StartAllSubscribers() {
	//default queue. may be removed in future
	//c.StartSubscriber(consts.DefaultQueue, "v1")
	apps, err := c.appService.GetAll()
	if err != nil {
		c.logger.Error("Failed to start all subscribers", logger.Extra("error", err.Error()))
		return
	}

	for _, app := range apps {
		if app.Active {
			c.StartSubscriber(app.Name, "v2")
		}
	}
}
