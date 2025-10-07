package service

import (
	"fmt"
	"sobezos/services/core-service/internal/models"
	"strings"
)

// TagClear удаляет все теги из состояния пользователя
func (s *Service) TagClear(telegramID int) (string, error) {
	patch := models.UserState{
		TheoryTags:           []string{},
		CompletedTheoryTasks: []string{},
	}

	fmt.Printf("\n\n\n - bot-service - TagClear patch: %+v\n\n\n", patch)
	err := s.UserStateEdit(telegramID, patch)
	if err != nil {
		return "⚠️Ошибка очистки тегов", err
	}

	// Повторно читаем состояние пользователя
	updatedState, err := s.UserStateGet(telegramID)
	var allTags []string
	var completedTasks []string

	if err == nil && updatedState != nil {
		allTags = updatedState.TheoryTags
		completedTasks = updatedState.CompletedTheoryTasks
	}

	return fmt.Sprintf(
		"🧹Список ваших тэгов и список пройденных задач успешно очищены\n\nВаши тэги\\: %s \nВаши пройденные задачи\\: %s",
		strings.Join(allTags, ", "),
		strings.Join(completedTasks, ", "),
	), nil
}
