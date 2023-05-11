package entity

import (
	"envoyer/errors"
	"gorm.io/gorm"
	"time"
)

type Application struct {
	gorm.Model
	Name        string `json:"name" gorm:"unique"`
	Description string `json:"description"`
	AppKey      string `json:"app_key" gorm:"index:idx_app_key,unique:true"`
	Active      bool   `json:"active" gorm:"default:false"`
}

func (app *Application) AfterDelete(tx *gorm.DB) (err error) {
	err = tx.Where("app_id = ?", app.ID).Delete(&Provider{}).Error
	if err != nil {
		return err
	}

	currentTime := time.Now().String() + "@"
	err = tx.Model(User{}).Where("app_id = ?", app.ID).Update("user_name", gorm.Expr("CONCAT(? , user_name)", currentTime)).Delete(&User{}).Error
	if err != nil {
		return errors.NewInternalServerError("can not delete user", err)
	}

	err = tx.Where("app_id = ?", app.ID).Delete(&Client{}).Error
	if err != nil {
		return err
	}

	var events []uint
	err = tx.Model(Event{}).Select("id").Where("app_id = ?", app.ID).Find(&events).Error
	if err != nil {
		return err
	}
	err = tx.Where("app_id = ?", app.ID).Delete(&Event{}).Error
	if err != nil {
		return err
	}

	err = tx.Where("event_id in ?", events).Delete(&Template{}).Error
	if err != nil {
		return err
	}
	return nil
}
