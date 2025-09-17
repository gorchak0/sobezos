package handlers

import (
	"encoding/json"
	"net/http"
	"sobezos/services/user-service/internal/repository"
	"sobezos/services/user-service/pkg/models"

	"go.uber.org/zap"
)

// AddUser обрабатывает HTTP-запрос на добавление нового пользователя.
// Ожидает JSON с telegram_id, username и role в теле запроса.
func (h *UserServiceHandler) AddUser(w http.ResponseWriter, r *http.Request) {
	// Структура для декодирования входящего JSON
	var req struct {
		TelegramID int64  `json:"telegram_id"`
		Username   string `json:"username"`
		Role       string `json:"role"`
	}
	// Декодируем тело запроса в структуру req
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.Logger.Error("Invalid JSON in AddUser", zap.Error(err))
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	h.Logger.Info(
		"AddUser request received",
		zap.Int64("telegram_id", req.TelegramID),
		zap.String("username", req.Username),
		zap.String("role", req.Role),
	)
	// Формируем структуру пользователя для передачи в сервис
	user := models.User{
		TelegramID: req.TelegramID,
		Username:   req.Username,
		Role:       req.Role,
	}
	// Пытаемся добавить пользователя через сервис (без проверки администратора)
	err := h.Service.AddUser(user)
	if err == repository.ErrForbidden {
		h.Logger.Warn("Forbidden to add user", zap.Any("user", user))
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	if err == repository.ErrUserExists {
		h.Logger.Warn("User already exists", zap.Any("user", user))
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}
	if err != nil {
		h.Logger.Error("Failed to add user", zap.Any("user", user), zap.Error(err))
		http.Error(w, "Failed to add user", http.StatusInternalServerError)
		return
	}
	h.Logger.Info("User successfully added", zap.Any("user", user))
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{"ok": true})
}
