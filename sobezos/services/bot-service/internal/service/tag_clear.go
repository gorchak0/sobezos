package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"go.uber.org/zap"
)

func (s *Service) TagClear(telegramID int) (string, error) {
	s.Logger.Info("TagClear called", zap.Int("telegram_id", telegramID))
	url := fmt.Sprintf("%s/tagclear", s.CoreServiceUrl)
	reqBody, _ := json.Marshal(map[string]int{"telegram_id": telegramID})
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		s.Logger.Error("TagClear: http.Do error", zap.Error(err))
		return "", err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		s.Logger.Error("TagClear: non-200 response", zap.Int("status", resp.StatusCode), zap.ByteString("body", body))
		return "", fmt.Errorf("core-service error: %s", string(body))
	}
	var res commonSuccessResponse
	if err := json.Unmarshal(body, &res); err != nil {
		s.Logger.Error("TagClear: json.Unmarshal error", zap.Error(err))
		return "", err
	}
	s.Logger.Info("TagClear: success", zap.String("result", res.Result))
	return res.Result, nil
}
