package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"go.uber.org/zap"
)

func (s *Service) TaskAdd(telegramID int, args string) (string, error) {
	s.Logger.Info("TaskAdd called", zap.Int("telegram_id", telegramID), zap.String("args", args))

	requiredFields := []string{"question", "answer", "tags"}
	taskData, err := parseTaskMessage(args, requiredFields)
	if err != nil {
		s.Logger.Error("TaskAdd: parseTaskMessage error", zap.Error(err))
		return "", err
	}
	fmt.Printf("\n\nTaskAdd: parsed taskData: %+v\n", taskData)

	jsonBody, err := json.Marshal(taskData)
	if err != nil {
		s.Logger.Error("TaskAdd: json.Marshal error", zap.Error(err))
		return "", err
	}

	fmt.Printf("\n\nTaskAdd: marshaled jsonBody: %s\n", jsonBody)

	url := fmt.Sprintf("%s/taskadd?telegram_id=%d", s.CoreServiceUrl, telegramID)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		s.Logger.Error("TaskAdd: http.Do error", zap.Error(err))
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		s.Logger.Error("TaskAdd: non-200 response", zap.Int("status", resp.StatusCode), zap.ByteString("body", body))
		return "", fmt.Errorf("core-service error: %s", string(body))
	}

	var res commonSuccessResponse
	if err := json.Unmarshal(body, &res); err != nil {
		s.Logger.Error("TaskAdd: json.Unmarshal error", zap.Error(err))
		return "", err
	}

	s.Logger.Info("TaskAdd: success", zap.String("result", res.Result))
	return res.Result, nil
}
