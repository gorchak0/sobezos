package bot

import (
	"fmt"
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
		//msg.ParseMode = "MarkdownV2" // или "MarkdownV2"
		bot.Send(msg)
	}
}

// HandleCommand обрабатывает команды бота и возвращает текст ответа
func (h *Handler) HandleCommand(cmd string, telegramID int, chatID int64, args string) string {
	userInfo, exists := h.service.UserCheck(telegramID) // user-service query param telegram_id
	if !exists {
		return fmt.Sprintf("⛔ Вы не авторизованы \nПерешлите это сообщение администратору чтобы он вас добавил: \n /useradd %d <username> <role>\nАдмин: @sorokinengineer", telegramID)
		//INSERT INTO users (telegram_id, username, role) VALUES (<telegramID>, '<username>', 'admin');
	}
	//UserStateEdit используется только под капотом

	// команды, доступные всем пользователям
	commands := map[string]func() (string, error){
		"taskget":   func() (string, error) { return h.service.TaskGet(telegramID) },         //theory-service
		"taskgetid": func() (string, error) { return h.service.TaskGetID(telegramID, args) }, //theory-service query param task_id
		"answerget": func() (string, error) { return h.service.AnswerGet(telegramID) },       //theory-service query param task_id
		"statsget":  func() (string, error) { return "Статистика пока не реализована", nil },
		"tagset":    func() (string, error) { return h.service.TagSet(telegramID, args) }, //user-service query param user_id
		"tagclear":  func() (string, error) { return h.service.TagClear(telegramID) },     // user-service query param user_id
		"tagget":    func() (string, error) { return h.service.TagGet() },                 //theory-service
		"help":      func() (string, error) { return help() },
		"start":     func() (string, error) { return help() },
		"taskadd":   func() (string, error) { return h.service.TaskAdd(telegramID, args) },  //theory-service
		"taskedit":  func() (string, error) { return h.service.TaskEdit(telegramID, args) }, //theory-service
	}

	// команды только для админов
	adminCommands := map[string]func() (string, error){
		"useradd": func() (string, error) { return h.service.UserAdd(telegramID, args) }, //user-service
		//useredit todo

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

func help() (string, error) {
	return `
Доступно всем пользователям:
/taskget — Получить случайную задачу (с учётом ваших тегов, если заданы)
/taskgetid <id> — Получить задачу по её номеру (id)
/answerget — Получить ответ на последнюю полученную задачу
/statsget — Показать статистику (пока не реализовано)
/tagset <тег1, тег2,...> — Добавить теги для фильтрации задач (например: /tagset динамика, графы)
/tagclear — Очистить все ваши теги
/tagget — Показать список всех доступных тегов с описаниями
/taskadd <json> — Добавить новую задачу 
/taskedit <json> — Редактировать задачу 

Только для администраторов:
/useradd <id> <username> — Добавить пользователя
`, nil
}
