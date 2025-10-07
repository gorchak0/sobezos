package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

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

	if len(tags) == 0 {
		// Нет задач по выбранным тегам вообще
		msg := map[string]string{"message": "❌Доступных задач нет, так как тэги не выбраны\\. Используйте \\/tagset"}
		h.logger.Info("Нет задач по выбранным тегам", zap.Strings("tags", tags))
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(msg)
		return
	}

	// Получаем exclude из query (?exclude=...)
	excludeParam := r.URL.Query().Get("exclude")
	excludeSet := make(map[int]struct{})
	if excludeParam != "" {
		excludeStrs := splitAndTrim(excludeParam, ",")
		for _, s := range excludeStrs {
			if s == "" {
				continue
			}
			if id, err := strconv.Atoi(s); err == nil {
				excludeSet[id] = struct{}{}
			}
		}
	}

	// Получаем id задач по тегам (или все задачи, если теги не заданы)
	ids, err := h.service.GetTaskIDsByTags(tags)
	if err != nil {
		h.logger.Error("Ошибка получения id задач по тегам", zap.Error(err))
		http.Error(w, "Failed to get task", http.StatusInternalServerError)
		return
	}

	// Фильтруем задачи, исключая прорешённые
	var availableIDs []int
	for _, id := range ids {
		if _, found := excludeSet[id]; !found {
			availableIDs = append(availableIDs, id)
		}
	}

	if len(availableIDs) == 0 {
		// Все задачи по тегам прорешаны
		msg := map[string]string{"message": "Поздравляем, вы прорешали все задачи по доступным тэгам"}
		h.logger.Info("Все задачи по тегам прорешаны", zap.Strings("tags", tags), zap.Any("exclude", excludeSet))
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(msg)
		return
	}

	// Выбираем случайную задачу из доступных
	randIdx := time.Now().UnixNano() % int64(len(availableIDs))
	taskID := availableIDs[randIdx]
	task, err := h.service.GetTaskByID(taskID)
	if err != nil || task == nil {
		h.logger.Error("Ошибка получения задачи по id", zap.Error(err), zap.Int("id", taskID))
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
