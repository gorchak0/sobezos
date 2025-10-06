package service

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"sobezos/services/core-service/internal/models"

	"go.uber.org/zap"
)

func (s *Service) TaskGet(telegramID int) (string, error) {
	// Получаем состояние пользователя из user-service
	state, err := s.UserStateGet(telegramID)
	var tags []string
	var completedTasks []string
	if err == nil && state != nil {
		tags = state.TheoryTags
		completedTasks = state.CompletedTheoryTasks
	}

	// Формируем query-параметры
	params := url.Values{}
	if len(tags) > 0 {
		params.Add("tags", strings.Join(tags, ","))
	}
	if len(completedTasks) > 0 {
		params.Add("exclude", strings.Join(completedTasks, ","))
	}

	url := "http://theory-service:8081/taskget"
	if len(params) > 0 {
		url += "?" + params.Encode()
	}

	resp, err := http.Get(url)
	if err != nil {
		s.logger.Error("Не удалось выполнить GET-запрос к theory-service", zap.Error(err))
		return "⚠️ Сервис временно недоступен", err
	}
	defer resp.Body.Close()

	// Обработка ошибок
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		s.logger.Error("theory-service вернул ошибочный статус", zap.Int("status", resp.StatusCode), zap.String("body", string(body)))
		return "⚠️ Сервис временно недоступен", ErrServiceUnavailable
	}

	// Чтение тела
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.Error("Не удалось прочитать тело ответа от theory-service", zap.Error(err))
		return "⚠️ Сервис временно недоступен", ErrServiceUnavailable
	}

	// Получаем содержимое
	var task TaskResponse
	if err := json.Unmarshal(body, &task); err != nil {
		s.logger.Error("Не удалось распарсить JSON от theory-service", zap.Error(err), zap.String("raw_body", string(body)))
		return "⚠️ Сервис временно недоступен", ErrServiceUnavailable
	}

	// Добавляем текущую задачу в массив решенных задач
	if completedTasks == nil {
		completedTasks = []string{}
	}

	// Преобразуем ID задачи в строку и добавляем в массив, если его там еще нет
	taskIDStr := strconv.Itoa(task.ID)
	taskExists := false
	for _, completedTask := range completedTasks {
		if completedTask == taskIDStr {
			taskExists = true
			break
		}
	}

	if !taskExists {
		completedTasks = append(completedTasks, taskIDStr)
	}
	lastAction := "get_task"

	// Сохраняем обновленное состояние пользователя
	err = s.UserStateEdit(telegramID, models.UserState{
		UserID:               int64(telegramID),
		LastTheoryTaskID:     &task.ID,
		LastAction:           &lastAction,
		CompletedTheoryTasks: completedTasks,
	})

	if err != nil {
		s.logger.Error("Не удалось обновить состояние пользователя", zap.Error(err))
		// Продолжаем выполнение, так как основная задача уже получена
	}

	tagsText := ""
	if len(task.Tags) > 0 {
		tagsText = "Теги: " + strings.Join(task.Tags, ", ") + "\n"
	}
	return "Задача №" + strconv.Itoa(task.ID) + ":\n" + task.Question + "\n" + tagsText, nil
}
