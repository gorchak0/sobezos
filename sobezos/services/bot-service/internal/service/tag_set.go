package service

import (
	"fmt"
	"strings"
)

func (s *Service) TagSet(telegramID int, args string) (string, error) {
	// Получаем текущие тэги пользователя
	state, err := s.UserStateGet(telegramID)
	var currentTags []string
	if err == nil && state != nil {
		if tagsIface, ok := (*state)["theory_tags"]; ok {
			if tagsSlice, ok := tagsIface.([]interface{}); ok {
				for _, t := range tagsSlice {
					if str, ok := t.(string); ok {
						currentTags = append(currentTags, str)
					}
				}
			} else if tagsSlice, ok := tagsIface.([]string); ok {
				currentTags = tagsSlice
			}
		}
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
	s.UserStateEdit(telegramID, map[string]interface{}{
		"theory_tags": mergedTags,
	})

	return fmt.Sprintf("Тэги успешно добавлены: %v", mergedTags), nil
}
