package service

import (
	"encoding/json"
	"net/http"
	"sobezos/services/core-service/internal/models"
	"strconv"
)

// UserState API helpers
func (s *Service) UserStateGet(userID int) (*models.UserState, error) {
	url := "http://user-service:8082/userstateget?user_id=" + strconv.Itoa(userID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, nil
	}
	var state models.UserState
	if err := json.NewDecoder(resp.Body).Decode(&state); err != nil {
		return nil, err
	}
	return &state, nil
}
