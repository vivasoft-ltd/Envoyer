package notification

import (
	"context"
	"envoyer/logger"
	"envoyer/model/entity"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/option"
	"gorm.io/gorm"
	"strings"
)

type FirebaseFCM struct {
	provider entity.Provider
	logger   logger.Logger
	Db       *gorm.DB
}

func GetFirebaseFCMProvider(provider entity.Provider, logger logger.Logger, db *gorm.DB) *FirebaseFCM {
	return &FirebaseFCM{provider: provider, logger: logger, Db: db}
}

func (h FirebaseFCM) Handle(template entity.Template, pushReq PushReq, event entity.Event) (error, bool) {
	configJson := h.provider.Config

	templateMsg := template.Message
	templateTitle := template.Title
	var templateImage string
	if len(pushReq.ImageUrl) > 0 {
		templateImage = pushReq.ImageUrl
	} else {
		templateImage = template.File
	}

	opts := []option.ClientOption{option.WithCredentialsJSON(configJson)}

	app, err := firebase.NewApp(context.Background(), nil, opts...)
	if err != nil {
		h.logger.Error("error in creating firebase app", logger.Extra("firebaseError", err.Error()))
		ErrorLog(h.logger, h.Db, pushReq, nil, h.provider.ID, err.Error(), false)
		return nil, false
	}

	fcmClient, err := app.Messaging(context.Background())
	if err != nil {
		h.logger.Error("error in creating fcm client", logger.Extra("fcmClientError", err.Error()))
		ErrorLog(h.logger, h.Db, pushReq, nil, h.provider.ID, err.Error(), false)
		return nil, false
	}

	if len(pushReq.BulkMessage) > 0 {
		var messages []*messaging.Message
		for _, v := range pushReq.BulkMessage {
			// replace the variables
			var oldNew []string
			for _, variable := range v.TemplateVariables {
				if contains(event.VariableAsArray, variable.Name) {
					oldNew = append(oldNew, variable.Name, variable.Value)
				}
			}
			replacer := strings.NewReplacer(oldNew...)
			replacedMessage := replacer.Replace(templateMsg)
			replacedTitle := replacer.Replace(templateTitle)
			replacedImage := replacer.Replace(templateImage)

			msg := &messaging.Message{
				Notification: &messaging.Notification{
					Title: replacedTitle,
					Body:  replacedMessage,
				},
				Token: v.Receiver,
			}
			if len(replacedImage) != 0 {
				msg.Notification.ImageURL = replacedImage
			}
			if len(pushReq.Data) != 0 {
				msg.Data = pushReq.Data
			}

			messages = append(messages, msg)
		}
		response, err := fcmClient.SendAll(context.Background(), messages)
		if err != nil {
			h.logger.Error("error in sending firebase fcm push", logger.Extra("firebaseSendError", err.Error()), logger.Extra("response", response))
			requeue := true
			if strings.Contains(err.Error(), "registration-token-not-registered") || strings.Contains(err.Error(), "exactly one of token, topic or condition must be specified") {
				requeue = false
			}
			ErrorLog(h.logger, h.Db, pushReq, messages, h.provider.ID, err.Error(), requeue)
			return err, requeue
		}
	} else {
		// replace the variables
		var oldNew []string
		for _, v := range pushReq.TemplateVariables {
			if contains(event.VariableAsArray, v.Name) {
				oldNew = append(oldNew, v.Name, v.Value)
			}
		}
		replacer := strings.NewReplacer(oldNew...)
		replacedMessage := replacer.Replace(templateMsg)
		replacedTitle := replacer.Replace(templateTitle)
		replacedImage := replacer.Replace(templateImage)

		if len(pushReq.Receivers) <= 1 {
			message := &messaging.Message{
				Notification: &messaging.Notification{
					Title: replacedTitle,
					Body:  replacedMessage,
				},
			}
			if len(pushReq.Receivers) != 0 {
				message.Token = pushReq.Receivers[0]
			} else if len(pushReq.Topic) != 0 {
				message.Topic = pushReq.Topic
			}
			if len(pushReq.Data) != 0 {
				message.Data = pushReq.Data
			}
			if len(replacedImage) != 0 {
				message.Notification.ImageURL = replacedImage
			}
			if len(pushReq.Condition) != 0 {
				message.Condition = pushReq.Condition
			}

			response, err := fcmClient.Send(context.Background(), message)
			if err != nil {
				h.logger.Error("error in sending firebase fcm push", logger.Extra("firebaseSendError", err.Error()), logger.Extra("response", response))
				requeue := true
				if strings.Contains(err.Error(), "registration-token-not-registered") || strings.Contains(err.Error(), "exactly one of token, topic or condition must be specified") {
					requeue = false
				}
				ErrorLog(h.logger, h.Db, pushReq, message, h.provider.ID, err.Error(), requeue)
				return err, requeue
			}
		} else {
			message := &messaging.MulticastMessage{
				Notification: &messaging.Notification{
					Title: replacedTitle,
					Body:  replacedMessage,
				},
				Tokens: pushReq.Receivers,
			}
			if len(pushReq.Data) != 0 {
				message.Data = pushReq.Data
			}
			if len(replacedImage) != 0 {
				message.Notification.ImageURL = replacedImage
			}
			response, err := fcmClient.SendMulticast(context.Background(), message)
			if err != nil {
				h.logger.Error("error in sending firebase fcm push", logger.Extra("firebaseSendError", err.Error()), logger.Extra("response", response))
				requeue := true
				if strings.Contains(err.Error(), "registration-token-not-registered") ||
					strings.Contains(err.Error(), "exactly one of token, topic or condition must be specified") ||
					strings.Contains(err.Error(), "invalid-argument") {
					requeue = false
				}
				ErrorLog(h.logger, h.Db, pushReq, message, h.provider.ID, err.Error(), requeue)
				return err, requeue
			}
		}
	}

	return nil, false
}
