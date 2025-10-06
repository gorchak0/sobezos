package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"sobezos/services/core-service/internal/models"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

func (s *Service) UserAdd(TelegramID int, args string) (string, error) {

	parts := strings.Fields(args)
	if len(parts) < 3 {
		s.logger.Warn("Недостаточно аргументов для добавления пользователя", zap.String("args", args))
		return "Используйте: /add <telegram_id> <username> <role>", errors.New("Используйте: /add <telegram_id> <username> <role>")
	}
	newID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		s.logger.Warn("Некорректный telegram_id", zap.String("telegram_id", parts[0]), zap.Error(err))
		return "Некорректный telegram_id", errors.New("Некорректный telegram_id") //
	}
	username := parts[1]
	role := parts[2]
	s.logger.Info("Попытка добавить пользователя", zap.Int64("new_telegram_id", newID), zap.String("username", username), zap.String("role", role), zap.Int("admin_telegram_id", TelegramID))

	user := models.User{
		TelegramID: newID,
		Username:   username,
		Role:       role,
	}

	url := "http://user-service:8082/useradd" //
	body, err := json.Marshal(user)
	if err != nil {
		s.logger.Error("Ошибка маршалинга пользователя", zap.Any("user", user), zap.Error(err))
		return "Ошибка формирования запроса", ErrServiceUnavailable
	}
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		s.logger.Error("Ошибка создания запроса к user-service", zap.Error(err))
		return "Ошибка запроса к user-service", ErrServiceUnavailable
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Admin-Telegram-ID", strconv.Itoa(TelegramID))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		s.logger.Error("Ошибка отправки запроса к user-service", zap.Error(err))
		return "Ошибка запроса к user-service", ErrServiceUnavailable
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusCreated {
		s.logger.Info("Пользователь успешно добавлен", zap.Any("user", user))
		return "Пользователь успешно добавлен", nil
	}
	respMsg, _ := io.ReadAll(resp.Body)
	s.logger.Warn("user-service вернул ошибку при добавлении пользователя", zap.Int("status", resp.StatusCode), zap.String("body", string(respMsg)), zap.Any("user", user))
	return "Ошибка", ErrServiceUnavailable
}
