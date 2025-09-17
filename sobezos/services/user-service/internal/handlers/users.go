package handlers

import (
	"sobezos/services/user-service/internal/repository"
	"sobezos/services/user-service/internal/service"

	"go.uber.org/zap"
)

type UserServiceHandler struct {
	Service *service.UserService
	Logger  *zap.Logger
}

func NewUserServiceHandler(userRepo *repository.UserRepository, stateRepo *repository.UserStateRepository, logger *zap.Logger) *UserServiceHandler {
	return &UserServiceHandler{
		Service: service.NewUserService(userRepo, stateRepo),
		Logger:  logger,
	}
}
