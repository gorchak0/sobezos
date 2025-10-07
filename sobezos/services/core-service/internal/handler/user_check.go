package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"sobezos/services/core-service/internal/models"
)

// UserCheck godoc
// @Summary Проверить пользователя
// @Description Проверяет существование пользователя и возвращает его роль и username
// @Tags user
// @Produce json
// @Param telegram_id query int true "Telegram ID пользователя"
// @Success 200 {object} models.UserCheckResponse
// @Failure 400 {object} models.CommonErrorResponse "invalid telegram_id"
// @Failure 404 {object} models.CommonErrorResponse "user not found"
// @Router /usercheck [get]
func (h *Handler) UserCheck(w http.ResponseWriter, r *http.Request) {
	telegramIDStr := r.URL.Query().Get("telegram_id")
	telegramID, err := strconv.Atoi(telegramIDStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.CommonErrorResponse{Error: "invalid telegram_id"})
		return
	}
	res, exists := h.service.UserCheck(telegramID)
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(models.CommonErrorResponse{Error: "user not found"})
		return
	}
	json.NewEncoder(w).Encode(models.UserCheckResponse{
		Role:     res.Role,
		Username: res.Username,
	})
}
