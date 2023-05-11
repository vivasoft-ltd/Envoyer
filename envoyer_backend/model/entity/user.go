package entity

import "gorm.io/gorm"

type User struct {
	gorm.Model
	UserName string `json:"user_name" gorm:"unique"`
	Password string `json:"password"`
	AppId    uint   `json:"app_id"`
	Role     string `json:"role"`
}
