package handler

import (
	"encoding/json"
	"net/http"

	"sobezos/services/core-service/internal/models"
)

// TagGet godoc
// @Summary Получить список тегов
// @Description Получает список всех доступных тегов
// @Tags tag
// @Produce json
// @Success 200 {object} models.CommonSuccessResponse
// @Failure 500 {object} models.CommonErrorResponse "internal error"
// @Router /tagget [get]
func (h *Handler) TagGet(w http.ResponseWriter, r *http.Request) {
	result, err := h.service.TagGet()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.CommonErrorResponse{Error: err.Error()})
		return
	}
	json.NewEncoder(w).Encode(models.CommonSuccessResponse{Result: result})
}
