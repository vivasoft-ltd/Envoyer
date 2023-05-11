package notification

import (
	"encoding/json"
	"envoyer/errors"
	"envoyer/logger"
	"envoyer/model/entity"
	"github.com/vonage/vonage-go-sdk"
	"gorm.io/gorm"
	"gorm.io/gorm/utils"
	"strings"
)

type VonageConfig struct {
	ApiKey    string `json:"api_key"`
	ApiSecret string `json:"api_secret"`
	SenderId  string `json:"sender_id"`
}

type VonageSms struct {
	Body string `json:"body"`
	From string `json:"from"`
	To   string `json:"to"`
}

type Vonage struct {
	provider entity.Provider
	logger   logger.Logger
	Db       *gorm.DB
}

func GetVonageProvider(provider entity.Provider, logger logger.Logger, db *gorm.DB) *Vonage {
	return &Vonage{provider: provider, logger: logger, Db: db}
}

func (h Vonage) Handle(template entity.Template, smsReq SmsReqV2, event entity.Event) (error, *RequeueError) {
	smsConfig := new(VonageConfig)
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
	var messages []VonageSms
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

			messageParam := VonageSms{
				Body: replacedMessage,
				To:   v.Receiver,
			}
			if len(smsReq.Sender) != 0 {
				messageParam.From = smsReq.Sender
			} else {
				messageParam.From = smsConfig.SenderId
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
			messageParam := VonageSms{
				Body: replacedMessage,
				To:   receiver,
			}
			if len(smsReq.Sender) != 0 {
				messageParam.From = smsReq.Sender
			} else {
				messageParam.From = smsConfig.SenderId
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

	auth := vonage.CreateAuthFromKeySecret(smsConfig.ApiKey, smsConfig.ApiSecret)
	client := vonage.NewSMSClient(auth)

	var failedNumber []string

	for _, message := range messages {
		resp, errResp, err := client.Send(message.From, message.To, message.Body, vonage.SMSOpts{Type: "unicode"})
		if err != nil {
			h.logger.Error("error in sending sms", logger.Extra("vonageError", err.Error()))
			ErrorLog(h.logger, h.Db, &smsReq, message, h.provider.ID, err.Error(), true)
			failedNumber = append(failedNumber, message.To)
		} else {
			h.logger.Debug("", logger.Extra("vonageResponse", resp))
			if resp.Messages[0].Status != "0" {
				h.logger.Error("error in sending sms", logger.Extra("vonageError", errResp))
				ErrorLog(h.logger, h.Db, &smsReq, message, h.provider.ID, errResp.Messages[0].ErrorText, false)
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
