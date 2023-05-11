package service

import (
	"envoyer/config/service_name"
	"envoyer/errors"
	"envoyer/model/entity"
	"envoyer/model/repository"
	"envoyer/model/serializers"
	"github.com/google/uuid"
	"strings"
)

type AppServiceInterface interface {
	Get(id uint) (*entity.Application, *errors.RestErr)
	GetAll() ([]*entity.Application, *errors.RestErr)
	Create(appReq serializers.CreateAppReq) (*entity.Application, *errors.RestErr)
	Update(appReq serializers.UpdateAppReq, appId uint) (*entity.Application, *errors.RestErr)
	Delete(id uint) *errors.RestErr
	GetByAppKey(appKey string) (*entity.Application, *errors.RestErr)
}

type appService struct {
	*BaseService
	appRepo repository.AppRepositoryInterface
}

func NewAppService(baseService *BaseService) AppServiceInterface {
	return &appService{
		BaseService: baseService,
		appRepo:     baseService.container.Get(service_name.AppRepository).(repository.AppRepositoryInterface),
	}
}

func (s appService) Update(appReq serializers.UpdateAppReq, appId uint) (*entity.Application, *errors.RestErr) {
	app := &entity.Application{
		Description: appReq.Description,
		Active:      appReq.Active,
	}

	return s.appRepo.Update(app, appId)
}

func (s appService) GetAll() ([]*entity.Application, *errors.RestErr) {
	return s.appRepo.GetAll()
}

func (s appService) Get(id uint) (*entity.Application, *errors.RestErr) {
	return s.appRepo.Get(id)
}

func (s appService) GetByAppKey(appKey string) (*entity.Application, *errors.RestErr) {
	return s.appRepo.GetByAppKey(appKey)
}

func (s appService) Delete(id uint) *errors.RestErr {
	return s.appRepo.Delete(id)
}

func (s appService) Create(appReq serializers.CreateAppReq) (*entity.Application, *errors.RestErr) {
	appKey := uuid.New().String()
	appKey = strings.Replace(appKey, "-", "", -1)

	app := &entity.Application{
		Name:        appReq.Name,
		Description: appReq.Description,
		AppKey:      appKey,
		Active:      appReq.Active,
	}

	return s.appRepo.Create(app)
}
