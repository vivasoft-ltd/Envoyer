package service

import (
	"envoyer/config"
	"envoyer/config/consts"
	"envoyer/config/service_name"
	"envoyer/errors"
	"envoyer/model/entity"
	"envoyer/model/repository"
	"envoyer/model/serializers"
	"envoyer/utils"
	"gorm.io/gorm"
	"time"
)

type AuthServiceInterface interface {
	LoginUser(tokenRequest serializers.LoginReq) (*serializers.JWTTokenResponse, *errors.RestErr)
	GetAccessToken(user *entity.User) (*serializers.JWTTokenResponse, *errors.RestErr)
	FindUserById(id uint) (*entity.User, *errors.RestErr)
}

type authService struct {
	*BaseService
	userRepo repository.UserRepositoryInterface
}

func (s authService) FindUserById(id uint) (*entity.User, *errors.RestErr) {
	return s.userRepo.Get(id)
}

func NewAuthService(baseService *BaseService) AuthServiceInterface {
	return &authService{
		BaseService: baseService,
		userRepo:    baseService.container.Get(service_name.UserRepository).(repository.UserRepositoryInterface),
	}
}

func (s authService) LoginUser(tokenRequest serializers.LoginReq) (*serializers.JWTTokenResponse, *errors.RestErr) {
	if tokenRequest.Username == config.Config.SuperUserName && tokenRequest.Password == config.Config.SuperPassword {
		return s.GetAccessToken(&entity.User{
			Model: gorm.Model{
				ID: 0,
			},
			Role: consts.SuperAdminRole,
		})
	}
	user, err := s.userRepo.GetByUsername(tokenRequest.Username)
	if err != nil {
		return nil, err
	}
	if user.Password != tokenRequest.Password {
		return nil, errors.NewUnauthorizedError("invalid credentials", errors.NewError("invalid credentials"))
	}

	return s.GetAccessToken(user)
}

func (s authService) GetAccessToken(user *entity.User) (*serializers.JWTTokenResponse, *errors.RestErr) {
	accessToken, err := utils.CreateToken(config.Config.AccessTokenExpiresIn, user.ID, config.Config.AccessTokenPrivateKey)
	if err != nil {
		return nil, errors.NewInternalServerError("can not create access token", err)
	}

	refreshToken, err := utils.CreateToken(config.Config.RefreshTokenExpiresIn, user.ID, config.Config.RefreshTokenPrivateKey)
	if err != nil {
		return nil, errors.NewInternalServerError("can not create refresh token", err)
	}

	expirationTime := time.Now().Add(config.Config.RefreshTokenExpiresIn)

	return &serializers.JWTTokenResponse{
		Token:     accessToken,
		Refresh:   refreshToken,
		ExpiredAt: expirationTime.Unix(),
		Id:        user.ID,
		Role:      user.Role,
		AppId:     user.AppId,
	}, nil
}
