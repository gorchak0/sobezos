package service

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"go.uber.org/zap"
)

// TaskGetByID получает задачу по id из theory-service
func (s *Service) TaskGetID(telegramID int, args string) (string, error) {
	//получаем id задачи из args
	taskID, err := strconv.Atoi(args)
	if err != nil {
		s.logger.Error("Некорректный id задачи", zap.Error(err))
		return "Некорректный id задачи. Используйте /taskgetid <id_задачи>", nil
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

	tagsText := ""
	if len(task.Tags) > 0 {
		tagsText = "Теги: " + "\n" + "- " + (func(tags []string) string {
			res := ""
			for i, t := range tags {
				if i > 0 {
					res += ", "
				}
				res += t
			}
			return res
		})(task.Tags) + "\n"
	}
	return "Задача №" + strconv.Itoa(task.ID) + ":\n" + task.Question + "\n" + tagsText, nil
}
