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
	"gorm.io/gorm"
	"strings"
	"time"
)

type pushHandler struct {
	logger logger.Logger
	Db     *gorm.DB
}

type PushReq struct {
	serializers.PublishReq
	AppId             uint
	AppName           string
	EventType         string              `json:"event_name" binding:"required,_alpha-num,min=2,max=50"`
	DeliveryTime      *time.Time          `json:"delivery_time,omitempty"`
	Receivers         []string            `json:"receivers,omitempty"`
	TemplateVariables []TemplateVeriable  `json:"variables,omitempty"`
	BulkMessage       []MultiNotification `json:"receivers_with_variables,omitempty"`
	Data              map[string]string   `json:"data,omitempty"`
	ImageUrl          string              `json:"image_url,omitempty"`
	Topic             string              `json:"topic,omitempty"`
	Condition         string              `json:"condition,omitempty"`
	Language          string              `json:"language,omitempty"`
}

func (m PushReq) GetAppId() uint {
	return m.AppId
}

func (m PushReq) GetClientKey() string {
	return m.ClientKey
}

func (m PushReq) GetEventName() string {
	return m.EventType
}

var singletonPushHandler *pushHandler

func GetPushHandler(Db *gorm.DB, logger logger.Logger) Handler {
	if singletonPushHandler == nil {
		lock.Lock()
		defer lock.Unlock()
		if singletonPushHandler == nil {
			singletonPushHandler = &pushHandler{logger: logger, Db: Db}
		}
	}
	return singletonPushHandler
}

func (m pushHandler) GetRequest(context *gin.Context, messageType MessageType) (Request, *errors.RestErr) {
	var payload PushReq
	err := context.ShouldBindBodyWith(&payload, binding.JSON)
	if err != nil {
		m.logger.Error("failed to bind parameters")
		return Request{}, errors.NewBadRequestError("can not parse request parameters", err)
	}
	payload.AppId = context.GetUint("app_id")
	payload.AppName = context.GetString("app_name")

	deliveryTime := time.Now()
	if payload.DeliveryTime != nil && !payload.DeliveryTime.Before(time.Now()) {
		deliveryTime = *payload.DeliveryTime
	}

	smsContentJson, jsonErr := json.Marshal(payload)
	if jsonErr != nil {
		m.logger.Error("failed to marshal request", logger.Extra("jsonMarshalError", jsonErr.Error()))
		return Request{}, errors.NewBadRequestError("failed to marshal sms request", jsonErr)
	}

	request := Request{
		Message: Message{
			MessageType:  messageType,
			Body:         smsContentJson,
			RequeueCount: 0,
		},
		DeliveryTime: deliveryTime,
		Queue:        payload.AppName,
	}
	return request, nil
}

func (m pushHandler) Handle(message *Message) (error, bool) {
	var req PushReq
	err := json.Unmarshal(message.Body, &req)
	if err != nil {
		m.logger.Error("failed to unmarshal message", logger.Extra("jsonUnmarshalError", err.Error()))
		return err, false
	}

	var event entity.Event
	err = m.Db.Model(&entity.Event{}).Where("app_id = ? and name = ?", req.AppId, req.EventType).First(&event).Error
	if err != nil {
		m.logger.Error("failed to get event", logger.Extra("eventError", err.Error()))
		ErrorLog(m.logger, m.Db, req, nil, 0, "event not found", false)
		return err, false
	}
	if len(event.Variables) > 0 {
		event.VariableAsArray = strings.Split(event.Variables, ",")
	}

	language := "en"
	if len(req.Language) > 0 {
		language = req.Language
	}
	var template entity.Template
	err = m.Db.Model(&entity.Template{}).Where("event_id = ? and type = ? and active = true and language = ?", event.ID, consts.Push, language).First(&template).Error
	if err != nil {
		m.logger.Error("failed to get template", logger.Extra("templateError", err.Error()))
		ErrorLog(m.logger, m.Db, req, nil, 0, "template not found", false)
		return err, false
	}

	var providers []entity.Provider
	resultErr := m.Db.Where("app_id = ? and type = ? and active = true", req.AppId, consts.Push).Order("priority").Find(&providers).Error
	if resultErr != nil {
		m.logger.Error("failed to get provider", logger.Extra("providerError", resultErr.Error()))
		ErrorLog(m.logger, m.Db, req, nil, 0, "provider not found", false)
		return resultErr, false
	}

	if len(providers) == 0 {
		m.logger.Error("provider not found")
		ErrorLog(m.logger, m.Db, req, nil, 0, "provider not found", false)
		return errors.NewError("provider not found"), false
	}

	provider := providers[int(message.RequeueCount)%len(providers)]
	if provider.ProviderType == consts.FirebaseFCM {
		firebaseFCM := GetFirebaseFCMProvider(provider, m.logger, m.Db)
		err, requeue := firebaseFCM.Handle(template, req, event)
		if err != nil {
			return err, requeue
		}
	}

	return nil, false
}

func (m pushHandler) CanHandle(messageType MessageType) bool {
	return messageType == Push
}
