package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sobezos/services/core-service/internal/models"
	"strconv"

	"go.uber.org/zap"
)

func (s *Service) AnswerGet(telegramID int) (string, error) {
	// Получить состояние пользователя
	state, err := s.UserStateGet(telegramID)
	if err != nil || state == nil {
		s.logger.Error("Не удалось получить состояние пользователя", zap.Int("telegramID", telegramID), zap.Error(err))
		return "❌Нет информации о последней задаче", ErrServiceUnavailable
	}

	taskID := state.LastTheoryTaskID
	if taskID == nil || *taskID == 0 {
		s.logger.Warn("Нет информации о последней задаче", zap.Int("telegramID", telegramID))
		return "❌Нет информации о последней задаче", ErrServiceUnavailable
	}

	// Получение текста ответа
	url := fmt.Sprintf("http://theory-service:8081/answerget?task_id=%d", *taskID)
	resp, err := http.Get(url)
	if err != nil {
		s.logger.Error("Ошибка при запросе к theory-service", zap.String("url", url), zap.Error(err))
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		s.logger.Error("theory-service вернул ошибочный статус", zap.Int("status", resp.StatusCode), zap.String("url", url))
		return "", nil
	}

	var res struct {
		Answer string `json:"answer"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		s.logger.Error("Ошибка декодирования ответа от theory-service", zap.Error(err))
		return "", err
	}

	lastAction := "get_task"

	// Добавляем задачу в список решенных
	completedTasks := state.CompletedTheoryTasks
	if completedTasks == nil {
		completedTasks = []string{}
	}

	// Проверяем, нет ли уже этой задачи в списке решенных
	taskIDStr := strconv.Itoa(*taskID)
	taskExists := false
	for _, task := range completedTasks {
		if task == taskIDStr {
			taskExists = true
			break
		}
	}

	// Если задачи еще нет в списке решенных - добавляем
	if !taskExists {
		completedTasks = append(completedTasks, taskIDStr)
	}

	// Обновляем состояние пользователя
	err = s.UserStateEdit(telegramID, models.UserState{
		LastAction:           &lastAction,
		LastTheoryTaskID:     taskID,
		CompletedTheoryTasks: completedTasks,
	})
	if err != nil {
		s.logger.Error("Ошибка при обновлении состояния пользователя", zap.Int("telegramID", telegramID), zap.Error(err))
		return "", err
	}

	result := fmt.Sprintf("💡Ответ на вопрос №%d: \n%s", *taskID, res.Answer)
	s.logger.Info("Ответ успешно получен",
		zap.Int("telegramID", telegramID),
		zap.Int("taskID", *taskID),
		zap.Strings("completedTasks", completedTasks))

	return result, nil
}
