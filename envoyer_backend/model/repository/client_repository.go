package repository

import (
	"envoyer/errors"
	"envoyer/model/entity"
)

type ClientRepositoryInterface interface {
	Get(id uint) (*entity.Client, *errors.RestErr)
	Create(client *entity.Client) (*entity.Client, *errors.RestErr)
	Update(client *entity.Client, clientId uint) (*entity.Client, *errors.RestErr)
	Delete(id uint) *errors.RestErr
	GetByAppId(appId uint) ([]*entity.Client, *errors.RestErr)
	GetByClientKey(clientKey string) (*entity.Client, *errors.RestErr)
}

type clientRepository struct {
	*BaseRepository
}

func NewClientRepository(baseRepo *BaseRepository) ClientRepositoryInterface {
	return &clientRepository{BaseRepository: baseRepo}
}

func (r clientRepository) Get(id uint) (*entity.Client, *errors.RestErr) {
	client := &entity.Client{}
	if err := r.Db.First(&client, id).Error; err != nil {
		return nil, errors.NewInternalServerError("client not found", err)
	}
	return client, nil
}

func (r clientRepository) Create(client *entity.Client) (*entity.Client, *errors.RestErr) {
	if err := r.Db.Create(client).Error; err != nil {
		return nil, errors.NewInternalServerError("can not create client", err)
	}
	return client, nil
}

func (r clientRepository) Update(client *entity.Client, clientId uint) (*entity.Client, *errors.RestErr) {
	err := r.Db.Model(&entity.Client{}).Where("id = ?", clientId).Updates(client).Error
	if err != nil {
		return nil, errors.NewInternalServerError("can not update client", err)
	}
	return r.Get(clientId)
}

func (r clientRepository) Delete(id uint) *errors.RestErr {
	err := r.Db.Where("id = ?", id).Delete(&entity.Client{}).Error
	if err != nil {
		return errors.NewInternalServerError("can not delete client", err)
	}
	return nil
}

func (r clientRepository) GetByAppId(appId uint) ([]*entity.Client, *errors.RestErr) {
	var clients []*entity.Client
	err := r.Db.Where("app_id = ?", appId).Find(&clients).Error
	if err != nil {
		return nil, errors.NewInternalServerError("clients not found", err)
	}
	return clients, nil
}

func (r clientRepository) GetByClientKey(clientKey string) (*entity.Client, *errors.RestErr) {
	var client *entity.Client
	err := r.Db.Where("client_key = ?", clientKey).Find(&client).Error
	if err != nil {
		return nil, errors.NewInternalServerError("client not found", err)
	}
	return client, nil
}
