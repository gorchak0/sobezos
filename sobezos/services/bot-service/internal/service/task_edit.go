package service

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"go.uber.org/zap"
)

// SendEditTask отправляет PUT-запрос на theory-service для редактирования задачи
func (s *Service) TaskEdit(telegramID int, jsonText string) (string, error) {
	url := "http://theory-service:8081/taskedit"

	// Проверяем наличие id в полученном json
	var reqBody map[string]interface{}
	if err := json.Unmarshal([]byte(jsonText), &reqBody); err != nil {
		s.logger.Error("Некорректный JSON задачи", zap.Error(err))
		return "Некорректный JSON", err
	}
	if _, ok := reqBody["id"]; !ok {
		s.logger.Error("Не указан id задачи для редактирования")
		return "Для редактирования задачи необходимо указать id", nil
	}

	// Можно добавить тег пользователя, как в TaskAdd (опционально)
	userInfo, exists := s.UserCheck(telegramID)
	var tag string
	if exists && userInfo.Username != "" {
		tag = "@" + userInfo.Username
	} else {
		tag = "user:" + strconv.Itoa(telegramID)
	}
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

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(newJson))
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
	if resp.StatusCode == http.StatusOK {
		return "Задача успешно обновлена", nil
	}
	respMsg, _ := io.ReadAll(resp.Body)
	s.logger.Error("theory-service вернул ошибку при редактировании задачи", zap.Int("status", resp.StatusCode), zap.String("body", string(respMsg)))
	return "Ошибка: " + string(respMsg), nil
}
