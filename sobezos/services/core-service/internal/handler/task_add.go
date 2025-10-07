package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"sobezos/services/core-service/internal/models"
)

// TaskAdd godoc
// @Summary Добавить задачу
// @Description Добавляет новую задачу (требуется telegram_id в заголовке)
// @Tags task
// @Accept json
// @Produce json
// @Param telegram_id query int true "Telegram ID пользователя"
// @Param data body models.TaskAddRequest true "JSON задачи"
// @Success 200 {object} models.CommonSuccessResponse
// @Failure 400 {object} models.CommonErrorResponse "invalid telegram_id"
// @Failure 500 {object} models.CommonErrorResponse "internal error"
// @Router /taskadd [post]
func (h *Handler) TaskAdd(w http.ResponseWriter, r *http.Request) {
	telegramIDStr := r.URL.Query().Get("telegram_id")
	telegramID, err := strconv.Atoi(telegramIDStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.CommonErrorResponse{Error: "invalid telegram_id"})
		return
	}

	var req models.TaskAddRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.CommonErrorResponse{Error: "invalid body"})
		return
	}

	result, err := h.service.TaskAdd(telegramID, req.Question, req.Answer, req.Tags)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.CommonErrorResponse{Error: err.Error()})
		return
	}
	json.NewEncoder(w).Encode(models.CommonSuccessResponse{Result: result})
}
