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
	"sync"
	"time"
)

type SmsReqV2 struct {
	serializers.PublishReq
	AppId             uint
	AppName           string
	EventType         string              `json:"event_name" binding:"required,_alpha-num,min=2,max=50"`
	DeliveryTime      *time.Time          `json:"delivery_time,omitempty"`
	Receivers         []string            `json:"receivers,omitempty"`
	TemplateVariables []TemplateVeriable  `json:"variables,omitempty"`
	BulkMessage       []MultiNotification `json:"receivers_with_variables,omitempty"`
	Language          string              `json:"language,omitempty"`
	Sender            string              `json:"sender,omitempty"`
}

func (r *SmsReqV2) GetAppId() uint {
	return r.AppId
}

func (r *SmsReqV2) GetClientKey() string {
	return r.ClientKey
}

func (r *SmsReqV2) GetEventName() string {
	return r.EventType
}

type RequeueError struct {
	Request SmsReqV2
	Err     error
	Requeue bool
}

type ProviderReq struct {
	Provider entity.Provider
	Request  SmsReqV2
}

func (r *SmsReqV2) Validate() error {
	replacer := strings.NewReplacer("+", "", " ", "")

	if len(r.BulkMessage) > 0 {
		for i, msg := range r.BulkMessage {
			if len(msg.Receiver) <= 0 {
				return errors.NewError("Receiver not found")
			}

			err := v.Validate(&msg.Receiver, is.E164)
			if err != nil {
				return errors.NewError("Invalid mobile number: " + msg.Receiver)
			}

			// remove pluses and spaces
			r.BulkMessage[i].Receiver = replacer.Replace(r.BulkMessage[i].Receiver)
		}
	} else {
		if len(r.Receivers) <= 0 {
			return errors.NewError("Receiver not found")
		}

		for i, mobile := range r.Receivers {

			err := v.Validate(&mobile, is.E164)
			if err != nil {
				return errors.NewError("Invalid mobile number: " + mobile)
			}

			// remove pluses and spaces
			r.Receivers[i] = replacer.Replace(r.Receivers[i])
		}
	}

	return nil
}

type smsHandlerV2 struct {
	logger logger.Logger
	Db     *gorm.DB
}

var singletonSmsHandlerV2 *smsHandlerV2
var lock sync.Mutex

func GetSmsHandlerV2(Db *gorm.DB, logger logger.Logger) Handler {
	if singletonSmsHandlerV2 == nil {
		lock.Lock()
		defer lock.Unlock()
		if singletonSmsHandlerV2 == nil {
			singletonSmsHandlerV2 = &smsHandlerV2{logger: logger, Db: Db}
		}
	}
	return singletonSmsHandlerV2
}

