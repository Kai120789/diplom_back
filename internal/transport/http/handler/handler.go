package handler

import "materials/internal/config"

type Handler struct {
	UserHandler *UserHandler
}

type Service struct {
	UserService UserService
}

func NewHandler(s Service, c *config.Config) *Handler {
	return &Handler{
		UserHandler: NewUserHandler(s.UserService, c),
	}
}
