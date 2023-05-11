package service

import (
	"envoyer/logger"
	"github.com/sarulabs/di/v2"
)

type BaseService struct {
	logger    logger.Logger
	container di.Container
}

func NewBaseService(container di.Container, logger logger.Logger) *BaseService {
	return &BaseService{container: container, logger: logger}
}
