package handler

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

// TaskGetAll обрабатывает GET /taskgetall и возвращает общее количество задач
func (h *TaskHandler) TaskGetAll(w http.ResponseWriter, r *http.Request) {
	// Логируем вызов хендлера с методом и URL
	h.logger.Info("TaskGetAll handler called", zap.String("method", r.Method), zap.String("url", r.URL.String()))

	// Проверяем, что используется метод GET
	if r.Method != http.MethodGet {
		h.logger.Warn("Method not allowed", zap.String("method", r.Method))
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Получаем количество задач через сервис
	h.logger.Debug("Calling service.GetTaskCount()")
	count, err := h.service.GetTaskCount()
	if err != nil {
		// Логируем ошибку получения количества задач
		h.logger.Error("Ошибка получения количества задач", zap.Error(err))
		http.Error(w, "Failed to get task count", http.StatusInternalServerError)
		return
	}
	// Логируем успешное получение количества задач
	h.logger.Info("Task count successfully received", zap.Int("count", count))

	// Формируем ответ
	resp := map[string]interface{}{
		"count": count,
	}
	w.Header().Set("Content-Type", "application/json")

	// Кодируем и отправляем ответ в формате JSON
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.logger.Error("Ошибка кодирования ответа в JSON", zap.Error(err))
	} else {
		h.logger.Debug("Response sent", zap.Any("response", resp))
	}
}
