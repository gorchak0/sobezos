package service

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
)

func (s *Service) UserStateEdit(userID int, patch map[string]interface{}) error {
	url := "http://user-service:8082/states/" + strconv.Itoa(userID)
	body, _ := json.Marshal(patch)
	req, err := http.NewRequest("PATCH", url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return err
	}
	return nil
}
