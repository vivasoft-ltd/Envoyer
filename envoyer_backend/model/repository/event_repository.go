package repository

import (
	"envoyer/errors"
	"envoyer/model/entity"
	"gorm.io/gorm"
)

type EventRepositoryInterface interface {
	Get(id uint) (*entity.Event, *errors.RestErr)
	Create(event *entity.Event) (*entity.Event, *errors.RestErr)
	Update(event *entity.Event, eventId uint) (*entity.Event, *errors.RestErr)
	Delete(id uint) *errors.RestErr
	GetByAppId(appId uint) ([]*entity.Event, *errors.RestErr)
}

type eventRepository struct {
	*BaseRepository
}

func NewEventRepository(baseRepo *BaseRepository) EventRepositoryInterface {
	return &eventRepository{BaseRepository: baseRepo}
}

func (r eventRepository) Get(id uint) (*entity.Event, *errors.RestErr) {
	event := &entity.Event{}
	if err := r.Db.First(&event, id).Error; err != nil {
		return nil, errors.NewInternalServerError("event not found", err)
	}
	return event, nil
}

func (r eventRepository) Create(event *entity.Event) (*entity.Event, *errors.RestErr) {
	if err := r.Db.Create(event).Error; err != nil {
		return nil, errors.NewInternalServerError("can not create event", err)
	}
	return event, nil
}

func (r eventRepository) Update(event *entity.Event, eventId uint) (*entity.Event, *errors.RestErr) {
	err := r.Db.Model(&entity.Event{}).Where("id = ?", eventId).Updates(event).Error
	if err != nil {
		return nil, errors.NewInternalServerError("can not update event", err)
	}
	return r.Get(eventId)
}

func (r eventRepository) Delete(id uint) *errors.RestErr {
	err := r.Db.Where("id = ?", id).Delete(&entity.Event{Model: gorm.Model{ID: id}}).Error
	if err != nil {
		return errors.NewInternalServerError("can not delete event", err)
	}
	return nil
}

func (r eventRepository) GetByAppId(appId uint) ([]*entity.Event, *errors.RestErr) {
	var events []*entity.Event
	err := r.Db.Where("app_id = ?", appId).Find(&events).Error
	if err != nil {
		return nil, errors.NewInternalServerError("events not found", err)
	}
	return events, nil
}
