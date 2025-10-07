package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"sobezos/services/core-service/internal/models"
)

// AnswerGet godoc
// @Summary Получить ответ на последний вопрос пользователя
// @Description Возвращает ответ на последний теоретический вопрос пользователя
// @Tags answer
// @Produce json
// @Param telegram_id query int true "Telegram ID пользователя"
// @Success 200 {object} models.CommonSuccessResponse
// @Failure 400 {object} models.CommonErrorResponse "invalid telegram_id"
// @Failure 500 {object} models.CommonErrorResponse "internal error"
// @Router /answerget [get]
func (h *Handler) AnswerGet(w http.ResponseWriter, r *http.Request) {
	telegramIDStr := r.URL.Query().Get("telegram_id")
	telegramID, err := strconv.Atoi(telegramIDStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.CommonErrorResponse{Error: "invalid telegram_id"})
		return
	}
	result, err := h.service.AnswerGet(telegramID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.CommonErrorResponse{Error: err.Error()})
		return
	}
	json.NewEncoder(w).Encode(models.CommonSuccessResponse{Result: result})
}
