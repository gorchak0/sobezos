package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

func (h *TaskHandler) TaskGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.logger.Warn("Method not allowed", zap.String("method", r.Method))
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	h.logger.Info("Получен GET-запрос на задачу", zap.String("method", r.Method))

	// Получаем тэги из query (?tags=...)
	tagsParam := r.URL.Query().Get("tags")
	var tags []string
	if tagsParam != "" {
		for _, t := range splitAndTrim(tagsParam, ",") {
			if t != "" {
				tags = append(tags, t)
			}
		}
	}

	fmt.Printf("\n\ntheory-service  TaskGet %s \n", tags)
	//
	// Получаем задачу с учётом тегов
	task, err := h.service.GetRandomTask(tags)
	if err != nil || task == nil {
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

// splitAndTrim разбивает строку по разделителю и обрезает пробелы
func splitAndTrim(s, sep string) []string {
	var res []string
	for _, part := range strings.Split(s, sep) {
		trimmed := strings.TrimSpace(part)
		res = append(res, trimmed)
	}
	return res
}
