package notification

import (
	"encoding/json"
	"envoyer/errors"
	"envoyer/logger"
	"envoyer/model/entity"
	"github.com/twilio/twilio-go"
	api "github.com/twilio/twilio-go/rest/api/v2010"
	"gorm.io/gorm"
	"gorm.io/gorm/utils"
	"strings"
)

type TwilioConfig struct {
	AccountSid string `json:"account_sid"`
	AuthToken  string `json:"auth_token"`
	SenderId   string `json:"sender_id"`
}

type Twilio struct {
	provider entity.Provider
	logger   logger.Logger
	Db       *gorm.DB
}

func GetTwilioProvider(provider entity.Provider, logger logger.Logger, db *gorm.DB) *Twilio {
	return &Twilio{provider: provider, logger: logger, Db: db}
}

func (h Twilio) Handle(template entity.Template, smsReq SmsReqV2, event entity.Event) (error, *RequeueError) {
	smsConfig := new(TwilioConfig)
	configJson := h.provider.Config
	if err := json.Unmarshal(configJson, smsConfig); err != nil {
		h.logger.Error("Json unmarshal failed", logger.Extra("", err.Error()))
		return err, &RequeueError{
			Request: smsReq,
			Err:     err,
			Requeue: false,
		}
	}

	templateMsg := template.Message
	var messages []*api.CreateMessageParams
	isBulk := false

	if len(smsReq.BulkMessage) > 0 {
		isBulk = true
		for _, v := range smsReq.BulkMessage {
			// replace the variables
			var oldNew []string
			for _, variable := range v.TemplateVariables {
				if contains(event.VariableAsArray, variable.Name) {
					oldNew = append(oldNew, variable.Name, variable.Value)
				}
			}
			replacer := strings.NewReplacer(oldNew...)
			replacedMessage := replacer.Replace(templateMsg)

			messageParam := &api.CreateMessageParams{}
			messageParam.SetTo("+" + v.Receiver)
			messageParam.SetBody(replacedMessage)
			if len(smsReq.Sender) != 0 {
				messageParam.SetFrom(smsReq.Sender)
			} else {
				messageParam.SetFrom(smsConfig.SenderId)
			}
			messages = append(messages, messageParam)
		}
	} else {
		// replace the variables
		var oldNew []string
		for _, v := range smsReq.TemplateVariables {
			if contains(event.VariableAsArray, v.Name) {
				oldNew = append(oldNew, v.Name, v.Value)
			}
		}
		replacer := strings.NewReplacer(oldNew...)
		replacedMessage := replacer.Replace(templateMsg)

		for _, receiver := range smsReq.Receivers {
			messageParam := &api.CreateMessageParams{}
			messageParam.SetTo("+" + receiver)
			messageParam.SetBody(replacedMessage)
			if len(smsReq.Sender) != 0 {
				messageParam.SetFrom(smsReq.Sender)
			} else {
				messageParam.SetFrom(smsConfig.SenderId)
			}
			messages = append(messages, messageParam)
		}
	}
	if len(messages) == 0 {
		return nil, &RequeueError{
			Request: smsReq,
			Err:     errors.NewError("no message found"),
			Requeue: false,
		}
	}

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: smsConfig.AccountSid,
		Password: smsConfig.AuthToken,
	})

	var failedNumber []string

	for _, message := range messages {
		resp, err := client.Api.CreateMessage(message)
		if err != nil {
			h.logger.Error("error in sending sms", logger.Extra("twilioError", err.Error()))
			ErrorLog(h.logger, h.Db, &smsReq, message, h.provider.ID, err.Error(), true)
			failedNumber = append(failedNumber, strings.Replace(*message.To, "+", "", -1))
		} else {
			h.logger.Debug("", logger.Extra("twilioResponse", resp))
			if resp.ErrorCode != nil {
				h.logger.Error("error in sending sms", logger.Extra("twilioError", resp.ErrorMessage))
				ErrorLog(h.logger, h.Db, &smsReq, message, h.provider.ID, err.Error(), false)
			}
		}
	}

	if len(failedNumber) > 0 {
		if isBulk {
			var bulkEntry []MultiNotification
			for _, item := range smsReq.BulkMessage {
				if utils.Contains(failedNumber, item.Receiver) {
					bulkEntry = append(bulkEntry, item)
				}
			}
			smsReq.BulkMessage = bulkEntry
		} else {
			smsReq.Receivers = failedNumber
		}
		return errors.NewError("failed to send some sms"), &RequeueError{
			Request: smsReq,
			Err:     errors.NewError("failed to send some sms"),
			Requeue: true,
		}
	}

	return nil, &RequeueError{
		Request: smsReq,
		Err:     nil,
		Requeue: false,
	}
}
