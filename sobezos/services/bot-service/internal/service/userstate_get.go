package service

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// UserState API helpers
func (s *Service) UserStateGet(userID int) (*map[string]interface{}, error) {
	url := "http://user-service:8082/states/" + strconv.Itoa(userID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, nil
	}
	var state map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&state); err != nil {
		return nil, err
	}
	return &state, nil
}
