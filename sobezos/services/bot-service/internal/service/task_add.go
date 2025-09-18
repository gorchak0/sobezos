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

// SendAddTask отправляет задачу в theory-service
func (s *Service) TaskAdd(telegramID int, args string) (string, error) {
	// Парсим JSON задачи
	var reqBody map[string]interface{}
	if err := json.Unmarshal([]byte(args), &reqBody); err != nil {
		s.logger.Error("Некорректный JSON задачи", zap.Error(err))
		return "Некорректный JSON задачи", err
	}

	// --- Логируем распарсенный answer ---
	if ans, ok := reqBody["answer"].(string); ok {
		fmt.Printf("\n\n\nПосле Unmarshal, answer = %s\n", ans)
	} else {
		fmt.Printf("\n\n\nПосле Unmarshal, answer не найден или не строка\n")
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

	// --- Логируем JSON перед отправкой ---

	fmt.Printf("\n\n\nПеред отправкой в theory-service, JSON = %s\n", string(newJson))

	// Эндпоинт
	url := "http://theory-service:8081/taskadd"
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
		var respData struct {
			Status string `json:"status"`
			ID     int    `json:"id"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
			s.logger.Error("Ошибка декодирования ответа theory-service", zap.Error(err))
			return "Задача добавлена, но не удалось получить id", nil
		}
		return "Задача успешно добавлена. Номер задачи: " + strconv.Itoa(respData.ID), nil
	}
	respMsg, _ := io.ReadAll(resp.Body)
	s.logger.Error("theory-service вернул ошибку при добавлении задачи", zap.Int("status", resp.StatusCode), zap.String("body", string(respMsg)))
	return "Ошибка: " + string(respMsg), ErrServiceUnavailable
}
