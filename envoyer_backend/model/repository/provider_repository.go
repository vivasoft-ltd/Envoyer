package repository

import (
	"envoyer/errors"
	"envoyer/model/entity"
	"envoyer/model/serializers"
)

type ProviderRepositoryInterface interface {
	Get(id uint) (*entity.Provider, *errors.RestErr)
	Create(provider *entity.Provider) (*entity.Provider, *errors.RestErr)
	Update(provider *entity.Provider, providerId uint) (*entity.Provider, *errors.RestErr)
	Delete(id uint) *errors.RestErr
	GetByAppId(appId uint) ([]*entity.Provider, *errors.RestErr)
	GetByAppIdAndType(appId uint, Type string) ([]*entity.Provider, *errors.RestErr)
	GetByAppIdAndTypeWithTopPriority(appId uint, Type string) (*entity.Provider, *errors.RestErr)
	UpdatePriority(appId uint, Type string, priorityReq serializers.UpdatePriorityReq) ([]*entity.Provider, *errors.RestErr)
}

type providerRepository struct {
	*BaseRepository
}

func NewProviderRepository(baseRepo *BaseRepository) ProviderRepositoryInterface {
	return &providerRepository{BaseRepository: baseRepo}
}

func (r providerRepository) UpdatePriority(appId uint, Type string, priorityReq serializers.UpdatePriorityReq) ([]*entity.Provider, *errors.RestErr) {
	for _, provider := range priorityReq.ProviderPriority {
		err := r.Db.Model(&entity.Provider{}).Where("id = ?", provider.Id).Update("priority", provider.Priority).Error
		if err != nil {
			return nil, errors.NewInternalServerError("can not change priority provider", err)
		}
	}
	return r.GetByAppIdAndType(appId, Type)
}

func (r providerRepository) Get(id uint) (*entity.Provider, *errors.RestErr) {
	provider := &entity.Provider{}
	if err := r.Db.First(&provider, id).Error; err != nil {
		return nil, errors.NewInternalServerError("provider not found", err)
	}
	return provider, nil
}

func (r providerRepository) Create(provider *entity.Provider) (*entity.Provider, *errors.RestErr) {
	if err := r.Db.Create(provider).Error; err != nil {
		return nil, errors.NewInternalServerError("can not create provider", err)
	}
	return provider, nil
}

func (r providerRepository) Update(provider *entity.Provider, providerId uint) (*entity.Provider, *errors.RestErr) {
	err := r.Db.Model(&entity.Provider{}).Where("id = ?", providerId).Updates(provider).Updates(map[string]interface{}{"active": provider.Active}).Error
	if err != nil {
		return nil, errors.NewInternalServerError("can not update provider", err)
	}
	err = r.Db.Model(&entity.Provider{}).Where("id = ?", providerId).Updates(map[string]interface{}{"active": provider.Active}).Error
	if err != nil {
		return nil, errors.NewInternalServerError("can not update provider", err)
	}
	return r.Get(providerId)
}

func (r providerRepository) Delete(id uint) *errors.RestErr {
	err := r.Db.Where("id = ?", id).Delete(&entity.Provider{}).Error
	if err != nil {
		return errors.NewInternalServerError("can not delete provider", err)
	}
	return nil
}

func (r providerRepository) GetByAppId(appId uint) ([]*entity.Provider, *errors.RestErr) {
	var providers []*entity.Provider
	err := r.Db.Where("app_id = ?", appId).Find(&providers).Error
	if err != nil {
		return nil, errors.NewInternalServerError("providers not found", err)
	}
	return providers, nil
}

func (r providerRepository) GetByAppIdAndType(appId uint, Type string) ([]*entity.Provider, *errors.RestErr) {
	var providers []*entity.Provider
	err := r.Db.Where("app_id = ? and type = ?", appId, Type).Order("priority").Find(&providers).Error
	if err != nil {
		return nil, errors.NewInternalServerError("providers not found", err)
	}
	return providers, nil
}

func (r providerRepository) GetByAppIdAndTypeWithTopPriority(appId uint, Type string) (*entity.Provider, *errors.RestErr) {
	var provider entity.Provider
	err := r.Db.Where("app_id = ? and type = ? and active = ?", appId, Type, true).Order("priority").Find(&provider).Limit(1).Error
	if err != nil {
		return nil, errors.NewInternalServerError("providers not found", err)
	}
	return &provider, nil
}
