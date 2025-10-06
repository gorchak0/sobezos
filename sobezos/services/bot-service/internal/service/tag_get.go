package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"go.uber.org/zap"
)

func (s *Service) TagGet() (string, error) {
	s.Logger.Info("TagGet called")
	url := fmt.Sprintf("%s/tagget", s.CoreServiceUrl)
	resp, err := http.Get(url)
	if err != nil {
		s.Logger.Error("TagGet: http.Get error", zap.Error(err))
		return "", err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		s.Logger.Error("TagGet: non-200 response", zap.Int("status", resp.StatusCode), zap.ByteString("body", body))
		return "", fmt.Errorf("core-service error: %s", string(body))
	}
	var res commonSuccessResponse
	if err := json.Unmarshal(body, &res); err != nil {
		s.Logger.Error("TagGet: json.Unmarshal error", zap.Error(err))
		return "", err
	}
	s.Logger.Info("TagGet: success", zap.String("result", res.Result))
	return res.Result, nil
}
