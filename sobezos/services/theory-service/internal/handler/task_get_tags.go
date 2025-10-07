package handler

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

// TaskGetTags обрабатывает GET /taskgettags?tags=tag1,tag2 и возвращает id задач по тегам
func (h *TaskHandler) TaskGetTags(w http.ResponseWriter, r *http.Request) {
	// Логируем вызов хендлера с методом и URL
	h.logger.Info("TaskGetTags handler called", zap.String("method", r.Method), zap.String("url", r.URL.String()))

	// Проверяем, что используется метод GET
	if r.Method != http.MethodGet {
		h.logger.Warn("Method not allowed", zap.String("method", r.Method))
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Получаем параметр tags из URL (?tags=tag1,tag2)
	tagsParam := r.URL.Query().Get("tags")
	h.logger.Debug("Received tagsParam", zap.String("tagsParam", tagsParam))

	// Разбиваем строку тегов на слайс строк
	var tags []string
	if tagsParam != "" {
		for _, t := range splitAndTrim(tagsParam, ",") {
			if t != "" {
				tags = append(tags, t)
			}
		}
	}
	h.logger.Debug("Parsed tags", zap.Strings("tags", tags))

	// Получаем id задач по тегам через сервис
	h.logger.Debug("Calling service.GetTaskIDsByTags", zap.Strings("tags", tags))
	ids, err := h.service.GetTaskIDsByTags(tags)
	if err != nil {
		// Логируем ошибку получения id задач
		h.logger.Error("Ошибка получения id задач по тегам", zap.Error(err), zap.Strings("tags", tags))
		http.Error(w, "Failed to get task ids", http.StatusInternalServerError)
		return
	}
	// Логируем успешное получение id задач
	h.logger.Info("Task IDs successfully received", zap.Ints("ids", ids))

	// Формируем ответ
	resp := map[string]interface{}{
		"ids": ids,
	}
	w.Header().Set("Content-Type", "application/json")

	// Кодируем и отправляем ответ в формате JSON
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.logger.Error("Ошибка кодирования ответа в JSON", zap.Error(err))
	} else {
		h.logger.Debug("Response sent", zap.Any("response", resp))
	}
}
