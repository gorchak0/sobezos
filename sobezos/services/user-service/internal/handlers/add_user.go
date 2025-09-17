package handlers

import (
	"encoding/json"
	"net/http"
	"sobezos/services/user-service/internal/repository"
	"sobezos/services/user-service/pkg/models"
	"strconv"
)

func (h *UserServiceHandler) AddUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TelegramID int64  `json:"telegram_id"`
		Username   string `json:"username"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	adminIDStr := r.Header.Get("X-Admin-Telegram-ID")
	adminID, err := strconv.ParseInt(adminIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	user := models.User{
		TelegramID: req.TelegramID,
		Username:   req.Username,
	}
	err = h.Service.AddUser(adminID, user)
	if err == repository.ErrForbidden {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	if err == repository.ErrUserExists {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}
	if err != nil {
		http.Error(w, "Failed to add user", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{"ok": true})
}
