package service

import (
	"envoyer/config/service_name"
	"envoyer/errors"
	"envoyer/model/entity"
	"envoyer/model/repository"
	"envoyer/model/serializers"
	"strings"
)

type EventServiceInterface interface {
	Get(id uint) (*entity.Event, *errors.RestErr)
	Create(eventReq serializers.CreateEventReq) (*entity.Event, *errors.RestErr)
	Update(eventReq serializers.UpdateEventReq, eventId uint) (*entity.Event, *errors.RestErr)
	Delete(id uint) *errors.RestErr
	GetByAppId(appId uint) ([]*entity.Event, *errors.RestErr)
}

type eventService struct {
	*BaseService
	eventRepo repository.EventRepositoryInterface
}

func NewEventService(baseService *BaseService) EventServiceInterface {
	return &eventService{
		BaseService: baseService,
		eventRepo:   baseService.container.Get(service_name.EventRepository).(repository.EventRepositoryInterface),
	}
}

func (s eventService) Get(id uint) (*entity.Event, *errors.RestErr) {
	event, err := s.eventRepo.Get(id)
	if err != nil {
		return nil, err
	}
	if len(event.Variables) > 0 {
		event.VariableAsArray = strings.Split(event.Variables, ",")
	}
	return event, nil
}

func (s eventService) Create(eventReq serializers.CreateEventReq) (*entity.Event, *errors.RestErr) {
	event := &entity.Event{
		Name:        eventReq.Name,
		Description: eventReq.Description,
		AppId:       eventReq.AppId,
		Variables:   strings.Join(eventReq.Variables, ","),
	}

	eventResp, err := s.eventRepo.Create(event)
	if err != nil {
		return nil, err
	}
	if len(eventResp.Variables) > 0 {
		eventResp.VariableAsArray = strings.Split(eventResp.Variables, ",")
	}
	return eventResp, nil
}

func (s eventService) Update(eventReq serializers.UpdateEventReq, eventId uint) (*entity.Event, *errors.RestErr) {
	event := &entity.Event{
		Name:        eventReq.Name,
		Description: eventReq.Description,
		Variables:   strings.Join(eventReq.Variables, ","),
	}

	eventResp, err := s.eventRepo.Update(event, eventId)
	if err != nil {
		return nil, err
	}

	if len(eventResp.Variables) > 0 {
		eventResp.VariableAsArray = strings.Split(eventResp.Variables, ",")
	}
	return eventResp, nil
}

func (s eventService) Delete(id uint) *errors.RestErr {
	return s.eventRepo.Delete(id)
}

func (s eventService) GetByAppId(appId uint) ([]*entity.Event, *errors.RestErr) {
	events, err := s.eventRepo.GetByAppId(appId)
	if err != nil {
		return nil, err
	}
	for i, event := range events {
		if len(event.Variables) > 0 {
			events[i].VariableAsArray = strings.Split(event.Variables, ",")
		}
	}
	return events, nil
}
