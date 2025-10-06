package service

import (
	"errors"

	"go.uber.org/zap"
)

type Service struct {
	logger *zap.Logger
}

func NewService(logger *zap.Logger) *Service {
	return &Service{logger: logger}
}

type UserCheckResponse struct {
	Exists   bool   `json:"exists"`
	Role     string `json:"role,omitempty"`
	Username string `json:"username,omitempty"`
}

type TaskResponse struct {
	Exist    int      `json:"exist"`
	ID       int      `json:"id"`
	Question string   `json:"question"`
	Tags     []string `json:"tags"`
}

var ErrServiceUnavailable = errors.New("theory-service unavailable")
