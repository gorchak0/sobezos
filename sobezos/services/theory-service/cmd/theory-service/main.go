package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"go.uber.org/zap"

	"sobezos/services/theory-service/internal/handler"
	"sobezos/services/theory-service/internal/migrations"
	"sobezos/services/theory-service/internal/repository"
	"sobezos/services/theory-service/internal/service"
)

func main() {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dsn := "host=" + dbHost + " port=" + dbPort + " user=" + dbUser + " password=" + dbPassword + " dbname=" + dbName + " sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.Close()

	migrations.Migrate(db, "internal/migrations/001_create_tasks_table.sql")

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to initialize zap logger: %v", err)
	}
	defer logger.Sync()

	repo := repository.NewTaskRepository(db)
	service := service.NewTaskService(repo)
	handler := handler.NewTaskHandlerWithLogger(service, logger)

	http.HandleFunc("/taskget", handler.GetRandomTask)
	http.HandleFunc("/taskadd", handler.CreateTask)
	http.HandleFunc("/taskedit", handler.EditTask)
	http.HandleFunc("/answerget", handler.GetTaskAnswer)
	logger.Info("theory-service listening on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
