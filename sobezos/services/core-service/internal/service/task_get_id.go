package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sobezos/services/core-service/internal/models"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

// TaskGetByID получает задачу по id из theory-service
func (s *Service) TaskGetID(telegramID int, args string) (string, error) {
	// Получаем id задачи из args
	taskID, err := strconv.Atoi(args)
	if err != nil {
		s.logger.Error("Некорректный id задачи", zap.Error(err))
		return "❌Некорректный id задачи\\. Используйте /taskgetid <id\\_задачи>", nil
	}

	// Сначала получаем текущее состояние пользователя
	userState, err := s.UserStateGet(int(telegramID))
	if err != nil {
		s.logger.Error("Не удалось получить состояние пользователя", zap.Error(err))
		return "⚠️ Ошибка получения состояния пользователя", err
	}

	url := "http://theory-service:8081/taskgetid?task_id=" + strconv.Itoa(taskID)

	resp, err := http.Get(url)
	if err != nil {
		s.logger.Error("Не удалось выполнить GET-запрос к theory-service", zap.Error(err))
		return "⚠️ Сервис временно недоступен", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		s.logger.Error("theory-service вернул ошибочный статус", zap.Int("status", resp.StatusCode), zap.String("body", string(body)))
		return "⚠️ Сервис временно недоступен", ErrServiceUnavailable
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.Error("Не удалось прочитать тело ответа от theory-service", zap.Error(err))
		return "⚠️ Сервис временно недоступен", ErrServiceUnavailable
	}

	var task TaskResponse
	if err := json.Unmarshal(body, &task); err != nil {
		s.logger.Error("Не удалось распарсить JSON от theory-service", zap.Error(err), zap.String("raw_body", string(body)))
		return "⚠️ Сервис временно недоступен", ErrServiceUnavailable
	}

	if task.Exist == 0 {
		return "❌Задача с таким id не найдена", nil
	}

	// Добавляем задачу в completed_theory_tasks, если её там еще нет
	completedTasks := make([]string, len(userState.CompletedTheoryTasks))
	copy(completedTasks, userState.CompletedTheoryTasks)

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

	// Обновляем состояние пользователя
	s.UserStateEdit(int(telegramID), models.UserState{
		UserID:               int64(telegramID),
		LastTheoryTaskID:     &task.ID,
		LastAction:           &lastAction,
		CompletedTheoryTasks: completedTasks,
	})

	var tagsText string
	if len(task.Tags) > 0 {
		tagsText = "Теги:\n \\- " + strings.Join(task.Tags, ", ") + "\n"
	}

	return fmt.Sprintf("📌Задача №%d \n%s \n%s", task.ID, task.Question, tagsText), nil
}
