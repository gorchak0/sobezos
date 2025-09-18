package main

import (
	"bufio"
	"encoding/json"
	"os"
	"strconv"
	"strings"
)

// Хардкодим тег для всех задач
const hardcodedTag = "мод_4_БД"

type Task struct {
	Tags     []string `json:"tags"`
	Question string   `json:"question"`
	Answer   string   `json:"answer"`
}

func main() {
	file, err := os.Open("мод.4 БД.md")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var tasks []Task
	var currentQuestion string
	var currentAnswer []string
	//var currentTag string

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		// Пропускаем пустые строки
		if line == "" {
			continue
		}

		// Игнорируем картинки (markdown-изображения)
		if strings.HasPrefix(line, "[image") || strings.HasPrefix(line, "!") {
			continue
		}

		// Заголовок 1 -> новая тема (тэг)
		if strings.HasPrefix(line, "# ") {
			// currentTag = strings.ToLower(strings.ReplaceAll(line[2:], " ", "_"))
			// currentTag = strings.Trim(currentTag, "*{}#")
			continue
		}

		// Заголовок 2 -> новый вопрос
		if strings.HasPrefix(line, "## ") {
			// Если есть предыдущий вопрос, сохраняем его
			if currentQuestion != "" {
				tasks = append(tasks, Task{
					Tags:     []string{hardcodedTag},
					Question: currentQuestion,
					Answer:   strings.ReplaceAll(strings.Join(currentAnswer, "\n"), "\n", "\n\n"),
				})
			}

			// Обрезаем хвост вида ** {#...
			q := line[3:]
			if idx := strings.Index(q, "** {#"); idx != -1 {
				q = q[:idx]
			}
			// Удаляем ведущие звёздочки и пробелы
			q = strings.TrimLeft(q, "* ")
			currentQuestion = strings.TrimSpace(q)
			currentAnswer = []string{}
			continue
		}

		// Всё остальное — это текст ответа
		currentAnswer = append(currentAnswer, line)
	}

	// Сохраняем последний вопрос
	if currentQuestion != "" {
		tasks = append(tasks, Task{
			Tags:     []string{hardcodedTag},
			Question: currentQuestion,
			Answer:   strings.ReplaceAll(strings.Join(currentAnswer, "\n"), "\n", "\n\n"),
		})
	}

	// Сохраняем задачи по 10 штук в отдельные файлы
	batchSize := 10
	total := len(tasks)
	for i := 0; i < total; i += batchSize {
		start := i + 1
		end := i + batchSize
		if end > total {
			end = total
		}
		fileName := "tasks/tasks_" + itoa(start) + "_" + itoa(end) + ".txt"
		outFile, err := os.Create(fileName)
		if err != nil {
			panic(err)
		}
		for j := i; j < end; j++ {
			data, _ := json.MarshalIndent(tasks[j], "", "  ")
			outFile.WriteString("/taskadd\n")
			outFile.WriteString(string(data) + "\n")
		}
		outFile.Close()
	}
}

// Преобразует int в строку
func itoa(i int) string {
	return strconv.Itoa(i)
}
