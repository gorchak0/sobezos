package models

import "time"

type UserState struct {
	UserID               int64      `json:"user_id"`
	LastTheoryTaskID     *int       `json:"last_theory_task_id,omitempty"`
	TheoryTags           []string   `json:"theory_tags"`
	CompletedTheoryTasks []string   `json:"completed_theory_tasks"`
	LastAction           *string    `json:"last_action,omitempty"`
	UpdatedAt            *time.Time `json:"updated_at,omitempty"`
}

// User описывает пользователя из таблицы users
type User struct {
	ID         int64  `json:"id"`
	TelegramID int64  `json:"telegram_id"`
	Username   string `json:"username"`
	Role       string `json:"role"`
	CreatedAt  string `json:"created_at"`
}

// --- Request models for handlers ---

type UserAddRequest struct {
	Args string `json:"args"`
}

type TaskAddRequest struct {
	Question string   `json:"question"`
	Answer   string   `json:"answer"`
	Tags     []string `json:"tags"`
}

type TaskEditRequest struct {
	ID       string   `json:"id"`
	Question string   `json:"question"`
	Answer   string   `json:"answer"`
	Tags     []string `json:"tags"`
}

type TagSetRequest struct {
	Args string `json:"args"`
}

type TagClearRequest struct {
	TelegramID int `json:"telegram_id"`
}

// CommonSuccessResponse — универсальный ответ с результатом (например, для успешных операций)
type CommonSuccessResponse struct {
	Result string `json:"result"`
}

// CommonErrorResponse — универсальный ответ с ошибкой
type CommonErrorResponse struct {
	Error string `json:"error"`
}

// UserCheckResponse — ответ для проверки пользователя
type UserCheckResponse struct {
	Role     string `json:"role"`
	Username string `json:"username"`
}
