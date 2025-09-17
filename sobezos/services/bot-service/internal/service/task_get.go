package service

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

func (s *Service) TaskGet(telegramID int) (string, error) {
	//эндпоинт
	url := "http://theory-service:8081/task"

	resp, err := http.Get(url)
	if err != nil {
		s.logger.Error("Не удалось выполнить GET-запрос к theory-service", zap.Error(err))
		return "⚠️ Сервис временно недоступен", err
	}
	defer resp.Body.Close()

	//обработка ошибок
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		s.logger.Error("theory-service вернул ошибочный статус", zap.Int("status", resp.StatusCode), zap.String("body", string(body)))
		return "⚠️ Сервис временно недоступен", ErrServiceUnavailable
	}

	//чтение тела
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.Error("Не удалось прочитать тело ответа от theory-service", zap.Error(err))
		return "⚠️ Сервис временно недоступен", ErrServiceUnavailable
	}

	//получаем содержимое
	var task TaskResponse
	if err := json.Unmarshal(body, &task); err != nil {
		s.logger.Error("Не удалось распарсить JSON от theory-service", zap.Error(err), zap.String("raw_body", string(body)))
		return "⚠️ Сервис временно недоступен", ErrServiceUnavailable
	}

	// Сохраняем состояние пользователя
	s.UserStateEdit(telegramID, map[string]interface{}{
		"last_theory_task_id": task.ID,
		"last_action":         "get_task",
	})
	tagsText := ""
	if len(task.Tags) > 0 {
		tagsText = "Теги: " + strings.Join(task.Tags, ", ") + "\n"
	}
	return "Задача №" + strconv.Itoa(task.ID) + ":\n" + task.Question + "\n" + tagsText, nil
}
