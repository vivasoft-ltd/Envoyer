package notification

import (
	"encoding/json"
	"envoyer/config/consts"
	"envoyer/errors"
	"envoyer/logger"
	"envoyer/model/entity"
	"envoyer/model/serializers"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	v "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"gorm.io/gorm"
	"strings"
	"time"
)

type EmailReqV2 struct {
	serializers.PublishReq
	AppId             uint
	AppName           string
	EventType         string              `json:"event_name" binding:"required,_alpha-num,min=2,max=50"`
	DeliveryTime      *time.Time          `json:"delivery_time,omitempty"`
	Receivers         []string            `json:"receivers,omitempty"`
	TemplateVariables []TemplateVeriable  `json:"variables,omitempty"`
	BulkMessage       []MultiNotification `json:"receivers_with_variables,omitempty"`
	Sender            string              `json:"sender,omitempty"`
	Cc                []string            `json:"cc,omitempty"`
	Bcc               []string            `json:"bcc,omitempty"`
	Language          string              `json:"language,omitempty"`
}

func (r *EmailReqV2) GetAppId() uint {
	return r.AppId
}

func (r *EmailReqV2) GetClientKey() string {
	return r.ClientKey
}

func (r *EmailReqV2) GetEventName() string {
	return r.EventType
}

func (r *EmailReqV2) Validate() error {

	if len(r.BulkMessage) > 0 {
		for _, msg := range r.BulkMessage {
			if len(msg.Receiver) <= 0 {
				return errors.NewError("Receiver not found")
			}

			err := v.Validate(&msg.Receiver, v.Required, is.Email)
			if err != nil {
				return errors.NewError("Invalid email address: " + err.Error())
			}
		}
	} else {
		if len(r.Receivers) <= 0 {
			return errors.NewError("Receiver not found")
		}

		for _, email := range r.Receivers {
			err := v.Validate(&email, v.Required, is.Email)
			if err != nil {
				return errors.NewError("Invalid email address: " + err.Error())
			}
		}

		for _, email := range r.Cc {
			err := v.Validate(&email, v.Required, is.Email)
			if err != nil {
				return errors.NewError("Invalid email address: " + err.Error())
			}
		}

		for _, email := range r.Bcc {
			err := v.Validate(&email, v.Required, is.Email)
			if err != nil {
				return errors.NewError("Invalid email address: " + err.Error())
			}
		}
	}

	return nil
}

type emailHandlerV2 struct {
	logger logger.Logger
	Db     *gorm.DB
}

var singletonEmailHandlerV2 *emailHandlerV2

func GetEmailHandlerV2(Db *gorm.DB, logger logger.Logger) Handler {
	if singletonEmailHandlerV2 == nil {
		lock.Lock()
		defer lock.Unlock()
		if singletonEmailHandlerV2 == nil {
			singletonEmailHandlerV2 = &emailHandlerV2{Db: Db, logger: logger}
		}
	}
	return singletonEmailHandlerV2
}

func (m emailHandlerV2) GetRequest(context *gin.Context, messageType MessageType) (Request, *errors.RestErr) {
	var payload EmailReqV2
	err := context.ShouldBindBodyWith(&payload, binding.JSON)
	if err != nil {
		m.logger.Error("failed to bind parameters")
		return Request{}, errors.NewBadRequestError("can not parse request parameters", err)
	}
	payload.AppId = context.GetUint("app_id")
	payload.AppName = context.GetString("app_name")

	validateErr := payload.Validate()
	if validateErr != nil {
		m.logger.Error("invalid payload", logger.Extra("validationError", validateErr))
		return Request{}, errors.NewBadRequestError("validation error", validateErr)
	}

	deliveryTime := time.Now()
	if payload.DeliveryTime != nil && !payload.DeliveryTime.Before(time.Now()) {
		deliveryTime = *payload.DeliveryTime
	}

	emailContentJson, jsonErr := json.Marshal(payload)
	if jsonErr != nil {
		m.logger.Error("failed to marshal request", logger.Extra("jsonMarshalError", jsonErr.Error()))
		return Request{}, errors.NewBadRequestError("failed to marshal email request", jsonErr)
	}
	request := Request{
		Message: Message{
			MessageType:  messageType,
			Body:         emailContentJson,
			RequeueCount: 0,
		},
		DeliveryTime: deliveryTime,
		Queue:        payload.AppName,
	}
	return request, nil
}

func (m emailHandlerV2) Handle(message *Message) (error, bool) {
	var emailReq EmailReqV2
	err := json.Unmarshal(message.Body, &emailReq)
	if err != nil {
		m.logger.Error("failed to unmarshal message", logger.Extra("jsonUnmarshalError", err.Error()))
		return err, false
	}

	var event entity.Event
	err = m.Db.Model(&entity.Event{}).Where("app_id = ? and name = ?", emailReq.AppId, emailReq.EventType).First(&event).Error
	if err != nil {
		m.logger.Error("failed to get event", logger.Extra("eventError", err.Error()))
		ErrorLog(m.logger, m.Db, &emailReq, nil, 0, "event not found", false)
		return err, false
	}
	if len(event.Variables) > 0 {
		event.VariableAsArray = strings.Split(event.Variables, ",")
	}

	language := "en"
	if len(emailReq.Language) > 0 {
		language = emailReq.Language
	}

	var template entity.Template
	err = m.Db.Model(&entity.Template{}).Where("event_id = ? and type = ? and active = true and language = ?", event.ID, consts.Email, language).First(&template).Error
	if err != nil {
		m.logger.Error("failed to get template", logger.Extra("templateError", err.Error()))
		ErrorLog(m.logger, m.Db, &emailReq, nil, 0, "template not found", false)
		return err, false
	}

	var providers []entity.Provider
	resultErr := m.Db.Where("app_id = ? and type = ? and active = true", emailReq.AppId, consts.Email).Order("priority").Find(&providers).Error
	if resultErr != nil {
		m.logger.Error("failed to get provider", logger.Extra("providerError", resultErr.Error()))
		ErrorLog(m.logger, m.Db, &emailReq, nil, 0, "provider not found", false)
		return resultErr, false
	}

	if len(providers) == 0 {
		m.logger.Error("provider not found")
		ErrorLog(m.logger, m.Db, &emailReq, nil, 0, "provider not found", false)
		return errors.NewError("provider not found"), false
	}

	provider := providers[int(message.RequeueCount)%len(providers)]
	if provider.ProviderType == consts.Smtp {
		smtpProvider := GetSmtpProvider(provider, m.logger, m.Db)
		err, requeue := smtpProvider.Handle(template, emailReq, event)
		if err != nil {
			return err, requeue
		}
	}

	return nil, false
}

func (m emailHandlerV2) CanHandle(messageType MessageType) bool {
	return messageType == Email
}
