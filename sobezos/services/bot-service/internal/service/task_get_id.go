package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"go.uber.org/zap"
)

func (s *Service) TaskGetID(telegramID int, args string) (string, error) {
	s.Logger.Info("TaskGetID called", zap.Int("telegram_id", telegramID), zap.String("args", args))
	url := fmt.Sprintf("%s/taskgetid?telegram_id=%d&args=%s", s.CoreServiceUrl, telegramID, args)
	resp, err := http.Get(url)
	if err != nil {
		s.Logger.Error("TaskGetID: http.Get error", zap.Error(err))
		return "", err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		s.Logger.Error("TaskGetID: non-200 response", zap.Int("status", resp.StatusCode), zap.ByteString("body", body))
		return "", fmt.Errorf("core-service error: %s", string(body))
	}
	var res commonSuccessResponse
	if err := json.Unmarshal(body, &res); err != nil {
		s.Logger.Error("TaskGetID: json.Unmarshal error", zap.Error(err))
		return "", err
	}
	s.Logger.Info("TaskGetID: success", zap.String("result", res.Result))
	return res.Result, nil
}
