package service

import (
	"fmt"
	"strings"
)

func (s *Service) TagSet(telegramID int, args string) (string, error) {
	// Установить theory_tags для пользователя
	tags := strings.Split(args, ",")
	for i := range tags {
		tags[i] = strings.TrimSpace(tags[i])
	}

	// Сохраняем состояние пользователя
	s.UserStateEdit(telegramID, map[string]interface{}{
		"theory_tags": tags,
	})

	result, err := s.UserStateGet(telegramID)
	if err != nil {
		//
	}

	return fmt.Sprintf("Тэги успешно устанвлены \n %v", result), nil
}
