package service

import (
	"envoyer/config/service_name"
	"envoyer/errors"
	"envoyer/model/entity"
	"envoyer/model/repository"
	"envoyer/model/serializers"
)

type TemplateServiceInterface interface {
	Get(id uint) (*entity.Template, *errors.RestErr)
	Create(templateReq serializers.CreateTemplateReq) (*entity.Template, *errors.RestErr)
	Update(templateReq serializers.UpdateTemplateReq, templateId uint) (*entity.Template, *errors.RestErr)
	Delete(id uint) *errors.RestErr
	GetByEventId(eventId uint) ([]*entity.Template, *errors.RestErr)
}

type templateService struct {
	*BaseService
	templateRepo repository.TemplateRepositoryInterface
}

func NewTemplateService(baseService *BaseService) TemplateServiceInterface {
	return &templateService{
		BaseService:  baseService,
		templateRepo: baseService.container.Get(service_name.TemplateRepository).(repository.TemplateRepositoryInterface),
	}
}

func (s templateService) Get(id uint) (*entity.Template, *errors.RestErr) {
	return s.templateRepo.Get(id)
}

func (s templateService) Create(templateReq serializers.CreateTemplateReq) (*entity.Template, *errors.RestErr) {
	template := &entity.Template{
		Type:              templateReq.Type,
		Description:       templateReq.Description,
		Message:           templateReq.Message,
		EmailSubject:      templateReq.EmailSubject,
		EmailMarkup:       templateReq.EmailMarkup,
		EmailRenderedHTML: templateReq.EmailRenderedHTML,
		EventId:           templateReq.EventId,
		Active:            templateReq.Active,
		Title:             templateReq.Title,
		Link:              templateReq.Link,
		File:              templateReq.File,
		Language:          templateReq.Language,
	}
	return s.templateRepo.Create(template)
}

func (s templateService) Update(templateReq serializers.UpdateTemplateReq, templateId uint) (*entity.Template, *errors.RestErr) {
	template := &entity.Template{
		EventId:           templateReq.EventId,
		Type:              templateReq.Type,
		Description:       templateReq.Description,
		Message:           templateReq.Message,
		EmailSubject:      templateReq.EmailSubject,
		EmailMarkup:       templateReq.EmailMarkup,
		EmailRenderedHTML: templateReq.EmailRenderedHTML,
		Active:            templateReq.Active,
		Title:             templateReq.Title,
		Link:              templateReq.Link,
		File:              templateReq.File,
		Language:          templateReq.Language,
	}
	return s.templateRepo.Update(template, templateId)
}

func (s templateService) Delete(id uint) *errors.RestErr {
	return s.templateRepo.Delete(id)
}

func (s templateService) GetByEventId(eventId uint) ([]*entity.Template, *errors.RestErr) {
	return s.templateRepo.GetByEventId(eventId)
}
