package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// UserState API helpers
func (s *Service) StatsGet(userID int) (string, error) {
	// Получаем user-state (теги и решённые задачи)
	state, err := s.UserStateGet(userID)
	if err != nil {
		return "⚠️Ошибка получения состояния пользователя", err
	}
	if state == nil {
		return "❌Нет данных о пользователе", nil
	}
	userTags := state.TheoryTags
	completedTasks := make(map[string]struct{}, len(state.CompletedTheoryTasks))
	for _, t := range state.CompletedTheoryTasks {
		completedTasks[t] = struct{}{}
	}

	// Получаем все id задач по тегам пользователя (task_get_tags)
	tagTaskIDs := []int{}
	if len(userTags) > 0 {
		tagsParam := strings.Join(userTags, ",")
		urlTags := fmt.Sprintf("http://theory-service:8081/taskgettags?tags=%s", tagsParam)
		respTags, err := http.Get(urlTags)
		if err != nil {
			return "⚠️Ошибка получения задач по тегам", err
		}
		defer respTags.Body.Close()
		if respTags.StatusCode == http.StatusOK {
			var tagIDsResp struct {
				Ids []int `json:"ids"`
			}
			if err := json.NewDecoder(respTags.Body).Decode(&tagIDsResp); err == nil {
				tagTaskIDs = tagIDsResp.Ids
			}
		}
	}

	// Получаем общее количество задач (task_get_all)
	totalTasks := 0
	urlAll := "http://theory-service:8081/taskgetall" //
	respAll, err := http.Get(urlAll)
	if err != nil {
		return "⚠️Ошибка получения общего количества задач", err
	}
	defer respAll.Body.Close()
	if respAll.StatusCode == http.StatusOK {
		var allCountResp struct {
			Count int `json:"count"`
		}
		if err := json.NewDecoder(respAll.Body).Decode(&allCountResp); err == nil {
			totalTasks = allCountResp.Count
		}
	}

	// Статистика по всем задачам
	completedTotal := 0
	completedByTag := 0
	completedSet := make(map[int]struct{})
	for _, id := range tagTaskIDs {
		if _, ok := completedTasks[fmt.Sprint(id)]; ok {
			completedByTag++
			completedSet[id] = struct{}{}
		}
	}
	// completedTotal - по всем задачам (если completedTasks содержит id, которых нет в tagTaskIDs)
	completedTotal = len(completedTasks)

	percentTotal := 0
	if totalTasks > 0 {
		percentTotal = completedTotal * 100 / totalTasks
	}

	// Статистика по тегам
	percentByTags := 0
	tagTasksCount := len(tagTaskIDs)
	if tagTasksCount > 0 {
		percentByTags = completedByTag * 100 / tagTasksCount
	}

	// Оставшиеся задачи по тегам
	left := make([]int, 0, tagTasksCount)
	for _, id := range tagTaskIDs {
		if _, ok := completedTasks[fmt.Sprint(id)]; !ok {
			left = append(left, id)
		}
	}

	// Формируем ответ
	stat := fmt.Sprintf(
		"📊 Ваша статистика\\:\n\nВсего просмотрено %d%% задач \\(%d из %d\\)\nПо тэгам \\[%s \\] %d%% \\(%d из %d\\)\nНомера оставшихся\\: \\[%v \\]",
		percentTotal, completedTotal, totalTasks,
		strings.Join(userTags, ","), percentByTags, completedByTag, tagTasksCount, left,
	)
	return stat, nil
}
