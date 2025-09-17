package service

// TagClear удаляет все теги из состояния пользователя
func (s *Service) TagClear(telegramID int) (string, error) {
	patch := map[string]interface{}{
		"tags": []string{},
	}
	err := s.UserStateEdit(telegramID, patch)
	if err != nil {
		return "Ошибка очистки тегов", err
	}
	return "Теги успешно очищены", nil
}
