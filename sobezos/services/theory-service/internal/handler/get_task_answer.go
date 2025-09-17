package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"go.uber.org/zap"
)

// GET /answer?task_id=...
func (h *TaskHandler) GetTaskAnswer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.logger.Warn("Method not allowed", zap.String("method", r.Method))
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	taskIDStr := r.URL.Query().Get("task_id")
	if taskIDStr == "" {
		http.Error(w, "Missing task_id", http.StatusBadRequest)
		return
	}
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		http.Error(w, "Invalid task_id", http.StatusBadRequest)
		return
	}
	task, err := h.service.GetTaskByID(taskID)
	if err != nil || task == nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}
	resp := map[string]interface{}{
		"answer": task.Answer,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
