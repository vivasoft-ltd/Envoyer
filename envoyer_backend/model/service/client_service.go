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

type ClientServiceInterface interface {
	Get(id uint) (*entity.Client, *errors.RestErr)
	Create(clientReq serializers.CreateClientReq) (*entity.Client, *errors.RestErr)
	Update(clientReq serializers.UpdateClientReq, clientId uint) (*entity.Client, *errors.RestErr)
	Delete(id uint) *errors.RestErr
	GetByAppId(appId uint) ([]*entity.Client, *errors.RestErr)
	GetByClientKey(clientKey string) (*entity.Client, *errors.RestErr)
}

type clientService struct {
	*BaseService
	clientRepo repository.ClientRepositoryInterface
}

func NewClientService(baseService *BaseService) ClientServiceInterface {
	return &clientService{
		BaseService: baseService,
		clientRepo:  baseService.container.Get(service_name.ClientRepository).(repository.ClientRepositoryInterface),
	}
}

func (s clientService) Get(id uint) (*entity.Client, *errors.RestErr) {
	return s.clientRepo.Get(id)
}

func (s clientService) Create(clientReq serializers.CreateClientReq) (*entity.Client, *errors.RestErr) {
	clientKey := uuid.New().String()
	clientKey = strings.Replace(clientKey, "-", "", -1)

	client := &entity.Client{
		Name:        clientReq.Name,
		Description: clientReq.Description,
		AppId:       clientReq.AppId,
		ClientKey:   clientKey,
	}

	return s.clientRepo.Create(client)
}

func (s clientService) Update(clientReq serializers.UpdateClientReq, clientId uint) (*entity.Client, *errors.RestErr) {
	client := &entity.Client{
		Name:        clientReq.Name,
		Description: clientReq.Description,
	}

	return s.clientRepo.Update(client, clientId)
}

func (s clientService) Delete(id uint) *errors.RestErr {
	return s.clientRepo.Delete(id)
}

func (s clientService) GetByAppId(appId uint) ([]*entity.Client, *errors.RestErr) {
	return s.clientRepo.GetByAppId(appId)
}

func (s clientService) GetByClientKey(clientKey string) (*entity.Client, *errors.RestErr) {
	return s.clientRepo.GetByClientKey(clientKey)
}
