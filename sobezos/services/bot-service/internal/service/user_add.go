package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"go.uber.org/zap"
)

func (s *Service) UserAdd(telegramID int, args string) (string, error) {
	s.Logger.Info("UserAdd called", zap.Int("admin_telegram_id", telegramID), zap.String("args", args))
	url := fmt.Sprintf("%s/useradd?telegram_id=%d", s.CoreServiceUrl, telegramID)
	reqBody, _ := json.Marshal(map[string]string{"args": args})
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		s.Logger.Error("UserAdd: http.Do error", zap.Error(err))
		return "", err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 { //
		s.Logger.Error("UserAdd: non-200 response", zap.Int("status", resp.StatusCode), zap.ByteString("body", body))

		// Парсим JSON и извлекаем только сообщение об ошибке
		var errorResp struct {
			Error string `json:"error"`
		}
		if err := json.Unmarshal(body, &errorResp); err == nil && errorResp.Error != "" {
			return "", fmt.Errorf(errorResp.Error) // Вернет только "Используйте: /add <telegram\_id> <username> <role>"
		}

		return "", fmt.Errorf("service unavailable")
	}
	var res commonSuccessResponse
	if err := json.Unmarshal(body, &res); err != nil {
		s.Logger.Error("UserAdd: json.Unmarshal error", zap.Error(err))
		return "", err
	}
	s.Logger.Info("UserAdd: success", zap.String("result", res.Result))
	return res.Result, nil
}
