package service

import (
	"sobezos/services/user-service/internal/repository"
	"sobezos/services/user-service/pkg/models"
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
	return s.UserRepo.Add(user)
}

func (s *UserService) GetState(userID int64) (interface{}, error) {
	return s.StateRepo.Get(userID)
}

func (s *UserService) PatchState(userID int64, patch models.UserState) error {
	return s.StateRepo.Patch(userID, patch)
}
