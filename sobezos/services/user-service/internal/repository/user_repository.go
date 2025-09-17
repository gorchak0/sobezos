package repository

import (
	"database/sql"
	"errors"
	"sobezos/services/user-service/pkg/models"
)

var (
	ErrForbidden  = errors.New("forbidden")
	ErrUserExists = errors.New("user already exists")
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) GetByTelegramID(telegramID int64) (*models.User, error) {
	row := r.DB.QueryRow("SELECT id, telegram_id, username, role, created_at FROM users WHERE telegram_id=$1", telegramID)
	var user models.User
	if err := row.Scan(&user.ID, &user.TelegramID, &user.Username, &user.Role, &user.CreatedAt); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Exists(telegramID int64) bool {
	row := r.DB.QueryRow("SELECT 1 FROM users WHERE telegram_id=$1", telegramID)
	var exists int
	return row.Scan(&exists) == nil
}

func (r *UserRepository) Add(user models.User) error {
	_, err := r.DB.Exec("INSERT INTO users (telegram_id, username, role) VALUES ($1, $2, $3)", user.TelegramID, user.Username, user.Role)
	return err
}

func (r *UserRepository) IsAdmin(telegramID int64) bool {
	row := r.DB.QueryRow("SELECT role FROM users WHERE telegram_id=$1", telegramID)
	var role string
	if err := row.Scan(&role); err != nil {
		return false
	}
	return role == "admin"
}
