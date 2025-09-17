package service

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

// SendAddTask отправляет задачу в theory-service
func (s *Service) TaskAdd(jsonText string) (string, error) {

	//

	// Добавляем тег с username или telegram_id
	var reqBody map[string]interface{}
	if err := json.Unmarshal([]byte(args), &reqBody); err != nil {
		return "Некорректный JSON задачи"
	}
	var tag string
	if userInfo.Username != "" {
		tag = "@" + userInfo.Username
	} else {
		tag = "user:" + strconv.FormatInt(telegramID, 10)
	}
	tags, ok := reqBody["tags"].([]interface{})
	if !ok {
		tags = []interface{}{}
	}
	tags = append(tags, tag)
	reqBody["tags"] = tags
	newJson, err := json.Marshal(reqBody)
	if err != nil {
		return "Ошибка формирования задачи"
	}

	//
	//эндпоинт
	url := "http://theory-service:8081/createtask"
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(jsonText))
	if err != nil {
		return "Ошибка формирования запроса к theory-service"
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "Ошибка запроса к theory-service"
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusCreated {
		return "Задача успешно добавлена"
	}
	respMsg, _ := io.ReadAll(resp.Body)
	return "Ошибка: " + string(respMsg)
}
