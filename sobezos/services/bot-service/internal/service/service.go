package service

import (
	"go.uber.org/zap"
)

type ServiceInterface interface {
	Help() (string, error)
	TaskGet(telegramID int) (string, error)
	TaskGetID(telegramID int, args string) (string, error)
	AnswerGet(telegramID int) (string, error)
	TagSet(telegramID int, args string) (string, error)
	TagClear(telegramID int) (string, error)
	TagGet() (string, error)
	TaskAdd(telegramID int, args string) (string, error)
	TaskEdit(telegramID int, args string) (string, error)
	UserAdd(telegramID int, args string) (string, error)
	UserCheck(telegramID int) (userInfo struct{ Username, Role string }, exists bool)
	StatsGet(telegramID int) (string, error)
}

type Service struct {
	Logger         *zap.Logger
	CoreServiceUrl string
}

func NewService(logger *zap.Logger) ServiceInterface {
	return &Service{
		Logger:         logger,
		CoreServiceUrl: "http://core-service:8083",
	}
}

type commonSuccessResponse struct {
	Result string `json:"result"`
}
type userCheckResponse struct {
	Role     string `json:"role"`
	Username string `json:"username"`
}
