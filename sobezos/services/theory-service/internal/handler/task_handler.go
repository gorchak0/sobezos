package handler

import (
	"sobezos/services/theory-service/internal/service"

	"go.uber.org/zap"
)

type TaskHandler struct {
	service *service.TaskService
	logger  *zap.Logger
}

func NewTaskHandler(service *service.TaskService, logger *zap.Logger) *TaskHandler {
	return &TaskHandler{service: service, logger: logger}
}
