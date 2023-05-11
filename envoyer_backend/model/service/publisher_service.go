package service

import (
	"context"
	"encoding/json"
	"envoyer/config"
	"envoyer/config/consts"
	"envoyer/logger"
	"envoyer/notification"
	"errors"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"sync"
	"sync/atomic"
	"time"
)

type Publisher interface {
	Publish(notification.Request) error
	Close() error
	NotifyErrorf(chan error) chan error
}

type rabbitMQPublisher struct {
	logger          logger.Logger
	conn            *amqp.Connection
	channel         *amqp.Channel
	closed          int32
	m               sync.Mutex
	errorChannel    chan error
	publisherConfig config.PublisherConfig
}

// NewRabbitMQPublisher will Create a new connection and channel for the rabbitmq publisher
func NewRabbitMQPublisher(logger logger.Logger) (Publisher, error) {
	logger.Info("Creating new rabbit publisher ...")

	rmqp := &rabbitMQPublisher{
		logger:          logger,
		publisherConfig: config.Config.PublisherConfig,
	}

	connection, err := amqp.Dial(rmqp.publisherConfig.ServerUrl)
	if err != nil {
		logger.Error("Error creating connection: " + err.Error())
		return nil, fmt.Errorf(":Error connecting:%s", err)
	}
	rmqp.conn = connection

	channel, err := rmqp.conn.Channel()
	if err != nil {
		logger.Error("Error connecting channel: " + err.Error())
		return nil, fmt.Errorf(":Error creating channel:%s", err)
	}
	rmqp.channel = channel

	err = rmqp.connect()
	if err != nil {
		return nil, fmt.Errorf(":Error from load function:%s", err)
	}

	go handleDisconnect(rmqp)

	return rmqp, nil
}

func (rmqp *rabbitMQPublisher) isClosed() bool {
	return atomic.LoadInt32(&rmqp.closed) == 1
}

func (rmqp *rabbitMQPublisher) connect() error {
	if rmqp.conn.IsClosed() {
		connection, err := amqp.Dial(rmqp.publisherConfig.ServerUrl)
		if err != nil {
			rmqp.logger.Error("Error creating connection", logger.Extra("error", err.Error()))
			return fmt.Errorf(":Error connecting:%s", err)
		}
		rmqp.conn = connection
	}

	if rmqp.channel.IsClosed() {
		channel, err := rmqp.conn.Channel()
		if err != nil {
			rmqp.logger.Error("Error creating channel", logger.Extra("error", err.Error()))
			return fmt.Errorf(":Error creating channel:%s", err)
		}
		rmqp.channel = channel
	}

	// arg for delayed message
	args := make(amqp.Table)
	args["x-delayed-type"] = "direct"

	err := rmqp.channel.ExchangeDeclare(
		rmqp.publisherConfig.ExchangeConfig.Name,
		rmqp.publisherConfig.ExchangeConfig.Type,
		rmqp.publisherConfig.ExchangeConfig.Durable,
		rmqp.publisherConfig.ExchangeConfig.AutoDelete,
		rmqp.publisherConfig.ExchangeConfig.Internal,
		rmqp.publisherConfig.ExchangeConfig.NoWait,
		args)
	if err != nil {
		rmqp.logger.Error("Error declaring exchange", logger.Extra("error", err.Error()))
		return fmt.Errorf(":Error declaring exchange:%s", err)
	}
	atomic.StoreInt32(&rmqp.closed, 0)
	return nil
}

func (rmqp *rabbitMQPublisher) errorf(err error) {
	rmqp.m.Lock()
	defer rmqp.m.Unlock()

	if rmqp.errorChannel != nil {
		rmqp.errorChannel <- err
	}
}

func (rmqp *rabbitMQPublisher) NotifyErrorf(errChannel chan error) chan error {
	if errChannel != nil {
		rmqp.errorChannel = errChannel
	}

	return rmqp.errorChannel
}

func handleDisconnect(rmqp *rabbitMQPublisher) {
Loop:
	for {
		select {
		case errChannel := <-rmqp.channel.NotifyClose(make(chan *amqp.Error)):
			if rmqp.isClosed() {
				rmqp.logger.Info("...rabbitMQ has been shut down")
				return
			}

			if errChannel != nil {
				rmqp.logger.Error("RabbitMQ disconnection", logger.Extra("errChannel", errChannel.Error()))
				rmqp.logger.Info("...trying to reconnect to rabbitMQ...")
				for _, timeout := range rmqp.publisherConfig.BackoffPolicy {
					if err := rmqp.connect(); err != nil {
						time.Sleep(timeout * time.Second)
						continue
					} else {
						rmqp.logger.Info("successfully reconnected to RabbitMQ")
						go handleDisconnect(rmqp)
						break Loop
					}
				}
				rmqp.logger.Error("failed to reconnect rabbitmq")
				rmqp.errorf(errors.New("failed to reconnect rabbitmq"))
				_ = rmqp.Close()
			}
		}
	}
}

func (rmqp *rabbitMQPublisher) Publish(request notification.Request) error {
	delayTimeUnix := request.DeliveryTime.Unix()
	currentTimeUnix := time.Now().Unix()
	delay := (delayTimeUnix - currentTimeUnix) * 1000

	requestByte, err := json.Marshal(request.Message)
	if err != nil {
		rmqp.logger.Error("failed to marshal message", logger.Extra("jsonMarshalError", err.Error()))
		return fmt.Errorf(":Error Marshaling the struct request:%s", err)
	}

	if len(request.Queue) <= 0 {
		request.Queue = consts.DefaultQueue
	}

	return rmqp.publishToRabbitMQ(requestByte, delay, request.Queue)
}

func (rmqp *rabbitMQPublisher) publishToRabbitMQ(body []byte, delay int64, queue string) error {
	headers := make(amqp.Table)
	rmqp.logger.Debug("publishing to Rabbitmq", logger.Extra("routing key", queue), logger.Extra("body", body))

	if delay > 0 {
		rmqp.logger.Debug("request has delay", logger.Extra("delay", delay))
		headers["x-delay"] = delay
	}

	if rmqp.isClosed() {
		if err := rmqp.connect(); err != nil {
			return err
		} else {
			rmqp.logger.Info("successfully reconnected to RabbitMQ")
			go handleDisconnect(rmqp)
		}
	}

	return rmqp.channel.PublishWithContext(context.Background(),
		rmqp.publisherConfig.ExchangeConfig.Name,
		queue,
		rmqp.publisherConfig.Mandatory,
		rmqp.publisherConfig.Immidiate,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
			ContentType:  rmqp.publisherConfig.ContentType,
			Body:         body,
			Headers:      headers,
		})
}

func (rmqp *rabbitMQPublisher) Close() error {
	defer rmqp.stop()

	if rmqp.isClosed() {
		rmqp.logger.Info("publisher is already closed")
		return errors.New("publisher is not open")
	}

	if err := rmqp.channel.Close(); err != nil {
		rmqp.logger.Error("error closing rabbitmq channel..", logger.Extra("channelCloseError", err.Error()))
		return err
	}

	if err := rmqp.conn.Close(); err != nil {
		rmqp.logger.Error("error closing rabbitmq connection..", logger.Extra("connCloseError", err.Error()))
		return err
	}

	return nil
}

func (rmqp *rabbitMQPublisher) stop() {
	atomic.StoreInt32(&rmqp.closed, 1)
	rmqp.logger.Info("shutting down rabbitMQ connection...")
}
