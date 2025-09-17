package service

import (
	"fmt"
	"sobezos/services/bot-service/internal/models"
	"strings"
)

func (s *Service) TagSet(telegramID int, args string) (string, error) {
	// Получаем текущие тэги пользователя
	state, err := s.UserStateGet(telegramID)
	var currentTags []string
	if err == nil && state != nil {
		currentTags = state.TheoryTags
	}

	// Парсим новые тэги
	newTags := strings.Split(args, ",")
	for i := range newTags {
		newTags[i] = strings.TrimSpace(newTags[i])
	}

	// Добавляем только уникальные новые тэги
	tagSet := make(map[string]struct{})
	for _, t := range currentTags {
		tagSet[t] = struct{}{}
	}
	for _, t := range newTags {
		if t != "" {
			tagSet[t] = struct{}{}
		}
	}
	mergedTags := make([]string, 0, len(tagSet))
	for t := range tagSet {
		mergedTags = append(mergedTags, t)
	}

	// Сохраняем состояние пользователя
	s.UserStateEdit(telegramID, models.UserState{
		TheoryTags: mergedTags,
	})

	// Повторно читаем состояние пользователя
	updatedState, err := s.UserStateGet(telegramID)
	var allTags []string
	if err == nil && updatedState != nil {
		allTags = updatedState.TheoryTags
	}

	return fmt.Sprintf(
		"Тэги успешно добавлены: %s\nВсе ваши тэги: %s",
		strings.Join(newTags, ", "),
		strings.Join(allTags, ", "),
	), nil
}
