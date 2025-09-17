package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

// GET /states/{user_id}
func (h *UserServiceHandler) GetState(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		h.Logger.Warn("Invalid path", zap.String("path", r.URL.Path))
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}
	userID, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		h.Logger.Warn("Invalid user_id", zap.String("user_id", parts[2]), zap.Error(err))
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
