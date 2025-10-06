package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"sobezos/services/core-service/internal/models"
)

func (s *Service) UserStateEdit(userID int, patch models.UserState) error {
	url := "http://user-service:8082/userstateedit?user_id=" + strconv.Itoa(userID)
	body, _ := json.Marshal(patch)
	fmt.Printf("\n\n\nbot-service - UserStateEdit raw JSON: %s\n\n\n", string(body))
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

//
