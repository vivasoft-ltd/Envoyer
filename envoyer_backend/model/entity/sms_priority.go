package entity

import "gorm.io/gorm"

type SmsPriority struct {
	gorm.Model
	AppId       uint `json:"app_id"`
	CountryCode uint `json:"country_code"`
	Priority    uint `json:"priority"`
	ProviderId  uint `json:"provider_id"`
}
