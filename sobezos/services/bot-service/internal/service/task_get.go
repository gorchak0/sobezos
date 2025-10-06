package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"go.uber.org/zap"
)

func (s *Service) TaskGet(telegramID int) (string, error) {
	s.Logger.Info("TaskGet called", zap.Int("telegram_id", telegramID))
	url := fmt.Sprintf("%s/taskget?telegram_id=%d", s.CoreServiceUrl, telegramID)
	resp, err := http.Get(url)
	if err != nil {
		s.Logger.Error("TaskGet: http.Get error", zap.Error(err))
		return "", err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		s.Logger.Error("TaskGet: non-200 response", zap.Int("status", resp.StatusCode), zap.ByteString("body", body))
		return "", fmt.Errorf("core-service error: %s", string(body))
	}
	var res commonSuccessResponse
	if err := json.Unmarshal(body, &res); err != nil {
		s.Logger.Error("TaskGet: json.Unmarshal error", zap.Error(err))
		return "", err
	}
	s.Logger.Info("TaskGet: success", zap.String("result", res.Result))
	return res.Result, nil
}
