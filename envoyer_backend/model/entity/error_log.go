package entity

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

type ErrorLog struct {
	gorm.Model
	AppId      uint           `json:"app_id,omitempty"`
	ClientKey  string         `json:"client_key,omitempty"`
	EventName  string         `json:"event_name,omitempty"`
	ProviderId uint           `json:"provider_id,omitempty"`
	Message    string         `json:"message,omitempty"`
	Data       datatypes.JSON `json:"data,omitempty"`
	Request    datatypes.JSON `json:"request,omitempty"`
	Type       string         `json:"type,omitempty"`
	Date       *time.Time     `json:"date,omitempty"`
	IsRequeue  bool           `json:"is_requeue,omitempty"`
}
