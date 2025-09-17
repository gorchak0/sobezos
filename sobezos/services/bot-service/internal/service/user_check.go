package service

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func (s *Service) UserCheck(telegramID int) (UserCheckResponse, bool) {
	url := "http://user-service:8082/users/check?telegram_id=" + strconv.Itoa(telegramID)
	resp, err := http.Get(url)
	if err != nil {
		return UserCheckResponse{}, false
	}
	defer resp.Body.Close()
	var res UserCheckResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return UserCheckResponse{}, false
	}
	return res, res.Exists
}
