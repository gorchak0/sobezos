package handler

import (
	"encoding/json"
	"net/http"

	"sobezos/services/core-service/internal/models"
)

// TagClear godoc
// @Summary Очистить теги пользователя
// @Description Очищает все теги пользователя
// @Tags tag
// @Accept json
// @Produce json
// @Param data body models.TagClearRequest true "Telegram ID пользователя"
// @Success 200 {object} models.CommonSuccessResponse
// @Failure 400 {object} models.CommonErrorResponse "invalid body"
// @Failure 500 {object} models.CommonErrorResponse "internal error"
// @Router /tagclear [post]
func (h *Handler) TagClear(w http.ResponseWriter, r *http.Request) {
	var req models.TagClearRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.CommonErrorResponse{Error: "invalid body"})
		return
	}
	result, err := h.service.TagClear(req.TelegramID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.CommonErrorResponse{Error: err.Error()})
		return
	}
	json.NewEncoder(w).Encode(models.CommonSuccessResponse{Result: result})
}
