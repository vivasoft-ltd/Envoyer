package repository

import (
	"envoyer/errors"
	"envoyer/model/entity"
)

type LogRepositoryInterface interface {
	Get(id uint) (*entity.ErrorLog, *errors.RestErr)
	Create(errorLog *entity.ErrorLog) (*entity.ErrorLog, *errors.RestErr)
	Update(errorLog *entity.ErrorLog, errorLogId uint) (*entity.ErrorLog, *errors.RestErr)
	Delete(id uint) *errors.RestErr
	GetByAppId(appId uint) ([]*entity.ErrorLog, *errors.RestErr)
}

type logRepository struct {
	*BaseRepository
}

func NewLogRepository(baseRepo *BaseRepository) LogRepositoryInterface {
	return &logRepository{BaseRepository: baseRepo}
}

func (r logRepository) Get(id uint) (*entity.ErrorLog, *errors.RestErr) {
	errorLog := &entity.ErrorLog{}
	if err := r.Db.First(&errorLog, id).Error; err != nil {
		return nil, errors.NewInternalServerError("error log not found", err)
	}
	return errorLog, nil
}

func (r logRepository) Create(errorLog *entity.ErrorLog) (*entity.ErrorLog, *errors.RestErr) {
	if err := r.Db.Create(errorLog).Error; err != nil {
		return nil, errors.NewInternalServerError("can not create error log", err)
	}
	return errorLog, nil
}

func (r logRepository) Update(errorLog *entity.ErrorLog, errorLogId uint) (*entity.ErrorLog, *errors.RestErr) {
	err := r.Db.Model(&entity.ErrorLog{}).Where("id = ?", errorLogId).Updates(errorLog).Error
	if err != nil {
		return nil, errors.NewInternalServerError("can not update errorLog", err)
	}
	return r.Get(errorLogId)
}

func (r logRepository) Delete(id uint) *errors.RestErr {
	err := r.Db.Where("id = ?", id).Delete(&entity.ErrorLog{}).Error
	if err != nil {
		return errors.NewInternalServerError("can not delete ErrorLog", err)
	}
	return nil
}

func (r logRepository) GetByAppId(appId uint) ([]*entity.ErrorLog, *errors.RestErr) {
	var ErrorLogs []*entity.ErrorLog
	err := r.Db.Where("app_id = ?", appId).Order("date desc").Limit(100).Find(&ErrorLogs).Error
	if err != nil {
		return nil, errors.NewInternalServerError("ErrorLogs not found", err)
	}
	return ErrorLogs, nil
}
