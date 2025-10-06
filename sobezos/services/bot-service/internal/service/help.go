package service

import "go.uber.org/zap"

func (s *Service) Help() (string, error) {
	s.Logger.Info("Help called")
	resp := `
	🔸🔸🔸*Доступно всем пользователям*

	🔹/taskget 
		Получить случайную задачу \(с учётом ваших тегов, если заданы\)

	🔹/taskgetid \<id\> 
		Получить задачу по её номеру \(id\)

	🔹/answerget 
		Получить ответ на последнюю полученную задачу

	🔹/statsget 
		Показать статистику \(пока не реализовано\)

	🔹/tagset \<тег1, тег2,\.\.\.\> 
		Добавить теги для фильтрации задач \(например\: /tagset динамика, графы\)

	🔹/tagclear 
		Очистить все ваши теги

	🔹/tagget 
		Показать список всех доступных тегов с описаниями

	🔹/taskadd 
		Добавить новую задачу\. Пример синтаксиса\:
	
	/taskadd
	$tags
	тэг1
	$question
	Вопрос
	$answer
	Ответ

	🔹/taskedit
		Редактировать задачу\. Пример синтаксиса\:

	/taskadd
	$id
	4
	$tags
	тэг1
	$question
	Вопрос
	$answer
	Ответ

	🔸🔸🔸*Только для администраторов*
	
	🔹/useradd \<id\> \<username\> 
	Добавить пользователя
	

	`

	s.Logger.Info("\n Help text ", zap.String("text", resp))
	return resp, nil
}
