package handler

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

// GetAllTags handles GET /tags
func (h *TaskHandler) GetAllTags(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.logger.Warn("Method not allowed", zap.String("method", r.Method))
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tags, err := h.service.GetAllTags()
	if err != nil {
		h.logger.Error("Failed to get tags", zap.Error(err))
		http.Error(w, "Failed to get tags", http.StatusInternalServerError)
		return
	}

	response := make([]map[string]interface{}, 0, len(tags))
	for _, tag := range tags {
		response = append(response, map[string]interface{}{
			"id":          tag.ID,
			"name":        tag.Name,
			"description": tag.Description,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Error("Failed to encode response", zap.Error(err))
	}
}
