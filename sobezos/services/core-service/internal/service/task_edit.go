package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"go.uber.org/zap"
)

// SendEditTask отправляет PUT-запрос на theory-service для редактирования задачи
func (s *Service) TaskEdit(telegramID int, id string, question string, answer string, tags []string) (string, error) {
	url := "http://theory-service:8081/taskedit"

	// Проверяем обязательное поле id
	if id == "" {
		s.logger.Error("Не указан id задачи для редактирования")
		return "📝Для редактирования задачи необходимо указать id", nil
	}

	// Преобразовать строку в число
	intid, err := strconv.Atoi(id)
	if err != nil {
		return "", fmt.Errorf("invalid id format: %v", err)
	}

	// Формируем JSON только с переданными полями
	requestData := map[string]interface{}{
		"id": intid,
	}

	// Добавляем опциональные поля, если они не пустые
	if question != "" {
		requestData["question"] = question
	}
	if answer != "" {
		requestData["answer"] = answer
	}
	if len(tags) > 0 {
		requestData["tags"] = tags
	}

	newJson, err := json.Marshal(requestData)
	if err != nil {
		s.logger.Error("Ошибка формирования задачи", zap.Error(err))
		return "⚠️Ошибка формирования задачи", err
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(newJson))
	if err != nil {
		s.logger.Error("Ошибка формирования запроса к theory-service", zap.Error(err))
		return "⚠️Ошибка формирования запроса к theory\\-service", err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		s.logger.Error("Ошибка запроса к theory-service", zap.Error(err))
		return "⚠️Ошибка запроса к theory\\-service", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return "✅Задача успешно обновлена", nil
	}

	respMsg, _ := io.ReadAll(resp.Body)
	s.logger.Error("theory-service вернул ошибку при редактировании задачи", zap.Int("status", resp.StatusCode), zap.String("body", string(respMsg)))
	return "⚠️Ошибка: " + string(respMsg), nil
}
