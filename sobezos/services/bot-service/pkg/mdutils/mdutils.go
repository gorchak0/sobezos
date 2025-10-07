package mdutils

import (
	"fmt"
	"sort"
	"strings"
	"unicode/utf8"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type MarkdownV2Processor struct{}

func NewMarkdownV2Processor() *MarkdownV2Processor {
	return &MarkdownV2Processor{}
}

// escapeMDV2 экранирует все специальные символы MarkdownV2
func (p *MarkdownV2Processor) EscapeMD(text string) string {
	replacer := strings.NewReplacer(
		"_", "\\_",
		"*", "\\*",
		"[", "\\[",
		"]", "\\]",
		"(", "\\(",
		")", "\\)",
		"~", "\\~",
		"`", "\\`",
		">", "\\>",
		"#", "\\#",
		"+", "\\+",
		"-", "\\-",
		"=", "\\=",
		"|", "\\|",
		"{", "\\{",
		"}", "\\}",
		".", "\\.",
		"!", "\\!",
	)
	return replacer.Replace(text)
}

// chunk описывает кусок текста: с разметкой или без
type chunk struct {
	start  int
	end    int
	entity *tgbotapi.MessageEntity
}

func (p *MarkdownV2Processor) AddMD(message string, entities []tgbotapi.MessageEntity) string {
	if len(entities) == 0 {
		return p.EscapeMD(message)
	}

	utf16map := utf16OffsetMap(message)
	var chunks []chunk

	// сортируем entities по Offset
	sort.SliceStable(entities, func(i, j int) bool {
		return entities[i].Offset < entities[j].Offset
	})

	// разбиваем текст на куски
	pos := 0
	for _, e := range entities {
		start := e.Offset
		end := e.Offset + e.Length
		if start >= len(utf16map) {
			continue
		}
		if end > len(utf16map) {
			end = len(utf16map)
		}

		// кусок до entity (без разметки)
		if pos < start {
			chunks = append(chunks, chunk{
				start:  pos,
				end:    start,
				entity: nil,
			})
		}

		// кусок с entity
		chunks = append(chunks, chunk{
			start:  start,
			end:    end,
			entity: &e,
		})

		pos = end
	}

	// оставшийся текст после последней entity
	if pos < len(utf16map) {
		chunks = append(chunks, chunk{
			start:  pos,
			end:    len(utf16map),
			entity: nil,
		})
	}

	// собираем результат
	var sb strings.Builder
	for _, c := range chunks {
		// пропускаем команду
		if c.entity != nil && c.entity.Type == "bot_command" {
			continue
		}

		/*

			startByte := utf16map[c.start]
			endByte := utf16map[c.end-1]
			if c.end < len(utf16map) {
				_, size := utf8.DecodeRuneInString(message[endByte:])
				endByte += size
			}

		*/

		startByte := utf16map[c.start]
		var endByte int
		if c.end >= len(utf16map) {
			endByte = len(message)
		} else {
			endByte = utf16map[c.end]
		}

		text := p.EscapeMD(message[startByte:endByte])

		fmt.Printf("[AddMD] chunk text before MD: %q\n", text)

		if c.entity != nil {
			switch c.entity.Type {
			case "bold":
				text = "*" + text + "*"
			case "italic":
				text = "_" + text + "_"
			case "underline":
				text = "__" + text + "__"
			case "strikethrough":
				text = "~" + text + "~"
			case "code":
				text = "`" + text + "`"
			case "pre":
				text = "```\n" + text + "\n```"
			}
		}

		sb.WriteString(text)
	}

	result := strings.TrimSpace(sb.String())

	return result
}

// utf16OffsetMap строит соответствие UTF-16 кодовая единица → UTF-8 байт
func utf16OffsetMap(s string) []int {
	var map16 []int
	i := 0
	for _, r := range s {
		if r <= 0xFFFF {
			map16 = append(map16, i)
			i += utf8.RuneLen(r)
		} else {
			// суррогатная пара — две UTF-16 единицы
			map16 = append(map16, i)
			map16 = append(map16, i)
			i += utf8.RuneLen(r)
		}
	}
	return map16
}
