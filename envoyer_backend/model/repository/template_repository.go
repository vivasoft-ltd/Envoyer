package repository

import (
	"envoyer/errors"
	"envoyer/model/entity"
)

type TemplateRepositoryInterface interface {
	Get(id uint) (*entity.Template, *errors.RestErr)
	Create(template *entity.Template) (*entity.Template, *errors.RestErr)
	Update(template *entity.Template, templateId uint) (*entity.Template, *errors.RestErr)
	Delete(id uint) *errors.RestErr
	GetByEventId(eventId uint) ([]*entity.Template, *errors.RestErr)
}

type templateRepository struct {
	*BaseRepository
}

func NewTemplateRepository(baseRepo *BaseRepository) TemplateRepositoryInterface {
	return &templateRepository{BaseRepository: baseRepo}
}

func (r templateRepository) Get(id uint) (*entity.Template, *errors.RestErr) {
	template := &entity.Template{}
	if err := r.Db.First(&template, id).Error; err != nil {
		return nil, errors.NewInternalServerError("template not found", err)
	}
	return template, nil
}

func (r templateRepository) Create(template *entity.Template) (*entity.Template, *errors.RestErr) {
	if template.Active {
		err := r.Db.Model(&entity.Template{}).Where("event_id = ? and type = ? and language = ?", template.EventId, template.Type, template.Language).Updates(map[string]interface{}{
			"active": false}).Error
		if err != nil {
			return nil, errors.NewInternalServerError("can not update template", err)
		}
	}
	if err := r.Db.Create(template).Error; err != nil {
		return nil, errors.NewInternalServerError("can not create template", err)
	}
	return template, nil
}

func (r templateRepository) Update(template *entity.Template, templateId uint) (*entity.Template, *errors.RestErr) {
	if template.Active {
		err := r.Db.Model(&entity.Template{}).Where("event_id = ? and type = ? and language = ?", template.EventId, template.Type, template.Language).Updates(map[string]interface{}{
			"active": false}).Error
		if err != nil {
			return nil, errors.NewInternalServerError("can not update template", err)
		}
	}
	err := r.Db.Model(&entity.Template{}).Where("id = ?", templateId).Updates(map[string]interface{}{
		"description":         template.Description,
		"message":             template.Message,
		"email_subject":       template.EmailSubject,
		"email_markup":        template.EmailMarkup,
		"email_rendered_html": template.EmailRenderedHTML,
		"title":               template.Title,
		"link":                template.Link,
		"file":                template.File,
		"language":            template.Language,
		"active":              template.Active}).Error
	if err != nil {
		return nil, errors.NewInternalServerError("can not update template", err)
	}

	return r.Get(templateId)
}

func (r templateRepository) Delete(id uint) *errors.RestErr {
	err := r.Db.Where("id = ?", id).Delete(&entity.Template{}).Error
	if err != nil {
		return errors.NewInternalServerError("can not delete template", err)
	}
	return nil
}

func (r templateRepository) GetByEventId(eventId uint) ([]*entity.Template, *errors.RestErr) {
	var templates []*entity.Template
	err := r.Db.Where("event_id = ?", eventId).Find(&templates).Error
	if err != nil {
		return nil, errors.NewInternalServerError("templates not found", err)
	}
	return templates, nil
}
