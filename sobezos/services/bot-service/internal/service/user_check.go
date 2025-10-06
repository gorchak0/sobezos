package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"go.uber.org/zap"
)

func (s *Service) UserCheck(telegramID int) (userInfo struct{ Username, Role string }, exists bool) {
	s.Logger.Info("UserCheck called", zap.Int("telegram_id", telegramID))
	url := fmt.Sprintf("%s/usercheck?telegram_id=%d", s.CoreServiceUrl, telegramID)
	resp, err := http.Get(url)
	if err != nil {
		s.Logger.Error("UserCheck: http.Get error", zap.Error(err))
		return userInfo, false
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode == 404 {
		s.Logger.Warn("UserCheck: user not found", zap.Int("telegram_id", telegramID))
		return userInfo, false
	}
	if resp.StatusCode != 200 {
		s.Logger.Error("UserCheck: non-200 response", zap.Int("status", resp.StatusCode), zap.ByteString("body", body))
		return userInfo, false
	}
	var res userCheckResponse
	if err := json.Unmarshal(body, &res); err != nil {
		s.Logger.Error("UserCheck: json.Unmarshal error", zap.Error(err))
		return userInfo, false
	}
	userInfo.Username = res.Username
	userInfo.Role = res.Role
	s.Logger.Info("UserCheck: user found", zap.String("username", userInfo.Username), zap.String("role", userInfo.Role))
	return userInfo, true
}
