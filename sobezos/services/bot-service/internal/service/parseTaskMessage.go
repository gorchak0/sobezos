package service

import (
	"fmt"
	"strings"
)

func parseTaskMessage(msg string, requiredFields []string) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	lines := strings.Split(msg, "\n")
	currentKey := ""

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Обработка нового блока (ключа)
		if strings.HasPrefix(line, "$") {
			currentKey = line[1:]
			result[currentKey] = initializeFieldValue(currentKey)
			continue
		}

		// Обработка значения для текущего ключа
		if currentKey == "" {
			return nil, fmt.Errorf("строка вне блока: %s", line)
		}

		result[currentKey] = appendFieldValue(result[currentKey], line, currentKey)
	}

	// Проверка обязательных полей
	if err := validateRequiredFields(result, requiredFields); err != nil {
		return nil, err
	}

	return result, nil
}

// initializeFieldValue создает начальное значение для поля в зависимости от его типа
func initializeFieldValue(field string) interface{} {
	if field == "tags" {
		return []string{}
	}
	return ""
}

// appendFieldValue добавляет значение к существующему полю с учетом его типа
func appendFieldValue(currentValue interface{}, line string, field string) interface{} {
	switch field {
	case "tags":
		return append(currentValue.([]string), line)
	default:
		if currentValue.(string) != "" {
			return currentValue.(string) + "\n" + line
		}
		return line
	}
}

// validateRequiredFields проверяет наличие всех обязательных полей
func validateRequiredFields(data map[string]interface{}, requiredFields []string) error {
	var missingFields []string

	for _, field := range requiredFields {
		if _, exists := data[field]; !exists {
			missingFields = append(missingFields, field)
		}
	}

	if len(missingFields) > 0 {
		return fmt.Errorf("не указаны следующие блоки: %s", strings.Join(missingFields, ", "))
	}

	return nil
}
