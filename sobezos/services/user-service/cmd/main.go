package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sobezos/services/user-service/internal/handlers"
	"sobezos/services/user-service/internal/migrations"
	"sobezos/services/user-service/internal/repository"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
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

	// Миграция: выполнить SQL из файлов
	migrations.Migrate(db, "internal/migrations/001_create_users_table.sql")
	migrations.Migrate(db, "internal/migrations/002_create_user_states_table.sql")

	userRepo := repository.NewUserRepository(db)
	stateRepo := repository.NewUserStateRepository(db)
	logger := zap.NewExample() // используйте ваш logger
	handler := handlers.NewUserServiceHandler(userRepo, stateRepo, logger)

	//для всех пользователей
	http.HandleFunc("/usercheck", handler.CheckUser)
	http.HandleFunc("/userstateget", handler.GetState)
	http.HandleFunc("/userstateedit", handler.PatchState)

	//для админов
	http.HandleFunc("/useradd", handler.AddUser)

	log.Println("user-service listening on :8082")
	log.Fatal(http.ListenAndServe(":8082", nil))
}

//
