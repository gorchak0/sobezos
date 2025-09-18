package service

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"sobezos/services/bot-service/internal/models"

	"go.uber.org/zap"
)

func (s *Service) TaskGet(telegramID int) (string, error) {
	// Получаем тэги пользователя из user-service
	state, err := s.UserStateGet(telegramID)
	var tags []string
	if err == nil && state != nil {
		tags = state.TheoryTags
	}

	// Формируем query-параметр
	tagsParam := ""
	if len(tags) > 0 {
		tagsParam = "?tags=" + strings.Join(tags, ",")
	}
	url := "http://theory-service:8081/taskget" + tagsParam

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

	//fmt.Printf("\n\n\n!!!Сохраняем состояние пользователя: %d, %d\n\n\n\n", state.UserID, telegramID)

	s.UserStateEdit(telegramID, models.UserState{
		UserID:           state.UserID,
		LastTheoryTaskID: task.ID,
		LastAction:       "get_task",
		LastTheoryAnswer: "", //
	})

	tagsText := ""
	if len(task.Tags) > 0 {
		tagsText = "Теги: " + strings.Join(task.Tags, ", ") + "\n"
	}
	return "Задача №" + strconv.Itoa(task.ID) + ":\n" + task.Question + "\n" + tagsText, nil
}
