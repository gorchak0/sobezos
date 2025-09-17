package service

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func (s *Service) UserAdd(TelegramID int, args string) (string, error) {

	parts := strings.Fields(args)
	if len(parts) < 2 {
		return "Используйте: /add <telegram_id> <username>", ErrServiceUnavailable
	}
	newID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return "Некорректный telegram_id", ErrServiceUnavailable
	}
	username := parts[1]

	url := "http://user-service:8082/users/add"
	reqBody := map[string]interface{}{
		"telegram_id": newID,
		"username":    username,
	}
	body, _ := json.Marshal(reqBody)
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return "Ошибка запроса к user-service", ErrServiceUnavailable
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Admin-Telegram-ID", strconv.Itoa(TelegramID))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "Ошибка запроса к user-service", ErrServiceUnavailable
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusCreated {
		return "Пользователь успешно добавлен", ErrServiceUnavailable
	}
	respMsg, _ := io.ReadAll(resp.Body)
	return "Ошибка: " + string(respMsg), ErrServiceUnavailable
}
