package repository

import (
	"envoyer/logger"
	"gorm.io/gorm"
)

type BaseRepository struct {
	logger logger.Logger
	Db     *gorm.DB
}

func NewBaseRepository(db *gorm.DB, logger logger.Logger) *BaseRepository {
	return &BaseRepository{Db: db, logger: logger}
}