func (m smsHandlerV2) GetRequest(context *gin.Context, messageType MessageType) (Request, *errors.RestErr) {
	var payload SmsReqV2
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

func (m smsHandlerV2) Handle(message *Message) (error, bool) {
	var smsReq SmsReqV2
	err := json.Unmarshal(message.Body, &smsReq)
	if err != nil {
		m.logger.Error("failed to unmarshal message", logger.Extra("jsonUnmarshalError", err.Error()))
		return err, false
	}

	var event entity.Event
	err = m.Db.Model(&entity.Event{}).Where("app_id = ? and name = ?", smsReq.AppId, smsReq.EventType).First(&event).Error
	if err != nil {
		m.logger.Error("failed to get event", logger.Extra("eventError", err.Error()))
		ErrorLog(m.logger, m.Db, &smsReq, nil, 0, "event not found", false)
		return err, false
	}
	if len(event.Variables) > 0 {
		event.VariableAsArray = strings.Split(event.Variables, ",")
	}

	language := "en"
	if len(smsReq.Language) > 0 {
		language = smsReq.Language
	}
	var template entity.Template
	err = m.Db.Model(&entity.Template{}).Where("event_id = ? and type = ? and active = true and language = ?", event.ID, consts.Sms, language).First(&template).Error
	if err != nil {
		m.logger.Error("failed to get template", logger.Extra("templateError", err.Error()))
		ErrorLog(m.logger, m.Db, &smsReq, nil, 0, "template not found", false)
		return err, false
	}

	var providers []entity.Provider
	resultErr := m.Db.Where("app_id = ? and type = ? and active = true", smsReq.AppId, consts.Sms).Order("priority").Find(&providers).Error
	if resultErr != nil {
		m.logger.Error("failed to get provider", logger.Extra("providerError", resultErr.Error()))
		ErrorLog(m.logger, m.Db, &smsReq, nil, 0, "provider not found", false)
		return resultErr, false
	}

	if len(providers) == 0 {
		m.logger.Error("provider not found")
		ErrorLog(m.logger, m.Db, &smsReq, nil, 0, "provider not found", false)
		return errors.NewError("provider not found"), false
	}

	var providerRequest []ProviderReq
	providerRequest = m.ApplyPolicy(providers, message.RequeueCount, smsReq)

	var Requeued []RequeueError
	for _, p := range providerRequest {
		handler := m.getSmsHandler(p)

		if handler != nil {
			err, requeue := handler.Handle(template, p.Request, event)
			if err != nil {
				Requeued = append(Requeued, *requeue)
			}
		}
	}
	if len(Requeued) > 0 {
		requeue := m.SetRequeue(message, Requeued)
		return errors.NewError("Message sent failed"), requeue
	}

	/*provider := providers[int(message.RequeueCount)%len(providers)]
	if provider.ProviderType == consts.SecureMxBd {
		secureMxBd := GetSecureMxBdProvider(provider, m.logger)
		err, requeue := secureMxBd.Handle(template, smsReq, event)
		if err != nil {
			return err, requeue
		}
	} else if provider.ProviderType == consts.SslWireless {
		sslWireless := GetSslWirelessProvider(provider, m.logger)
		err, requeue := sslWireless.Handle(template, smsReq, event)
		if err != nil {
			return err, requeue
		}
	}*/

	return nil, false
}

func (m smsHandlerV2) getSmsHandler(p ProviderReq) SmsHandler {
	var handler SmsHandler
	if p.Provider.ProviderType == consts.Twilio {
		handler = GetTwilioProvider(p.Provider, m.logger, m.Db)
	} else if p.Provider.ProviderType == consts.Vonage {
		handler = GetVonageProvider(p.Provider, m.logger, m.Db)
	}
	return handler
}

func (m smsHandlerV2) CanHandle(messageType MessageType) bool {
	return messageType == Sms
}

func (m smsHandlerV2) SetRequeue(message *Message, requeueErr []RequeueError) bool {
	var smsReq SmsReqV2
	err := json.Unmarshal(message.Body, &smsReq)
	if err != nil {
		m.logger.Error("failed to unmarshal message", logger.Extra("jsonUnmarshalError", err.Error()))
		return false
	}

	var bulkMessage []MultiNotification
	var receivers []string
	Requeue := false
	for _, r := range requeueErr {
		if r.Requeue {
			smsRequest := r.Request
			if len(smsRequest.BulkMessage) > 0 {
				bulkMessage = append(bulkMessage, smsRequest.BulkMessage...)
			} else {
				receivers = append(receivers, smsRequest.Receivers...)
			}
			Requeue = true
		}
	}
	if len(bulkMessage) > 0 {
		smsReq.BulkMessage = bulkMessage
	} else {
		smsReq.Receivers = receivers
	}

	smsContentJson, jsonErr := json.Marshal(smsReq)
	if jsonErr != nil {
		m.logger.Error("failed to marshal request", logger.Extra("jsonMarshalError", jsonErr.Error()))
		return false
	}
	message.Body = smsContentJson

	return Requeue
}

func (m smsHandlerV2) ApplyPolicy(providers []entity.Provider, requeueCount uint, req SmsReqV2) []ProviderReq {
	type ProviderReqArray struct {
		Providers []entity.Provider
		Request   SmsReqV2
	}
	var splitReq []ProviderReqArray
	isBulk := false

	//split receivers
	if len(req.BulkMessage) > 0 {
		isBulk = true
		for _, notification := range req.BulkMessage {
			newReq := SmsReqV2{
				PublishReq:   req.PublishReq,
				AppId:        req.AppId,
				AppName:      req.AppName,
				EventType:    req.EventType,
				DeliveryTime: req.DeliveryTime,
				BulkMessage:  []MultiNotification{notification},
				Language:     req.Language,
			}
			splitReq = append(splitReq, ProviderReqArray{Request: newReq})
		}
	} else {
		for _, receiver := range req.Receivers {
			newReq := SmsReqV2{
				PublishReq:        req.PublishReq,
				AppId:             req.AppId,
				AppName:           req.AppName,
				EventType:         req.EventType,
				DeliveryTime:      req.DeliveryTime,
				Receivers:         []string{receiver},
				TemplateVariables: req.TemplateVariables,
				Language:          req.Language,
			}
			splitReq = append(splitReq, ProviderReqArray{Request: newReq})
		}
	}

	//apply policy
	for _, provider := range providers {
		if len(provider.Policy) <= 2 {
			for i := range splitReq {
				splitReq[i].Providers = append(splitReq[i].Providers, provider)
			}
		} else {
			var policy map[string]interface{}
			if err := json.Unmarshal(provider.Policy, &policy); err != nil {
				m.logger.Error("Json unmarshal failed", logger.Extra("", err.Error()))
			}
			if len(policy) == 0 {
				continue
			}
			for i, smsReq := range splitReq {
				if IsApply(policy, smsReq.Request, isBulk) {
					splitReq[i].Providers = append(splitReq[i].Providers, provider)
				}
			}
		}
	}

	//prioritize
	var prioritized []ProviderReq
	for _, item := range splitReq {
		if len(item.Providers) != 0 {
			prioritized = append(prioritized, ProviderReq{
				Provider: item.Providers[int(requeueCount)%len(item.Providers)],
				Request:  item.Request,
			})
		}
	}

	// combine receivers
	formatted := make(map[uint]ProviderReq)
	for _, item := range prioritized {
		if _, ok := formatted[item.Provider.ID]; ok {
			temp := formatted[item.Provider.ID]
			if isBulk {
				temp.Request.BulkMessage = append(temp.Request.BulkMessage, item.Request.BulkMessage[0])
			} else {
				temp.Request.Receivers = append(temp.Request.Receivers, item.Request.Receivers[0])
			}
			formatted[item.Provider.ID] = temp
		} else {
			formatted[item.Provider.ID] = item
		}
	}

	var answers []ProviderReq
	for _, item := range formatted {
		answers = append(answers, item)
	}
	return answers
}

func IsApply(policy map[string]interface{}, req SmsReqV2, isBulk bool) bool {
	countrySatisfy := false
	if rule, ok := policy["receiver.country"]; ok {
		codes := rule.([]interface{})
		for _, code := range codes {
			if isBulk {
				if code == req.BulkMessage[0].Receiver[:len(code.(string))] {
					countrySatisfy = true
				}
			} else {
				if code == req.Receivers[0][:len(code.(string))] {
					countrySatisfy = true
				}
			}
		}
	} else {
		countrySatisfy = true
	}

	prefixSatisfy := false
	if rule, ok := policy["receiver.prefix"]; ok {
		codes := rule.([]interface{})
		for _, code := range codes {
			if isBulk {
				if code == req.BulkMessage[0].Receiver[:len(code.(string))] {
					prefixSatisfy = true
				}
			} else {
				if code == req.Receivers[0][:len(code.(string))] {
					prefixSatisfy = true
				}
			}
		}
	} else {
		prefixSatisfy = true
	}

	return countrySatisfy && prefixSatisfy
}
