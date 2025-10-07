package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"go.uber.org/zap"
)

func (s *Service) StatsGet(telegramID int) (string, error) {
	s.Logger.Info("AnswerGet called", zap.Int("telegram_id", telegramID))
	url := fmt.Sprintf("%s/statsget?telegram_id=%d", s.CoreServiceUrl, telegramID)
	resp, err := http.Get(url)
	if err != nil {
		s.Logger.Error("AnswerGet: http.Get error", zap.Error(err))
		return "", err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		s.Logger.Error("AnswerGet: non-200 response", zap.Int("status", resp.StatusCode), zap.ByteString("body", body))
		return "", fmt.Errorf("core-service error: %s", string(body))
	}
	var res commonSuccessResponse

	s.Logger.Info("AnswerGet: raw body", zap.ByteString("body", body))

	if err := json.Unmarshal(body, &res); err != nil {
		s.Logger.Error("AnswerGet: json.Unmarshal error", zap.Error(err))
		return "", err
	}
	s.Logger.Info("AnswerGet: success", zap.String("result", res.Result))
	return res.Result, nil
}
