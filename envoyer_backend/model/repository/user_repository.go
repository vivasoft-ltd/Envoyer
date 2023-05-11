package repository

import (
	"envoyer/errors"
	"envoyer/model/entity"
	"gorm.io/gorm"
	"time"
)

type UserRepositoryInterface interface {
	Get(id uint) (*entity.User, *errors.RestErr)
	Create(user *entity.User) (*entity.User, *errors.RestErr)
	Update(user *entity.User, userId uint) (*entity.User, *errors.RestErr)
	Delete(id uint) *errors.RestErr
	GetByAppId(appId uint) ([]*entity.User, *errors.RestErr)
	GetByUsername(username string) (*entity.User, *errors.RestErr)
}

type userRepository struct {
	*BaseRepository
}

func NewUserRepository(baseRepo *BaseRepository) UserRepositoryInterface {
	return &userRepository{BaseRepository: baseRepo}
}

func (r userRepository) GetByUsername(username string) (*entity.User, *errors.RestErr) {
	var user *entity.User
	err := r.Db.Where("user_name = ?", username).First(&user).Error
	if err != nil {
		return nil, errors.NewInternalServerError("user not found", err)
	}
	return user, nil
}

func (r userRepository) Get(id uint) (*entity.User, *errors.RestErr) {
	user := &entity.User{}
	if err := r.Db.First(&user, id).Error; err != nil {
		return nil, errors.NewInternalServerError("user not found", err)
	}
	return user, nil
}

func (r userRepository) Create(user *entity.User) (*entity.User, *errors.RestErr) {
	if err := r.Db.Create(user).Error; err != nil {
		return nil, errors.NewInternalServerError("can not create user", err)
	}
	return user, nil
}

func (r userRepository) Update(user *entity.User, userId uint) (*entity.User, *errors.RestErr) {
	err := r.Db.Model(&entity.User{}).Where("id = ?", userId).Updates(user).Error
	if err != nil {
		return nil, errors.NewInternalServerError("can not update user", err)
	}
	return r.Get(userId)
}

func (r userRepository) Delete(id uint) *errors.RestErr {
	currentTime := time.Now().String() + "@"
	err := r.Db.Model(entity.User{}).Where("id = ?", id).Update("user_name", gorm.Expr("CONCAT(? , user_name)", currentTime)).Delete(&entity.User{}).Error
	if err != nil {
		return errors.NewInternalServerError("can not delete user", err)
	}
	return nil
}

func (r userRepository) GetByAppId(appId uint) ([]*entity.User, *errors.RestErr) {
	var users []*entity.User
	err := r.Db.Where("app_id = ?", appId).Find(&users).Error
	if err != nil {
		return nil, errors.NewInternalServerError("users not found", err)
	}
	return users, nil
}
