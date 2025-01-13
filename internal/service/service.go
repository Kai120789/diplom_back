package service

import (
	"materials/internal/config"

	"go.uber.org/zap"
)

type Service struct {
	UserService *UserService
}

type Storage struct {
	UserStorage UserStorage
}

func NewService(db Storage, cfg *config.Config, log *zap.Logger) *Service {
	return &Service{
		UserService: NewUserService(db.UserStorage, cfg, log),
	}
}
