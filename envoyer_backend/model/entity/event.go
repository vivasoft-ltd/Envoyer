package entity

import (
	"gorm.io/gorm"
)

type Event struct {
	gorm.Model
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	Variables       string   `json:"-"`
	AppId           uint     `json:"app_id"`
	VariableAsArray []string `json:"variables" gorm:"-"`
}

func (event *Event) AfterDelete(tx *gorm.DB) (err error) {
	err = tx.Where("event_id = ?", event.ID).Delete(&Template{}).Error
	if err != nil {
		return err
	}
	return nil
}
