package service

import (
	"envoyer/config/service_name"
	"envoyer/errors"
	"envoyer/model/entity"
	"envoyer/model/repository"
	"envoyer/model/serializers"
)

type UserServiceInterface interface {
	Get(id uint) (*entity.User, *errors.RestErr)
	Create(userReq serializers.CreateUserReq) (*entity.User, *errors.RestErr)
	Update(userReq serializers.UpdateUserReq, userId uint) (*entity.User, *errors.RestErr)
	Delete(id uint) *errors.RestErr
	GetByAppId(appId uint) ([]*entity.User, *errors.RestErr)
}

type userService struct {
	*BaseService
	userRepo repository.UserRepositoryInterface
}

func NewUserService(baseService *BaseService) UserServiceInterface {
	return &userService{
		BaseService: baseService,
		userRepo:    baseService.container.Get(service_name.UserRepository).(repository.UserRepositoryInterface),
	}
}

func (s userService) Get(id uint) (*entity.User, *errors.RestErr) {
	return s.userRepo.Get(id)
}

func (s userService) Create(userReq serializers.CreateUserReq) (*entity.User, *errors.RestErr) {
	user := &entity.User{
		UserName: userReq.UserName,
		Password: userReq.Password,
		AppId:    userReq.AppId,
		Role:     userReq.Role,
	}
	return s.userRepo.Create(user)
}

func (s userService) Update(userReq serializers.UpdateUserReq, userId uint) (*entity.User, *errors.RestErr) {
	user := &entity.User{
		UserName: userReq.UserName,
		Password: userReq.Password,
		Role:     userReq.Role,
	}
	return s.userRepo.Update(user, userId)
}

func (s userService) Delete(id uint) *errors.RestErr {
	return s.userRepo.Delete(id)
}

func (s userService) GetByAppId(appId uint) ([]*entity.User, *errors.RestErr) {
	return s.userRepo.GetByAppId(appId)
}
