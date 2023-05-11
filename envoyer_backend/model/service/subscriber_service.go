package service

import (
	"envoyer/config"
	"envoyer/logger"
	"envoyer/notification"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"sync"
	"sync/atomic"
	"time"
)

type Subscriber interface {
	Subscribe(queueName string) (<-chan amqp.Delivery, error)
	CloseQueue(queueName string) error
	NotifyClose(chan error)
	CloseSubscriber() error
	Reject(deliveryId uint64, requeue bool) error
	Ack(deliveryId uint64, multiple bool) error
}

type rabbitMQSubscriber struct {
	logger         logger.Logger
	conn           *amqp.Connection
	channel        *amqp.Channel
	tag            string
	closed         int32
	m              sync.Mutex
	notifyCloseCh  chan<- error
	destructor     sync.Once
	consumerConfig config.ConsumerConfig
}

// NewRabbitMQSubscriber will Create a new connection and channel for the rabbitmq consumer
func NewRabbitMQSubscriber(logger logger.Logger) (Subscriber, error) {
	logger.Info("Creating new rabbit consumer ...")

	//create connection to RabbitMQ
	conn, err := amqp.Dial(config.Config.ConsumerConfig.ServerUrl)
	if err != nil {
		logger.Error("Error creating connection: " + err.Error())
		return nil, err
	}

	// Create channel
	channel, err := conn.Channel()
	if err != nil {
		logger.Error("Error connecting channel: " + err.Error())
		return nil, fmt.Errorf(":Error creating channel: %s", err)
	}

	r := &rabbitMQSubscriber{
		logger:         logger,
		channel:        channel,
		conn:           conn,
		consumerConfig: config.Config.ConsumerConfig,
	}

	err = r.connect()
	if err != nil {
		return nil, fmt.Errorf(":Error from connecting rabbitmq:%s", err)
	}

	go r.watchConnection()

	return r, nil
}

// connect will connect to the rabbitmq if not already connected and declare exchange
func (subscriber *rabbitMQSubscriber) connect() error {
	if subscriber.conn.IsClosed() {
		conn, err := amqp.Dial(subscriber.consumerConfig.ServerUrl)
		if err != nil {
			subscriber.logger.Error("Error creating connection", logger.Extra("error", err.Error()))
			return err
		}
		subscriber.conn = conn
	}

	if subscriber.channel.IsClosed() {
		channel, err := subscriber.conn.Channel()
		if err != nil {
			subscriber.logger.Error("Error creating channel", logger.Extra("error", err.Error()))
			return fmt.Errorf(":Channel: %s", err)
		}
		subscriber.channel = channel
	}

	args := make(amqp.Table)
	args["x-delayed-type"] = "direct"

	//declare exchange
	if err := subscriber.channel.ExchangeDeclare(
		subscriber.consumerConfig.ExchangeConfig.Name,       // name of the exchange
		subscriber.consumerConfig.ExchangeConfig.Type,       // type
		subscriber.consumerConfig.ExchangeConfig.Durable,    // durable
		subscriber.consumerConfig.ExchangeConfig.AutoDelete, // delete when complete
		subscriber.consumerConfig.ExchangeConfig.Internal,   // internal
		subscriber.consumerConfig.ExchangeConfig.NoWait,     // noWait
		args, // arguments
	); err != nil {
		subscriber.logger.Error("Error declaring exchange", logger.Extra("error", err.Error()))
		return fmt.Errorf(":Exchange Declare: %s", err)
	}

	atomic.StoreInt32(&subscriber.closed, 0)

	return nil
}

// watchConnection will watch the  connection of the rabbitmq and try to reconnect if necessary
func (subscriber *rabbitMQSubscriber) watchConnection() {
Loop:
	for {
		select {
		case _ = <-subscriber.channel.NotifyClose(make(chan *amqp.Error)):
			if subscriber.isClosed() {
				subscriber.logger.Info("...rabbitMQ has been shut down")
				break Loop
			}

			subscriber.logger.Error("RabbitMQ disconnection")
			subscriber.logger.Info("...trying to reconnect to rabbitMQ...")
			var err error
			for _, timeout := range subscriber.consumerConfig.BackoffPolicy {
				if err = subscriber.connect(); err != nil {
					time.Sleep(timeout * time.Second)
					continue
				} else {
					subscriber.logger.Info("successfully reconnected to RabbitMQ")
					continue Loop
				}
			}
			subscriber.logger.Error("failed to reconnect rabbitmq")
			subscriber.notifyError(fmt.Errorf(":Failed to reconnect rabbitmq: %s", err))
			_ = subscriber.CloseSubscriber()
			break Loop
		}
	}
}

