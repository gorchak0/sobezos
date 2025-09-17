package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sobezos/services/bot-service/internal/models"
)

func (s *Service) AnswerGet(telegramID int) (string, error) {

	// Получить id последнего теоретического вопроса
	state, err := s.UserStateGet(telegramID)
	if err != nil || state == nil {
		return "Нет информации о последней задаче", ErrServiceUnavailable
	}
	taskID := state.LastTheoryTaskID
	if taskID == 0 {
		return "Нет информации о последней задаче", ErrServiceUnavailable
	}

	//получение текста вопроса
	url := fmt.Sprintf("http://theory-service:8081/answerget?task_id=%d", taskID)
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
	err = s.UserStateEdit(telegramID, models.UserState{
		LastAction:       "get_answer",
		LastTheoryTaskID: taskID,
	})
	if err != nil {
		return "", err
	}

	return res.Answer, nil
}
