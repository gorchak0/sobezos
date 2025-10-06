package bot

import (
	"fmt"
	"sobezos/services/bot-service/internal/service"
	"sobezos/services/bot-service/pkg/mdutils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

type Handler struct {
	logger  *zap.Logger
	service service.ServiceInterface
	mdUtil  mdutils.MarkdownProcessor
}

func NewHandler(logger *zap.Logger, svc service.ServiceInterface) *Handler {
	return &Handler{
		logger: logger, service: svc,
		mdUtil: mdutils.NewMarkdownV2Processor(),
	}
}

func (h *Handler) HandleUpdate(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	if update.Message == nil {
		fmt.Println("Received non-message update")
		return
	}
	fmt.Println("Received message update")

	telegramID := int(update.Message.From.ID)
	fmt.Println("Get Telegram ID:", telegramID)
	chatID := update.Message.Chat.ID
	fmt.Println("Get Chat ID:", chatID)
	cmd := update.Message.Command()
	fmt.Println("Get Command:", cmd)

	args := update.Message.CommandArguments()
	fmt.Println("Get Arguments:", args)

	argsWithMarkdown := h.mdUtil.Restore(args, update.Message.Entities, cmd)
	fmt.Println("Get Arguments with Markdown:", argsWithMarkdown)

	h.logger.Info("HandleUpdate: received command",
		zap.String("command", cmd),
		zap.Int("telegram_id", telegramID),
		zap.Int64("chat_id", chatID),
		zap.String("args", argsWithMarkdown),
	)

	msgText := h.HandleCommand(cmd, telegramID, chatID, argsWithMarkdown)
	if msgText != "" {
		fmt.Println("Message from user is not empty")

		fmt.Printf("\n\nResponding to user %d in chat %d with message: %s\n", telegramID, chatID, msgText)

		msg := tgbotapi.NewMessage(chatID, msgText)
		msg.ParseMode = "MarkdownV2" // или "MarkdownV2"

		h.logger.Info("HandleUpdate: sending message",
			zap.Int64("chat_id", chatID),
			zap.String("text", msgText),
		)

		if _, err := bot.Send(msg); err != nil {
			h.logger.Error("Failed to send message", zap.Error(err))
		}
	}
}

// HandleCommand обрабатывает команды бота и возвращает текст ответа
func (h *Handler) HandleCommand(cmd string, telegramID int, chatID int64, args string) string {
	userInfo, exists := h.service.UserCheck(telegramID) // user-service query param telegram_id
	name := userInfo.Username
	if !exists {
		h.logger.Warn("User not authorized", zap.Int("telegram_id", telegramID), zap.String("username", name))
		return fmt.Sprintf("⛔ Вы не авторизованы \nПерешлите это сообщение администратору чтобы он вас добавил: \n `/useradd %d %s user`\nАдмин: @sorokinengineer", telegramID, name)
	}
	//UserStateEdit используется только под капотом

	// команды, доступные всем пользователям
	commands := map[string]func() (string, error){
		"taskget":   func() (string, error) { return h.service.TaskGet(telegramID) },         //theory-service
		"taskgetid": func() (string, error) { return h.service.TaskGetID(telegramID, args) }, //theory-service query param task_id
		"answerget": func() (string, error) { return h.service.AnswerGet(telegramID) },       //theory-service query param task_id
		"statsget":  func() (string, error) { return h.service.StatsGet(telegramID) },
		"tagset":    func() (string, error) { return h.service.TagSet(telegramID, args) }, //user-service query param user_id
		"tagclear":  func() (string, error) { return h.service.TagClear(telegramID) },     // user-service query param user_id
		"tagget":    func() (string, error) { return h.service.TagGet() },                 //theory-service
		"help":      func() (string, error) { return h.service.Help() },
		"start":     func() (string, error) { return h.service.Help() },
		"taskadd":   func() (string, error) { return h.service.TaskAdd(telegramID, args) },  //theory-service
		"taskedit":  func() (string, error) { return h.service.TaskEdit(telegramID, args) }, //theory-service
	}

	// команды только для админов
	adminCommands := map[string]func() (string, error){
		"useradd": func() (string, error) { return h.service.UserAdd(telegramID, args) }, //user-service
		//useredit todo

	}

	h.logger.Info("HandleCommand: processing command",
		zap.String("command", cmd),
		zap.Int("telegram_id", telegramID),
		zap.Int64("chat_id", chatID),
		zap.String("args", args),
		zap.String("role", userInfo.Role),
	)

	if fn, ok := commands[cmd]; ok {
		resp, err := fn()
		if err != nil {
			h.logger.Error("Service error in command", zap.String("command", cmd), zap.Error(err))
			return fmt.Sprintf("⚠️ Ошибка работы сервиса\n %s", h.mdUtil.Escape(err.Error()))
		}
		h.logger.Info("Command executed successfully", zap.String("command", cmd), zap.String("response", resp))
		return resp
	}

	if fn, ok := adminCommands[cmd]; ok {
		if userInfo.Role != "admin" {
			h.logger.Warn("Admin command attempted by non-admin", zap.String("command", cmd), zap.String("role", userInfo.Role))
			return "⛔ Только администратор может выполнять эту команду"
		}
		resp, err := fn()
		if err != nil {
			h.logger.Error("Service error in admin command", zap.String("command", cmd), zap.Error(err))
			return fmt.Sprintf("⚠️ Ошибка работы сервиса\n %s", h.mdUtil.Escape(err.Error()))
		}
		h.logger.Info("Admin command executed successfully", zap.String("command", cmd), zap.String("response", resp))
		return resp
	}

	h.logger.Warn("Unknown command received", zap.String("command", cmd), zap.Int("telegram_id", telegramID))
	return "❓ Неизвестная команда"
}
