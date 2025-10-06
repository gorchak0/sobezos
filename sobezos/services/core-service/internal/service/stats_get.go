package service

// UserState API helpers
func (s *Service) StatsGet(userID int) (string, error) {

	//запрос к theory-service для получения всех номеров задач по указанным тэгам, а также общего количества задач

	//запрос к user-service для получения user-state поля с решенными задачами

	/*


		📊 Ваша статистика:

		Всего просмотрено 2% задач (2 из 100)
		По тэгам [SQL,БД] 24% (12 из 50)
		Номера оставшихся	[2, 4, 5, 6, ...]

	*/

	return "Статистика пока не реализована", nil
}
