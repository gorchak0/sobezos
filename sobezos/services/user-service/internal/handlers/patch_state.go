package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"sobezos/services/user-service/pkg/models"

	"go.uber.org/zap"
)

// "http://user-service:8082/userstateedit?user_id="
func (h *UserServiceHandler) PatchState(w http.ResponseWriter, r *http.Request) {
	// Получаем user_id из query параметра
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		h.Logger.Warn("Missing user_id query param", zap.String("url", r.URL.String()))
		http.Error(w, "Missing user_id query param", http.StatusBadRequest)
		return
	}
	// Преобразуем user_id из строки в int64
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		h.Logger.Warn("Invalid user_id", zap.String("user_id", userIDStr), zap.Error(err))
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}

	// Декодируем JSON-патч из тела запроса в структуру
	var patch models.UserState
	if err := json.NewDecoder(r.Body).Decode(&patch); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	// Логируем JSON через fmt для отладки
	patchJson, _ := json.Marshal(patch)
	fmt.Printf("\n\n\nPatchState raw JSON: %s\n\n\n", string(patchJson))

	h.Logger.Info("PatchState have JSON: ", zap.Any("patch", patch))

	// Пытаемся обновить состояние пользователя через сервис
	if err := h.Service.PatchState(userID, patch); err != nil {
		h.Logger.Error("Failed to update user state", zap.Int64("user_id", userID), zap.Error(err))
		http.Error(w, "Failed to update state", http.StatusInternalServerError)
		return
	}

	//getstate получаем и логгируем для отладки
	state, _ := h.Service.GetState(userID)
	fmt.Printf("\n\n\nPatchState new state: %+v\n\n\n", state)

	// Логируем успешное обновление, возвращаем ok:true
	h.Logger.Info("User state updated", zap.Int64("user_id", userID), zap.Any("patch", patch))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]interface{}{"ok": true})
}
