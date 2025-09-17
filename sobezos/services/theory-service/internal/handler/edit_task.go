package handler

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

// EditTask handles PUT /edittask for updating a task
func (h *TaskHandler) EditTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		h.logger.Warn("Method not allowed", zap.String("method", r.Method))
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("Invalid request body", zap.Error(err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	idVal, ok := req["id"]
	if !ok {
		h.logger.Warn("Missing id in request")
		http.Error(w, "Missing id field", http.StatusBadRequest)
		return
	}
	id, ok := idVal.(float64)
	if !ok {
		h.logger.Warn("Invalid id type")
		http.Error(w, "Invalid id type", http.StatusBadRequest)
		return
	}
	// Получаем остальные поля, если они есть
	var tags []string
	if v, ok := req["tags"]; ok {
		arr, ok := v.([]interface{})
		if ok {
			for _, t := range arr {
				if str, ok := t.(string); ok {
					tags = append(tags, str)
				}
			}
		}
	}
	var question, answer string
	if v, ok := req["question"]; ok {
		question, _ = v.(string)
	}
	if v, ok := req["answer"]; ok {
		answer, _ = v.(string)
	}
	if err := h.service.UpdateTask(int(id), tags, question, answer); err != nil {
		h.logger.Error("Failed to update task", zap.Error(err))
		http.Error(w, "Failed to update task", http.StatusInternalServerError)
		return
	}
	h.logger.Info("Task updated", zap.Int("id", int(id)))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}
