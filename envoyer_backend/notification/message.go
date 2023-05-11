package notification

import (
	"errors"
	"time"
)

var (
	ErrClosed             = errors.New("subscriber/publisher is not open")
	ErrNoHandlerFound     = errors.New("dispatcher haven't found any handler")
	ErrUnknownMessageType = errors.New("dispatcher doesn't recognise the message type")
)

type MessageType string

const (
	Email   MessageType = "Email"
	Custom  MessageType = "Custom"
	Sms     MessageType = "Sms"
	Push    MessageType = "Push"
	Webhook MessageType = "Webhook"
	UnKnown MessageType = "UnKnown"
)

func StringToMessageType(messageType string) (MessageType, error) {
	switch messageType {
	case "email":
		return Email, nil
	case "custom":
		return Custom, nil
	case "sms":
		return Sms, nil
	case "push":
		return Push, nil
	case "webhook":
		return Webhook, nil
	default:
		return UnKnown, ErrUnknownMessageType
	}
}

type DispatcherEventType int

const (
	Dispatch DispatcherEventType = iota
	Ack
	RejectAndDelete
	RejectAndRequeue
	Terminate
)

type TemplateVeriable struct {
	Name  string `json:"name" binding:"variable_format"`
	Value string `json:"value"`
}

type MultiNotification struct {
	Receiver          string             `json:"receiver"`
	TemplateVariables []TemplateVeriable `json:"variables"`
}

type Message struct {
	MetaData     any
	RequeueCount uint
	MessageType  MessageType
	Body         []byte
}

type Request struct {
	Queue        string
	Message      Message
	DeliveryTime time.Time
}

type DispatcherEvent struct {
	DispatcherEventType DispatcherEventType
	Delivery            Delivery
	Error               error
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
