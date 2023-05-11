package repository

import (
	"envoyer/errors"
	"envoyer/model/entity"
	"gorm.io/gorm"
	"time"
)

type AppRepositoryInterface interface {
	Get(id uint) (*entity.Application, *errors.RestErr)
	GetAll() ([]*entity.Application, *errors.RestErr)
	Create(application *entity.Application) (*entity.Application, *errors.RestErr)
	Update(application *entity.Application, appId uint) (*entity.Application, *errors.RestErr)
	Delete(id uint) *errors.RestErr
	GetByAppKey(appKey string) (*entity.Application, *errors.RestErr)
}

type appRepository struct {
	*BaseRepository
}

func NewAppRepository(baseRepo *BaseRepository) AppRepositoryInterface {
	return &appRepository{BaseRepository: baseRepo}
}

func (r appRepository) Get(id uint) (*entity.Application, *errors.RestErr) {
	app := &entity.Application{}
	if err := r.Db.First(&app, id).Error; err != nil {
		return nil, errors.NewInternalServerError("app not found", err)
	}
	return app, nil
}

func (r appRepository) GetAll() ([]*entity.Application, *errors.RestErr) {
	var apps []*entity.Application
	if err := r.Db.Model(&entity.Application{}).Find(&apps).Error; err != nil {
		return nil, errors.NewInternalServerError("app list not found", err)
	}
	return apps, nil
}

func (r appRepository) Create(application *entity.Application) (*entity.Application, *errors.RestErr) {
	if err := r.Db.Create(application).Error; err != nil {
		return nil, errors.NewInternalServerError("can not create app", err)
	}
	return application, nil
}

func (r appRepository) Update(application *entity.Application, appId uint) (*entity.Application, *errors.RestErr) {
	err := r.Db.Model(&entity.Application{}).Where("id = ?", appId).Updates(map[string]interface{}{"description": application.Description, "active": application.Active}).Error
	if err != nil {
		return nil, errors.NewInternalServerError("can not update app", err)
	}
	return r.Get(appId)
}

func (r appRepository) Delete(id uint) *errors.RestErr {
	time := time.Now().String() + "@"
	err := r.Db.Model(entity.Application{}).Where("id = ?", id).Update("name", gorm.Expr("CONCAT(? , name)", time)).
		Delete(&entity.Application{Model: gorm.Model{ID: id}}).Error
	if err != nil {
		return errors.NewInternalServerError("can not delete app", err)
	}
	return nil
}

func (r appRepository) GetByAppKey(appKey string) (*entity.Application, *errors.RestErr) {
	var app *entity.Application
	err := r.Db.Where("app_key = ?", appKey).Find(&app).Error
	if err != nil {
		return nil, errors.NewInternalServerError("app not found", err)
	}
	return app, nil
}