func (subscriber *rabbitMQSubscriber) notifyError(err error) {
	subscriber.m.Lock()
	defer subscriber.m.Unlock()

	if subscriber.notifyCloseCh != nil {
		subscriber.notifyCloseCh <- err
	}
}

// CloseSubscriber returns an error that indicates the server may not have received this request to close but
// the connection should be treated as closed regardless.
func (subscriber *rabbitMQSubscriber) CloseSubscriber() error {
	if subscriber.isClosed() {
		subscriber.logger.Info("consumer is already closed")
		return notification.ErrClosed
	}

	defer subscriber.shutdown()

	if err := subscriber.channel.Cancel(subscriber.tag, false); err != nil {
		subscriber.logger.Error("error closing rabbitmq channel..", logger.Extra("channelCloseError", err.Error()))
		return err
	}

	if err := subscriber.conn.Close(); err != nil {
		subscriber.logger.Error("error closing rabbitmq connection..", logger.Extra("connCloseError", err.Error()))
		return err
	}

	return nil
}

// shutdown shuts down the subscriber and close all delivery channels of all queues
func (subscriber *rabbitMQSubscriber) shutdown() {
	atomic.StoreInt32(&subscriber.closed, 1)
	subscriber.logger.Info("shutting down rabbitMQ connection...")

	subscriber.destructor.Do(func() {
		if subscriber.notifyCloseCh != nil {
			close(subscriber.notifyCloseCh)
		}
	})
}

// Subscribe is used to declare a queue, bind and start consuming form it.
func (subscriber *rabbitMQSubscriber) Subscribe(queueName string) (<-chan amqp.Delivery, error) {
	//declare queue
	queue, err := subscriber.channel.QueueDeclare(
		queueName, // name of the queue
		subscriber.consumerConfig.QueueConfig.Durable,    // durable
		subscriber.consumerConfig.QueueConfig.AutoDelete, // delete when unused
		subscriber.consumerConfig.QueueConfig.Exclusive,  // exclusive
		subscriber.consumerConfig.QueueConfig.NoWait,     // noWait
		nil, // arguments
	)
	if err != nil {
		subscriber.logger.Error("Error declaring queue", logger.Extra("error", err.Error()))
		return nil, fmt.Errorf(":Queue Declare: %s", err)
	}

	//bind queue
	if err = subscriber.channel.QueueBind(
		queue.Name, // name of the queue
		queueName,  // bindingKey
		subscriber.consumerConfig.ExchangeConfig.Name, // sourceExchange
		false, // noWait
		nil,   // arguments
	); err != nil {
		subscriber.logger.Error("Error binding queue", logger.Extra("error", err.Error()))
		return nil, fmt.Errorf(":Queue Bind: %s", err)
	}

	//start consuming from the queue
	deliveries, err := subscriber.channel.Consume(
		queueName,                         // name
		queueName,                         // consumerTag,
		subscriber.consumerConfig.AutoAck, // autoAck
		false,                             // exclusive
		false,                             // noLocal
		false,                             // noWait
		nil,                               // arguments
	)

	return deliveries, err
}

func (subscriber *rabbitMQSubscriber) NotifyClose(notifyCloseCh chan error) {
	subscriber.m.Lock()
	defer subscriber.m.Unlock()
	subscriber.notifyCloseCh = notifyCloseCh
}

// CloseQueue close the delivery channel of that queue
func (subscriber *rabbitMQSubscriber) CloseQueue(queueName string) error {
	if subscriber.isClosed() {
		subscriber.logger.Info("consumer is already closed")
		return notification.ErrClosed
	}

	//cancel consuming from a queue
	err := subscriber.channel.Cancel(queueName, false)
	if err != nil {
		subscriber.logger.Error("closing consumer failed: ", logger.Extra("channelCancelErr", err.Error()))
		return err
	}

	return nil
}

// Ack will send acknowledgement to RabbitMQ
func (subscriber *rabbitMQSubscriber) Ack(deliveryId uint64, multiple bool) error {
	return subscriber.channel.Ack(deliveryId, multiple)
}

// Reject will send reject to the RabbitMQ
func (subscriber *rabbitMQSubscriber) Reject(deliveryId uint64, requeue bool) error {
	return subscriber.channel.Reject(deliveryId, requeue)
}

// isClosed returns true if the subscriber is closed
func (subscriber *rabbitMQSubscriber) isClosed() bool {
	return atomic.LoadInt32(&subscriber.closed) == 1
}
