package service

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

// SendEditTask отправляет PUT-запрос на theory-service для редактирования задачи
func (s *Service) TaskEdit(jsonText string) (string, error) {
	//эндпоинт
	url := "http://theory-service:8081/edittask"

	// Проверяем наличие id в полученном json
	var reqBody map[string]interface{}
	//если не получается размаршалить
	if err := json.Unmarshal([]byte(jsonText), &reqBody); err != nil {
		return "Некорректный JSON"
	}
	//если нет поля id
	if _, ok := reqBody["id"]; !ok {
		return "Для редактирования задачи необходимо указать id"
	}

	//формируем и выполняем запрос к theory-service
	req, err := http.NewRequest("PUT", url, bytes.NewBufferString(jsonText))
	if err != nil {
		return "Ошибка формирования запроса к theory-service"
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)

	//обработка ошибок
	if err != nil {
		return "Ошибка запроса к theory-service"
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		return "Задача успешно обновлена"
	}
	respMsg, _ := io.ReadAll(resp.Body)
	return "Ошибка: " + string(respMsg)
}
