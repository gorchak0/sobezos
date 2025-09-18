package handler

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

// TaskAdd handles POST /taskadd for adding a new task
func (h *TaskHandler) TaskAdd(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.logger.Warn("Method not allowed", zap.String("method", r.Method))
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Tags     []string `json:"tags"`
		Question string   `json:"question"`
		Answer   string   `json:"answer"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("Invalid request body", zap.Error(err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	id, err := h.service.CreateTask(req.Tags, req.Question, req.Answer)
	if err != nil {
		h.logger.Error("Failed to create task", zap.Error(err))
		http.Error(w, "Failed to create task", http.StatusInternalServerError)
		return
	}
	h.logger.Info("Task created", zap.String("question", req.Question), zap.Int("id", id))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{"status": "ok", "id": id})
}
