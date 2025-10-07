package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"sobezos/services/core-service/internal/models"
)

// UserAdd godoc
// @Summary Добавить пользователя
// @Description Добавляет нового пользователя (требуется telegram_id администратора в заголовке)
// @Tags user
// @Accept json
// @Produce json
// @Param telegram_id query int true "Telegram ID пользователя"
// @Param data body models.UserAddRequest true "Аргументы для добавления пользователя"
// @Success 200 {object} models.CommonSuccessResponse
// @Failure 400 {object} models.CommonErrorResponse "invalid admin telegram_id"
// @Failure 500 {object} models.CommonErrorResponse "internal error"
// @Router /useradd [post]
func (h *Handler) UserAdd(w http.ResponseWriter, r *http.Request) {
	telegramIDStr := r.URL.Query().Get("telegram_id")
	telegramID, err := strconv.Atoi(telegramIDStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.CommonErrorResponse{Error: "invalid admin telegram_id"})
		return
	}
	var req models.UserAddRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.CommonErrorResponse{Error: "invalid body"})
		return
	}
	result, err := h.service.UserAdd(telegramID, req.Args) //
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.CommonErrorResponse{Error: err.Error()})
		return
	}
	json.NewEncoder(w).Encode(models.CommonSuccessResponse{Result: result})
}
