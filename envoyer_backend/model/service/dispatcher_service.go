package service

import (
	"encoding/json"
	"envoyer/config/consts"
	"envoyer/config/service_name"
	"envoyer/logger"
	"envoyer/notification"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sarulabs/di/v2"
	"math"
	"sync"
	"sync/atomic"
	"time"
)

type Dispatcher struct {
	logger          logger.Logger
	destructor      sync.Once
	queues          map[string]bool
	m               sync.Mutex
	stopped         int32
	subscriber      Subscriber
	validityChecker notification.ValidityChecker
	publisher       Publisher
	handlerService  HandlerServiceInterface
}

func NewDispatcherService(subscriber Subscriber, logger logger.Logger, handlerService HandlerServiceInterface, container di.Container) *Dispatcher {
	dispatcher := &Dispatcher{
		logger:         logger,
		subscriber:     subscriber,
		handlerService: handlerService,
		queues:         make(map[string]bool),
		publisher:      container.Get(service_name.PublisherService).(Publisher),
	}
	notifyCloseCh := make(chan error)
	dispatcher.subscriber.NotifyClose(notifyCloseCh)
	atomic.StoreInt32(&dispatcher.stopped, 0)

	go dispatcher.watchConnection(notifyCloseCh)

	return dispatcher

}

func (d *Dispatcher) watchConnection(notifyCloseCh chan error) {
	err := <-notifyCloseCh
	d.logger.Error("Subscriber is closed", logger.Extra("error", err.Error()))
	d.m.Lock()
	atomic.StoreInt32(&d.stopped, 1)
	d.m.Unlock()
}

// Start will start the subscriber and keeps listening for messages and handle them and return necessary
// event like acknowledged, rejected, dispatched etc.
func (d *Dispatcher) Start(queueName string, version string) error {
	d.m.Lock()
	if _, found := d.queues[queueName]; found {
		d.m.Unlock()
		d.logger.Info("queue is already running")
		return nil
	} else {
		d.queues[queueName] = true
		d.m.Unlock()
	}

	deliveryCh, err := d.subscriber.Subscribe(queueName)
	if err != nil {
		d.logger.Error("subscribe failed", logger.Extra("subscribeError", err.Error()))
		return err
	}

	go d.ProcessDelivery(queueName, version, deliveryCh)

	return nil
}

func (d *Dispatcher) ProcessDelivery(queueName string, version string, deliveryCh <-chan amqp.Delivery) {
	waitChan := make(chan struct{}, consts.MaxConcurrentHandle)
Loop:
	for {
		amqpDelivery, ok := <-deliveryCh
		if ok {
			delivery := notification.Delivery{
				Acknowledger: d.subscriber,
				Id:           amqpDelivery.DeliveryTag,
				MessageType:  amqpDelivery.Type,
				Exchange:     amqpDelivery.Exchange,
				RoutingKey:   amqpDelivery.RoutingKey,
				Body:         amqpDelivery.Body,
				Headers:      amqpDelivery.Headers,
				Timestamp:    time.Now(),
			}

			var message notification.Message
			err := json.Unmarshal(delivery.Body, &message)
			if err != nil {
				d.logger.Error("Error decoding delivery message", logger.Extra("jsonError", err.Error()))
				d.sendDispatcherEvent(notification.Dispatch, delivery, fmt.Errorf("dispatcher couldn't parse the received message %w", err))
				d.reject(delivery, false)
				continue
			}

			handler := d.handlerService.GetHandler(message.MessageType, version)

			if handler != nil && handler.CanHandle(message.MessageType) {
				waitChan <- struct{}{}
				go d.handleDelivery(delivery, message, handler, waitChan)
			} else {
				d.sendDispatcherEvent(notification.Dispatch, delivery, notification.ErrNoHandlerFound)
				d.reject(delivery, false)
			}
		} else {
			d.logger.Info("consumer stopped for queue", logger.Extra("queueName", queueName))
			d.m.Lock()
			if _, found := d.queues[queueName]; found {
				delete(d.queues, queueName)
			}
			d.m.Unlock()
			break Loop
		}
	}
}

// Stop may get an error that indicates the dispatcher may not able to close the server connection but
// the dispatcher should be treated as stopped regardless.
func (d *Dispatcher) Stop() {
	_ = d.subscriber.CloseSubscriber()
}

