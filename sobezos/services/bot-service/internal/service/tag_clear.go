package service

import (
	"fmt"
	"sobezos/services/bot-service/internal/models"
	"strings"
)

// TagClear удаляет все теги из состояния пользователя
func (s *Service) TagClear(telegramID int) (string, error) {
	patch := models.UserState{
		TheoryTags: []string{},
	}

	fmt.Printf("\n\n\n - bot-service - TagClear patch: %+v\n\n\n", patch)
	err := s.UserStateEdit(telegramID, patch)
	if err != nil {
		return "Ошибка очистки тегов", err
	}

	// Повторно читаем состояние пользователя
	updatedState, err := s.UserStateGet(telegramID)
	var allTags []string
	if err == nil && updatedState != nil {
		allTags = updatedState.TheoryTags
	}

	return fmt.Sprintf(
		"Тэги успешно очищены\nВсе ваши тэги: %s",
		strings.Join(allTags, ", "),
	), nil
}
