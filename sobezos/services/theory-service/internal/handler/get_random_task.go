package handler

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

func (h *TaskHandler) GetRandomTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.logger.Warn("Method not allowed", zap.String("method", r.Method))
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	h.logger.Info("Получен GET-запрос на задачу", zap.String("method", r.Method))

	//тут обращение к user-states для получения тэгов пользователя
	task, err := h.service.GetRandomTask()
	if err != nil {
		h.logger.Error("Ошибка получения задачи из сервиса", zap.Error(err))
		http.Error(w, "Failed to get task", http.StatusInternalServerError)
		return
	}
	response := struct {
		ID       int      `json:"id"`
		Tags     []string `json:"tags"`
		Question string   `json:"question"`
	}{
		ID:       task.ID,
		Tags:     task.Tags,
		Question: task.Question,
	}
	h.logger.Info("Задача получена", zap.Int("id", task.ID), zap.String("question", task.Question), zap.Any("tags", task.Tags))
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Error("Ошибка кодирования ответа", zap.Error(err))
	} else {
		h.logger.Info("Ответ успешно отправлен", zap.Any("response", response))
	}
}