// ISQueueRunning will return if  the queue is active
func (d *Dispatcher) ISQueueRunning(queueName string) bool {
	d.m.Lock()
	defer d.m.Unlock()
	if _, found := d.queues[queueName]; found {
		return true
	}
	return false
}

// StopQueue will stop the queue
func (d *Dispatcher) StopQueue(queueName string) {
	_ = d.subscriber.CloseQueue(queueName)
}

// IsStopped will return true if the dispatcher is stopped
func (d *Dispatcher) IsStopped() bool {
	return atomic.LoadInt32(&d.stopped) == 1
}

// handleDelivery will handle the delivery from RabbitMQ and process them and requeue them if necessary
func (d *Dispatcher) handleDelivery(delivery notification.Delivery, message notification.Message, handler notification.Handler, waitChan chan struct{}) {
	defer func() { <-waitChan }()
	//handle the message
	err, requeue := handler.Handle(&message)
	if err != nil {
		d.sendDispatcherEvent(notification.Dispatch, delivery, err)
		d.reject(delivery, false)
		if requeue {
			message.RequeueCount += 1
			requeueCount := message.RequeueCount

			if requeueCount <= consts.RequeueLimit {
				delayTime := math.Pow(2, float64(requeueCount))
				pubErr := d.publisher.Publish(notification.Request{
					Message:      message,
					DeliveryTime: time.Now().Add(time.Duration(delayTime) * time.Second),
					Queue:        delivery.RoutingKey,
				})
				if pubErr != nil {
					d.logger.Error("Failed to publish notification", logger.Extra("publishError", pubErr.Error()))
					d.sendDispatcherEvent(notification.RejectAndDelete, delivery, pubErr)
				} else {
					d.sendDispatcherEvent(notification.RejectAndRequeue, delivery, err)
				}
			}
		}
		return
	}

	d.sendDispatcherEvent(notification.Dispatch, delivery, nil)
	err = delivery.Acknowledger.Ack(delivery.Id, false)
	d.sendDispatcherEvent(notification.Ack, delivery, err)
}

// reject will send reject to the RabbitMQ
func (d *Dispatcher) reject(delivery notification.Delivery, requeue bool) {
	err := delivery.Acknowledger.Reject(delivery.Id, requeue)

	var eventType notification.DispatcherEventType
	if requeue {
		eventType = notification.RejectAndRequeue
	} else {
		eventType = notification.RejectAndDelete
	}

	d.sendDispatcherEvent(eventType, delivery, err)
}

// sendDispatcherEvent will set necessary events to chan that will be logged
func (d *Dispatcher) sendDispatcherEvent(eventType notification.DispatcherEventType, delivery notification.Delivery, err error) {
	if d.IsStopped() {
		d.logger.Error("dispatcher is already closed")
		return
	}

	dispatcherEvent := notification.DispatcherEvent{
		DispatcherEventType: eventType,
		Delivery:            delivery,
		Error:               err,
	}

	d.logEvents(dispatcherEvent)
}

func (d *Dispatcher) logEvents(dispatcherEvent notification.DispatcherEvent) {
	switch dispatcherEvent.DispatcherEventType {
	case notification.Dispatch:
		d.logger.Debug("Dispatch", logger.Extra("event", dispatcherEvent), logger.Extra("queue", dispatcherEvent.Delivery.RoutingKey))
		d.logger.Info("Message is dispatched", logger.Extra("queue", dispatcherEvent.Delivery.RoutingKey))
	case notification.Ack:
		d.logger.Debug("Act", logger.Extra("event", dispatcherEvent))
		d.logger.Info("Message is acknowledged", logger.Extra("queue", dispatcherEvent.Delivery.RoutingKey))
	case notification.RejectAndDelete:
		d.logger.Debug("RejectAndDelete", logger.Extra("event", dispatcherEvent))
		d.logger.Info("Message is rejected and deleted", logger.Extra("queue", dispatcherEvent.Delivery.RoutingKey))
	case notification.RejectAndRequeue:
		d.logger.Debug("RejectAndRequeue", logger.Extra("event", dispatcherEvent))
		d.logger.Info("Message is rejected and re queued", logger.Extra("queue", dispatcherEvent.Delivery.RoutingKey))
	case notification.Terminate:
		d.logger.Debug("Terminate", logger.Extra("event", dispatcherEvent))
		d.logger.Info("Connection  is terminated", logger.Extra("queue", dispatcherEvent.Delivery.RoutingKey))
	}
}
