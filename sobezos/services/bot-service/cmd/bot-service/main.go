package main

import (
	"os"
	"sobezos/services/bot-service/internal/config"
	handler "sobezos/services/bot-service/internal/handler"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	config.LoadEnv()
	token := os.Getenv("TELEGRAM_TOKEN")
	if token == "" {
		logger.Fatal("TELEGRAM_TOKEN not set")
	}

	botAPI, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		logger.Fatal("Failed to create bot", zap.Error(err))
	}
	logger.Info("bot-service started", zap.String("username", botAPI.Self.UserName))

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := botAPI.GetUpdatesChan(u)
	handler := handler.NewHandler(logger)

	for update := range updates {
		handler.HandleUpdate(update, botAPI)
	}
}
