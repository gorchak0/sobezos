package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"go.uber.org/zap"
)

// TaskGetID обрабатывает GET /taskgetid?task_id=... и возвращает задачу по id
func (h *TaskHandler) TaskGetID(w http.ResponseWriter, r *http.Request) {
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
		h.logger.Warn("Invalid task_id", zap.String("task_id", taskIDStr))
		http.Error(w, "Некорректный id задачи", http.StatusBadRequest)
		return
	}
	task, err := h.service.GetTaskByID(taskID)
	if err != nil || task == nil {
		h.logger.Error("Task not found or error", zap.Int("task_id", taskID), zap.Error(err))
		http.Error(w, "Задача не найдена", http.StatusNotFound)
		return
	}
	resp := map[string]interface{}{
		"id":       task.ID,
		"question": task.Question,
		"tags":     task.Tags,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
