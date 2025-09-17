package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func (h *UserServiceHandler) CheckUser(w http.ResponseWriter, r *http.Request) {
	telegramIDStr := r.URL.Query().Get("telegram_id")
	telegramID, err := strconv.ParseInt(telegramIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid telegram_id", http.StatusBadRequest)
		return
	}
	user, err := h.Service.UserRepo.GetByTelegramID(telegramID)
	if err != nil || user == nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"exists": false})
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"exists": true, "role": user.Role, "username": user.Username})
}
