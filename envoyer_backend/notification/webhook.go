package notification

import (
	"encoding/json"
	"envoyer/errors"
	"envoyer/logger"
	"envoyer/model/entity"
	"fmt"
	"github.com/go-resty/resty/v2"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"
)

type WebhookConfig struct {
	Url   string `json:"url"`
	Token string `json:"token"`
}

type CustomWebhook struct {
	provider entity.Provider
	client   *resty.Client
	logger   logger.Logger
	Db       *gorm.DB
}

func GetWebhookProvider(provider entity.Provider, logger logger.Logger, db *gorm.DB) *CustomWebhook {
	return &CustomWebhook{provider: provider, logger: logger, client: resty.New(), Db: db}
}

func (h CustomWebhook) Handle(webhookReq WebhookReq) (error, bool) {
	fmt.Println("Handling")
	webhookConfig := new(WebhookConfig)
	configJson := h.provider.Config
	if err := json.Unmarshal(configJson, webhookConfig); err != nil {
		h.logger.Error("Json unmarshal failed", logger.Extra("", err.Error()))
		return err, false
	}

	Headers := make(map[string]string)
	if len(webhookConfig.Token) != 0 {
		Headers["Authorization"] = fmt.Sprintf("Bearer %s", webhookConfig.Token)
	}
	Headers["Content-Type"] = "application/json"
	resp, err := h.client.SetTimeout(5 * time.Second).R().SetHeaders(Headers).SetBody(webhookReq.Data).Post(webhookConfig.Url)
	fmt.Println(resp)
	if err != nil {
		h.logger.Error("error in sending resty request", logger.Extra("restyError", err.Error()))
		ErrorLog(h.logger, h.Db, &webhookReq, webhookReq.Data, h.provider.ID, err.Error(), true)
		return err, true
	}
	if resp.StatusCode() != http.StatusOK {
		h.logger.Error("Message send failed", logger.Extra("statusCode", resp.StatusCode()))
		ErrorLog(h.logger, h.Db, &webhookReq, webhookReq.Data, h.provider.ID, "message send failed. Status code : "+strconv.Itoa(resp.StatusCode()), true)
		return errors.NewError("Message send failed"), true
	}
	return nil, false
}
