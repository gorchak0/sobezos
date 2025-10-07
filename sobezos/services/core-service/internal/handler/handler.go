package handler

import (
	"sobezos/services/core-service/internal/service"
)

type Handler struct {
	service *service.Service
}

func NewHandler(s *service.Service) *Handler {
	return &Handler{service: s}
}
