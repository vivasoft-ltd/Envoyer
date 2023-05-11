package entity

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Provider struct {
	gorm.Model
	AppId        uint           `json:"app_id"`
	Type         string         `json:"type"`
	ProviderType string         `json:"provider_type"`
	Name         string         `json:"name"`
	Description  string         `json:"description"`
	Config       datatypes.JSON `json:"config"`
	Priority     uint           `json:"priority"`
	Active       bool           `json:"active"`
	Policy       datatypes.JSON `json:"policy"`
}
