package mdutils

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type MarkdownV2Processor struct{}

func NewMarkdownV2Processor() *MarkdownV2Processor {
	fmt.Println("NewMarkdownV2Processor called")
	return &MarkdownV2Processor{}
}

func (m *MarkdownV2Processor) Escape(text string) string {
	fmt.Printf("escapeMarkdownV2 called with text: %q\n", text)
	replacer := []struct{ old, new string }{
		{"_", "\\_"},
		{"*", "\\*"},
		{"[", "\\["},
		{"]", "\\]"},
		{"(", "\\("},
		{")", "\\)"},
		{"~", "\\~"},
		{"`", "\\`"},
		{">", "\\>"},
		{"#", "\\#"},
		{"+", "\\+"},
		{"-", "\\-"},
		{"=", "\\="},
		{"|", "\\|"},
		{"{", "\\{"},
		{"}", "\\}"},
		{".", "\\."},
		{"!", "\\!"},
		{"u003c", ""},
	}
	res := text
	for _, r := range replacer {
		res = strings.ReplaceAll(res, r.old, r.new)
	}
	return res
}

func (m *MarkdownV2Processor) Restore(text string, entities []tgbotapi.MessageEntity, cmd string) string {

	shift := len(cmd) + 2

	fmt.Printf("Restore called with text: %q, entities: %+v, shift: %d\n", text, entities, shift)
	if len(entities) == 0 {
		return m.Escape(text)
	}

	// Сначала экранируем текст
	escaped := m.Escape(text)

	// Создаем карту смещений для каждого символа
	originalRunes := []rune(text)
	escapedRunes := []rune(escaped)

	// Строим маппинг позиций из оригинального в экранированный текст
	positionMap := make([]int, len(originalRunes)+1)
	origIdx, escIdx := 0, 0

	for origIdx < len(originalRunes) && escIdx < len(escapedRunes) {
		if origIdx < len(originalRunes) && escIdx < len(escapedRunes) &&
			originalRunes[origIdx] == escapedRunes[escIdx] {
			positionMap[origIdx] = escIdx
			origIdx++
			escIdx++
		} else {
			// Пропускаем экранирующие символы
			escIdx++
		}
	}
	positionMap[len(originalRunes)] = len(escapedRunes)

	// Применяем entities с учетом новых позиций и сдвига
	result := make([]rune, len(escapedRunes))
	copy(result, escapedRunes)

	// Применяем entities в обратном порядке
	for i := len(entities) - 1; i >= 0; i-- {
		ent := entities[i]
		start, end := markdownSymbols(ent.Type)
		if start == "" && end == "" {
			continue
		}

		// Применяем сдвиг к позициям entity
		shiftedOffset := ent.Offset - shift
		if shiftedOffset < 0 {
			shiftedOffset = 0
		}
		shiftedEnd := shiftedOffset + ent.Length
		if shiftedEnd > len(originalRunes) {
			shiftedEnd = len(originalRunes)
		}

		// Преобразуем позиции с учетом экранирования
		newOffset := positionMap[shiftedOffset]
		newEnd := positionMap[shiftedEnd]

		if newOffset < 0 || newEnd > len(result) || newOffset >= newEnd {
			fmt.Printf("Skipping entity: offset=%d, end=%d, result_len=%d\n", newOffset, newEnd, len(result))
			continue
		}

		// Вставляем закрывающий тег
		if newEnd <= len(result) {
			result = append(result[:newEnd], append([]rune(end), result[newEnd:]...)...)
		}

		// Вставляем открывающий тег
		if newOffset <= len(result) {
			result = append(result[:newOffset], append([]rune(start), result[newOffset:]...)...)
		}

		fmt.Printf("Applied entity: type=%s, original_offset=%d, shifted_offset=%d, new_offset=%d, new_end=%d\n",
			ent.Type, ent.Offset, shiftedOffset, newOffset, newEnd)
	}

	return string(result)
}

func markdownSymbols(entType string) (start, end string) {
	fmt.Printf("markdownSymbols called with entType: %q\n", entType)
	switch entType {
	case "bold":
		return "*", "*"
	case "italic":
		return "_", "_"
	case "code":
		return "`", "`"
	case "pre":
		return "```", "```"
	case "underline":
		return "__", "__"
	case "strikethrough":
		return "~", "~"
	default:
		return "", ""
	}
}
