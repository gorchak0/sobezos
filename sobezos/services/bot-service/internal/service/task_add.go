package service

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"go.uber.org/zap"
)

// SendAddTask отправляет задачу в theory-service
func (s *Service) TaskAdd(telegramID int, args string) (string, error) {
	// Парсим JSON задачи
	var reqBody map[string]interface{}
	if err := json.Unmarshal([]byte(args), &reqBody); err != nil {
		s.logger.Error("Некорректный JSON задачи", zap.Error(err))
		return "Некорректный JSON задачи", err
	}

	// Получаем username пользователя
	userInfo, exists := s.UserCheck(telegramID)
	var tag string
	if exists && userInfo.Username != "" {
		tag = "@" + userInfo.Username
	} else {
		tag = "user:" + strconv.Itoa(telegramID)
	}

	// Добавляем тег
	tagsIface, ok := reqBody["tags"].([]interface{})
	var tags []interface{}
	if ok {
		tags = tagsIface
	} else {
		tags = []interface{}{}
	}
	tags = append(tags, tag)
	reqBody["tags"] = tags

	// Формируем JSON
	newJson, err := json.Marshal(reqBody)
	if err != nil {
		s.logger.Error("Ошибка формирования задачи", zap.Error(err))
		return "Ошибка формирования задачи", err
	}

	// Эндпоинт
	url := "http://theory-service:8081/createtask"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(newJson))
	if err != nil {
		s.logger.Error("Ошибка формирования запроса к theory-service", zap.Error(err))
		return "Ошибка формирования запроса к theory-service", err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		s.logger.Error("Ошибка запроса к theory-service", zap.Error(err))
		return "Ошибка запроса к theory-service", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated {
		// Можно сохранить состояние пользователя, если нужно
		return "Задача успешно добавлена", nil
	}
	respMsg, _ := io.ReadAll(resp.Body)
	s.logger.Error("theory-service вернул ошибку при добавлении задачи", zap.Int("status", resp.StatusCode), zap.String("body", string(respMsg)))
	return "Ошибка: " + string(respMsg), ErrServiceUnavailable
}
