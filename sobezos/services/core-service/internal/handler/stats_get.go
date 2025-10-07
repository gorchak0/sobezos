package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"sobezos/services/core-service/internal/models"
)

// StatsGet godoc
// @Summary Получить статистику пользователя
// @Description Получить статистику пользователя
// @Tags user
// @Produce json
// @Param telegram_id query int true "Telegram ID пользователя"
// @Success 200 {object} models.CommonSuccessResponse
// @Failure 400 {object} models.CommonErrorResponse "invalid telegram_id"
// @Failure 500 {object} models.CommonErrorResponse "internal error"
// @Router /statsget [get]
func (h *Handler) StatsGet(w http.ResponseWriter, r *http.Request) {
	telegramIDStr := r.URL.Query().Get("telegram_id")
	telegramID, err := strconv.Atoi(telegramIDStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.CommonErrorResponse{Error: "invalid telegram_id"})
		return
	}
	result, err := h.service.StatsGet(telegramID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.CommonErrorResponse{Error: err.Error()})
		return
	}
	json.NewEncoder(w).Encode(models.CommonSuccessResponse{Result: result})
}
