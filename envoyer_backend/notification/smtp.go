package notification

import (
	"encoding/json"
	"envoyer/logger"
	"envoyer/model/entity"
	"gopkg.in/mail.v2"
	"gorm.io/gorm"
	"strings"
	"time"
)

type SmtpConfig struct {
	SMTPHost     string `json:"smtp_host"`
	SMTPPort     int    `json:"smtp_port"`
	SMTPUserName string `json:"smtp_user_name"`
	SMTPPassword string `json:"smtp_password"`
	SMTPSender   string `json:"sender"`
}

type Smtp struct {
	provider entity.Provider
	logger   logger.Logger
	Db       *gorm.DB
}

func GetSmtpProvider(provider entity.Provider, logger logger.Logger, db *gorm.DB) *Smtp {
	return &Smtp{provider: provider, logger: logger, Db: db}
}

func (h Smtp) Handle(template entity.Template, emailReq EmailReqV2, event entity.Event) (error, bool) {
	smtpConfig := new(SmtpConfig)
	configJson := h.provider.Config
	if err := json.Unmarshal(configJson, smtpConfig); err != nil {
		h.logger.Error("Json unmarshal failed", logger.Extra("", err.Error()))
		return err, false
	}

	var messages []*mail.Message

	if len(emailReq.BulkMessage) > 0 {
		for _, v := range emailReq.BulkMessage {
			// replace the variables
			var oldNew []string
			for _, variable := range v.TemplateVariables {
				if contains(event.VariableAsArray, variable.Name) {
					oldNew = append(oldNew, variable.Name, variable.Value)
				}
			}
			replacer := strings.NewReplacer(oldNew...)
			replacedSubject := replacer.Replace(template.EmailSubject)
			replacedMessage := replacer.Replace(template.EmailRenderedHTML)

			msg := mail.NewMessage()
			msg.SetHeader("To", v.Receiver)
			if len(emailReq.Sender) != 0 {
				msg.SetHeader("From", emailReq.Sender)
			} else {
				msg.SetHeader("From", smtpConfig.SMTPSender)
			}
			msg.SetHeader("Subject", replacedSubject)
			msg.SetBody("text/html", replacedMessage)

			messages = append(messages, msg)
		}
	} else {
		// replace the variables
		var oldNew []string
		for _, v := range emailReq.TemplateVariables {
			if contains(event.VariableAsArray, v.Name) {
				oldNew = append(oldNew, v.Name, v.Value)
			}
		}
		replacer := strings.NewReplacer(oldNew...)
		replacedSubject := replacer.Replace(template.EmailSubject)
		replacedMessage := replacer.Replace(template.EmailRenderedHTML)

		msg := mail.NewMessage()
		msg.SetHeader("To", emailReq.Receivers...)
		if len(emailReq.Cc) > 0 {
			msg.SetHeader("Cc", emailReq.Cc...)
		}
		if len(emailReq.Bcc) > 0 {
			msg.SetHeader("Bcc", emailReq.Bcc...)
		}
		if len(emailReq.Sender) != 0 {
			msg.SetHeader("From", emailReq.Sender)
		} else {
			msg.SetHeader("From", smtpConfig.SMTPSender)
		}
		msg.SetHeader("Subject", replacedSubject)
		msg.SetBody("text/html", replacedMessage)

		messages = append(messages, msg)
	}

	// Call the DialAndSend() method on the dialer, passing in the message to send. This
	// opens a connection to the SMTP server, sends the message, then closes the
	// connection. If there is a timeout, it will return a "dial tcp: i/o timeout"
	// error
	dialer := mail.NewDialer(smtpConfig.SMTPHost, smtpConfig.SMTPPort, smtpConfig.SMTPUserName, smtpConfig.SMTPPassword)
	dialer.Timeout = 5 * time.Second
	err := dialer.DialAndSend(messages...)
	if err != nil {
		h.logger.Error("Email send failed", logger.Extra("emailDialerError", err.Error()))
		ErrorLog(h.logger, h.Db, &emailReq, messages, h.provider.ID, err.Error(), true)
		return err, true
	}

	return nil, false
}
