package service

import (
	"sobezos/services/user-service/internal/repository"
	"sobezos/services/user-service/pkg/models"
	"time"
)

type UserService struct {
	UserRepo  *repository.UserRepository
	StateRepo *repository.UserStateRepository
}

func NewUserService(userRepo *repository.UserRepository, stateRepo *repository.UserStateRepository) *UserService {
	return &UserService{
		UserRepo:  userRepo,
		StateRepo: stateRepo,
	}
}

func (s *UserService) CheckUser(telegramID int64) (exists bool, role string, err error) {
	user, err := s.UserRepo.GetByTelegramID(telegramID)
	if err != nil {
		return false, "", err
	}
	return true, user.Role, nil
}

func (s *UserService) AddUser(user models.User) error {
	if s.UserRepo.Exists(user.TelegramID) {
		return repository.ErrUserExists
	}
	err := s.UserRepo.Add(user)
	if err != nil {
		return err
	}
	// Создать запись в user_state с полями по умолчанию через Add
	now := time.Now()
	defaultState := models.UserState{
		UserID:           user.TelegramID,
		LastTheoryTaskID: nil,
		TheoryTags:       []string{},
		CompletedTheoryTasks: []string{},
		LastAction:       nil,
		UpdatedAt:        &now,
	}
	if err := s.StateRepo.AddState(defaultState); err != nil {
		return err
	}
	return nil
}

func (s *UserService) GetState(userID int64) (interface{}, error) {
	return s.StateRepo.Get(userID)
}

func (s *UserService) PatchState(userID int64, patch models.UserState) error {
	return s.StateRepo.Patch(userID, patch)
}
