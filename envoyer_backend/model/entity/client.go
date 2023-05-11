package entity

import "gorm.io/gorm"

type Client struct {
	gorm.Model
	Name        string `json:"name"`
	Description string `json:"description"`
	ClientKey   string `json:"client_key" gorm:"unique"`
	AppId       uint   `json:"app_id"`
}
