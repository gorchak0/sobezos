package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"sobezos/services/core-service/internal/models"
)

// TaskGetID godoc
// @Summary Получить задачу по ID
// @Description Получает задачу по её ID для пользователя
// @Tags task
// @Produce json
// @Param telegram_id query int true "Telegram ID пользователя"
// @Param args query int true "ID задачи"
// @Success 200 {object} models.CommonSuccessResponse
// @Failure 400 {object} models.CommonErrorResponse "invalid telegram_id"
// @Failure 500 {object} models.CommonErrorResponse "internal error"
// @Router /taskgetid [get]
func (h *Handler) TaskGetID(w http.ResponseWriter, r *http.Request) {
	telegramIDStr := r.URL.Query().Get("telegram_id")
	args := r.URL.Query().Get("args")
	telegramID, err := strconv.Atoi(telegramIDStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.CommonErrorResponse{Error: "invalid telegram_id"})
		return
	}
	result, err := h.service.TaskGetID(telegramID, args)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.CommonErrorResponse{Error: err.Error()})
		return
	}
	json.NewEncoder(w).Encode(models.CommonSuccessResponse{Result: result})
}
