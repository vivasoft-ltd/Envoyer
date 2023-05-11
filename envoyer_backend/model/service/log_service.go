package service

import (
	"envoyer/config/service_name"
	"envoyer/errors"
	"envoyer/model/entity"
	"envoyer/model/repository"
)

type LogServiceInterface interface {
	Delete(id uint) *errors.RestErr
	GetByAppId(appId uint) ([]*entity.ErrorLog, *errors.RestErr)
}

type logService struct {
	*BaseService
	logRepo repository.LogRepositoryInterface
}

func NewLogService(baseService *BaseService) LogServiceInterface {
	return &logService{
		BaseService: baseService,
		logRepo:     baseService.container.Get(service_name.LogRepository).(repository.LogRepositoryInterface),
	}
}

func (s logService) Delete(id uint) *errors.RestErr {
	return s.logRepo.Delete(id)
}

func (s logService) GetByAppId(appId uint) ([]*entity.ErrorLog, *errors.RestErr) {
	return s.logRepo.GetByAppId(appId)
}
