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
	migrations.Migrate(db, "internal/migrations/002_create_tags_table.sql")
	migrations.Migrate(db, "internal/migrations/003_create_tasks_tags_table.sql")

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to initialize zap logger: %v", err)
	}
	defer logger.Sync()

	repo := repository.NewTaskRepository(db)
	service := service.NewTaskService(repo)
	handler := handler.NewTaskHandler(service, logger)

	http.HandleFunc("/taskget", handler.TaskGet)
	http.HandleFunc("/taskgetid", handler.TaskGetID)
	http.HandleFunc("/answerget", handler.AnswerGet)
	http.HandleFunc("/tagget", handler.TagGet)
	http.HandleFunc("/taskadd", handler.TaskAdd)
	http.HandleFunc("/taskedit", handler.TaskEdit)
	http.HandleFunc("/taskgetall", handler.TaskGetAll)
	http.HandleFunc("/taskgettags", handler.TaskGetTags)

	logger.Info("theory-service listening on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
