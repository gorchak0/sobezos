package service

import (
	"encoding/json"
	"net/http"
)

func (s *Service) AnswerGet(telegramID int) (string, error) {

	// Получить id последнего теоретического вопроса
	state, err := s.UserStateGet(telegramID)
	//обработка ошибки
	if err != nil || state == nil {
		return "Нет информации о последней задаче", ErrServiceUnavailable
	}
	//достаем id
	taskID, ok := (*state)["last_theory_task_id"]
	if !ok {
		return "Нет информации о последней задаче", ErrServiceUnavailable
	}

	//получение текста вопроса
	url := "http://theory-service:8081/answer?task_id=" + taskID.(string)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", nil
	}
	var res struct {
		Answer string `json:"answer"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}

	//фиксируем
	err = s.UserStateEdit(telegramID, map[string]interface{}{
		"last_action":        "get_answer",
		"last_theory_answer": res.Answer,
	})
	if err != nil {
		return "", err
	}

	return res.Answer, nil
}
