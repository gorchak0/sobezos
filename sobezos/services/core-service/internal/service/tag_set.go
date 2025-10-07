package service

import (
	"fmt"
	"sobezos/services/core-service/internal/models"
	"strings"

	"go.uber.org/zap"
)

func (s *Service) TagSet(telegramID int, args string) (string, error) {
	state, err := s.UserStateGet(telegramID)
	if err != nil {
		s.logger.Error("Ошибка при получении состояния пользователя", zap.Int("telegramID", telegramID), zap.Error(err))
		return "", err
	}

	currentTags := make([]string, 0)
	completedTasks := make([]string, 0)
	if state != nil {
		currentTags = state.TheoryTags
		completedTasks = state.CompletedTheoryTasks
	}

	// Парсим входные теги
	rawTags := strings.Split(args, ",")
	newTags := make([]string, 0, len(rawTags))
	for _, t := range rawTags {
		if tag := strings.TrimSpace(t); tag != "" {
			newTags = append(newTags, tag)
		}
	}

	// Если пользователь не передал новые теги
	if len(newTags) == 0 {
		return fmt.Sprintf(
			"✨Новые тэги\\: отсутствуют\nВсе ваши тэги\\: %s\nПройденные задачи\\: %v шт\\.",
			strings.Join(currentTags, ", "),
			len(completedTasks),
		), nil
	}

	// Объединяем старые и новые теги в set
	tagSet := make(map[string]struct{}, len(currentTags)+len(newTags))
	for _, t := range currentTags {
		tagSet[t] = struct{}{}
	}

	hasNewTags := false
	for _, t := range newTags {
		if _, exists := tagSet[t]; !exists {
			hasNewTags = true
			tagSet[t] = struct{}{}
		}
	}

	// Если ничего нового нет, просто возвращаем текущее состояние
	if !hasNewTags {
		return fmt.Sprintf(
			"✨Новые тэги\\: нет новых\nВсе ваши тэги\\: %s\nПройденные задачи\\: %v шт\\.",
			strings.Join(currentTags, ", "),
			len(completedTasks),
		), nil
	}

	// Преобразуем map обратно в slice
	mergedTags := make([]string, 0, len(tagSet))
	for t := range tagSet {
		mergedTags = append(mergedTags, t)
	}

	// Обновляем состояние пользователя
	updateData := models.UserState{
		TheoryTags: mergedTags,
	}
	completedTasksCount := len(completedTasks)
	if completedTasksCount > 0 {
		updateData.CompletedTheoryTasks = []string{} // очищаем задачи при добавлении новых тегов
	}

	if err := s.UserStateEdit(telegramID, updateData); err != nil {
		s.logger.Error("⚠️Ошибка при обновлении состояния пользователя", zap.Int("telegramID", telegramID), zap.Error(err))
		return "", err
	}

	// Формируем финальное сообщение
	completedMsg := "0 шт\\."
	if completedTasksCount > 0 {
		completedMsg = fmt.Sprintf("0 шт\\. \\(очищено %d задач\\)", completedTasksCount)
	}

	return fmt.Sprintf(
		"✨Новые тэги\\: %s\nВсе ваши тэги\\: %s\nПройденные задачи\\: %s",
		strings.Join(newTags, ", "),
		strings.Join(mergedTags, ", "),
		completedMsg,
	), nil
}
