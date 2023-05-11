package service

import (
	"envoyer/config/service_name"
	"envoyer/errors"
	"envoyer/model/entity"
	"envoyer/model/repository"
	"envoyer/model/serializers"
)

type ProviderServiceInterface interface {
	Get(id uint) (*entity.Provider, *errors.RestErr)
	Create(providerReq serializers.CreateProviderReq) (*entity.Provider, *errors.RestErr)
	Update(providerReq serializers.UpdateProviderReq, providerId uint) (*entity.Provider, *errors.RestErr)
	Delete(id uint) *errors.RestErr
	GetByAppId(appId uint) ([]*entity.Provider, *errors.RestErr)
	GetByAppIdAndType(appId uint, Type string) ([]*entity.Provider, *errors.RestErr)
	GetByAppIdAndTypeWithTopPriority(appId uint, Type string) (*entity.Provider, *errors.RestErr)
	UpdatePriority(appId uint, Type string, priorityReq serializers.UpdatePriorityReq) ([]*entity.Provider, *errors.RestErr)
}

type providerService struct {
	*BaseService
	providerRepo repository.ProviderRepositoryInterface
}

func NewProviderService(baseService *BaseService) ProviderServiceInterface {
	return &providerService{
		BaseService:  baseService,
		providerRepo: baseService.container.Get(service_name.ProviderRepository).(repository.ProviderRepositoryInterface),
	}
}

func (s providerService) UpdatePriority(appId uint, Type string, priorityReq serializers.UpdatePriorityReq) ([]*entity.Provider, *errors.RestErr) {
	return s.providerRepo.UpdatePriority(appId, Type, priorityReq)
}

func (s providerService) GetByAppId(appId uint) ([]*entity.Provider, *errors.RestErr) {
	return s.providerRepo.GetByAppId(appId)
}

func (s providerService) GetByAppIdAndType(appId uint, Type string) ([]*entity.Provider, *errors.RestErr) {
	return s.providerRepo.GetByAppIdAndType(appId, Type)
}

func (s providerService) GetByAppIdAndTypeWithTopPriority(appId uint, Type string) (*entity.Provider, *errors.RestErr) {
	return s.providerRepo.GetByAppIdAndTypeWithTopPriority(appId, Type)
}

func (s providerService) Get(id uint) (*entity.Provider, *errors.RestErr) {
	return s.providerRepo.Get(id)
}

func (s providerService) Create(providerReq serializers.CreateProviderReq) (*entity.Provider, *errors.RestErr) {
	provider := &entity.Provider{
		AppId:        providerReq.AppId,
		Type:         providerReq.Type,
		ProviderType: providerReq.ProviderType,
		Name:         providerReq.Name,
		Description:  providerReq.Description,
		Config:       providerReq.Config,
		Priority:     providerReq.Priority,
		Active:       providerReq.Active,
		Policy:       providerReq.Policy,
	}
	return s.providerRepo.Create(provider)
}

func (s providerService) Update(providerReq serializers.UpdateProviderReq, providerId uint) (*entity.Provider, *errors.RestErr) {
	provider := &entity.Provider{
		Name:        providerReq.Name,
		Description: providerReq.Description,
		Config:      providerReq.Config,
		Priority:    providerReq.Priority,
		Active:      providerReq.Active,
		Policy:      providerReq.Policy,
	}
	return s.providerRepo.Update(provider, providerId)
}

func (s providerService) Delete(id uint) *errors.RestErr {
	return s.providerRepo.Delete(id)
}
