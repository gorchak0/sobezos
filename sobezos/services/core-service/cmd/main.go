package main

import (
	"log"
	"net/http"

	_ "sobezos/services/core-service/docs"
	"sobezos/services/core-service/internal/handler"
	"sobezos/services/core-service/internal/service"

	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

// @title           SOBEZOS API
// @version         1.0
// @description     Это пример API с использованием Swaggo
// @host            localhost:8083
// @BasePath        /
func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	serv := service.NewService(logger)
	h := handler.NewHandler(serv)

	mux := http.NewServeMux()
	mux.HandleFunc("/answerget", h.AnswerGet)
	mux.HandleFunc("/usercheck", h.UserCheck)
	mux.HandleFunc("/useradd", h.UserAdd)
	mux.HandleFunc("/taskadd", h.TaskAdd)
	mux.HandleFunc("/taskedit", h.TaskEdit)
	mux.HandleFunc("/taskget", h.TaskGet)
	mux.HandleFunc("/taskgetid", h.TaskGetID)
	mux.HandleFunc("/tagget", h.TagGet)
	mux.HandleFunc("/tagset", h.TagSet)
	mux.HandleFunc("/tagclear", h.TagClear)
	mux.HandleFunc("/statsget", h.StatsGet)
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	logger.Info("core-service listening on :8083")
	log.Fatal(http.ListenAndServe(":8083", withCORS(mux)))
}

//swag init -g main.go -o ./docs --parseInternal --parseDependency --dir ./cmd,./internal

// middleware для CORS
func withCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// разрешаем все домены
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// если браузер проверяет OPTIONS, просто возвращаем 204
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// передаём дальше
		h.ServeHTTP(w, r)
	})
}
