package mdutils

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type MarkdownProcessor interface {
	Escape(text string) string
	Restore(text string, entities []tgbotapi.MessageEntity, cmd string) string
}
