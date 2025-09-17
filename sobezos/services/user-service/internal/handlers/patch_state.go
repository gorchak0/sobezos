package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

// PATCH /states/{user_id}
func (h *UserServiceHandler) PatchState(w http.ResponseWriter, r *http.Request) {
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
	var patchMap map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&patchMap); err != nil {
		h.Logger.Warn("Invalid JSON", zap.Error(err))
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if len(patchMap) == 0 {
		h.Logger.Warn("No fields to update", zap.Int64("user_id", userID))
		http.Error(w, "No fields to update", http.StatusBadRequest)
		return
	}
	if err := h.Service.PatchState(userID, patchMap); err != nil {
		h.Logger.Error("Failed to update user state", zap.Int64("user_id", userID), zap.Error(err))
		http.Error(w, "Failed to update state", http.StatusInternalServerError)
		return
	}
	h.Logger.Info("User state updated", zap.Int64("user_id", userID), zap.Any("patch", patchMap))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"ok": true})
}
