package bot

import (
	"sobezos/services/bot-service/internal/service"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

type Handler struct {
	logger  *zap.Logger
	service *service.Service
}

func NewHandler(logger *zap.Logger) *Handler {
	return &Handler{logger: logger, service: service.NewService(logger)}
}

func (h *Handler) HandleUpdate(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	if update.Message == nil {
		return
	}

	telegramID := int(update.Message.From.ID)
	chatID := update.Message.Chat.ID
	cmd := update.Message.Command()
	args := update.Message.CommandArguments()

	msgText := h.HandleCommand(cmd, telegramID, chatID, args)
	if msgText != "" {
		msg := tgbotapi.NewMessage(chatID, msgText)
		bot.Send(msg)
	}
}

// HandleCommand обрабатывает команды бота и возвращает текст ответа
func (h *Handler) HandleCommand(cmd string, telegramID int, chatID int64, args string) string {
	userInfo, exists := h.service.UserCheck(telegramID)
	if !exists {
		return "⛔ Вы не авторизованы"
	}
	//UserStateEdit используется только под капотом

	// команды, доступные всем пользователям
	commands := map[string]func() (string, error){
		"taskget": func() (string, error) { return h.service.TaskGet(telegramID) },
		//taskgetid todo
		"answerget": func() (string, error) { return h.service.AnswerGet(telegramID) },
		"statsget":  func() (string, error) { return "Статистика пока не реализована", nil },
		"tagset":    func() (string, error) { return h.service.TagSet(telegramID, args) },
		//tagsset без аргументов - tagsclear todo
		//tagget todo
	}

	// команды только для админов
	adminCommands := map[string]func() (string, error){
		"useradd": func() (string, error) { return h.service.UserAdd(telegramID, args) },
		//useredit todo
		"taskadd":  func() (string, error) { return h.service.TaskAdd(args) },
		"taskedit": func() (string, error) { return h.service.TaskEdit(args) },
	}

	if fn, ok := commands[cmd]; ok {
		resp, err := fn()
		if err != nil {
			return "⚠️ Ошибка работы сервиса"
		}
		return resp
	}

	if fn, ok := adminCommands[cmd]; ok {
		if userInfo.Role != "admin" {
			return "⛔ Только администратор может выполнять эту команду"
		}
		resp, err := fn()
		if err != nil {
			return "⚠️ Ошибка работы сервиса"
		}
		return resp
	}

	return "❓ Неизвестная команда"
}
