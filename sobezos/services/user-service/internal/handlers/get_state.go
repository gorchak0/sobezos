package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"go.uber.org/zap"
)

// "http://user-service:8082/userstateget?user_id="
func (h *UserServiceHandler) GetState(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		h.Logger.Warn("Missing user_id query param", zap.String("url", r.URL.String()))
		http.Error(w, "Missing user_id query param", http.StatusBadRequest)
		return
	}
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		h.Logger.Warn("Invalid user_id", zap.String("user_id", userIDStr), zap.Error(err))
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}
	state, err := h.Service.GetState(userID)
	if err != nil {
		h.Logger.Info("User state not found", zap.Int64("user_id", userID))
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	h.Logger.Info("User state returned", zap.Int64("user_id", userID), zap.Any("state", state))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(state)
}
