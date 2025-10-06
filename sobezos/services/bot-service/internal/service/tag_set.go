package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"go.uber.org/zap"
)

func (s *Service) TagSet(telegramID int, args string) (string, error) {
	s.Logger.Info("TagSet called", zap.Int("telegram_id", telegramID), zap.String("args", args))

	url := fmt.Sprintf("%s/tagset?telegram_id=%d", s.CoreServiceUrl, telegramID)

	reqBody, _ := json.Marshal(map[string]string{"args": args})
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		s.Logger.Error("TagSet: http.Do error", zap.Error(err))
		return "", err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		s.Logger.Error("TagSet: non-200 response", zap.Int("status", resp.StatusCode), zap.ByteString("body", body))
		return "", fmt.Errorf("core-service error: %s", string(body))
	}
	var res commonSuccessResponse
	if err := json.Unmarshal(body, &res); err != nil {
		s.Logger.Error("TagSet: json.Unmarshal error", zap.Error(err))
		return "", err
	}
	s.Logger.Info("TagSet: success", zap.String("result", res.Result))
	return res.Result, nil
}
