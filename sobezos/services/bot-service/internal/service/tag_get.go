package service

import (
	"encoding/json"
	"net/http"
	"strings"
)

// TagGet получает список всех доступных тегов с описаниями из theory-service
func (s *Service) TagGet() (string, error) {
	url := "http://theory-service:8081/tagget"
	resp, err := http.Get(url)
	if err != nil {
		return "Ошибка запроса к theory-service", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "Ошибка theory-service", nil
	}
	var tags []struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tags); err != nil {
		return "Ошибка разбора ответа от theory-service", err
	}
	if len(tags) == 0 {
		return "Нет доступных тегов", nil
	}
	var sb strings.Builder
	sb.WriteString("Доступные теги:\n")
	for _, t := range tags {
		sb.WriteString("- " + t.Name)
		if t.Description != "" {
			sb.WriteString(": " + t.Description)
		}
		sb.WriteString("\n")
	}
	return sb.String(), nil
}
