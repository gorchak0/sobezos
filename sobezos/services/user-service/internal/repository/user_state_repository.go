package repository

import (
	"database/sql"
	"log"
	"strconv"
	"strings"

	"sobezos/services/user-service/pkg/models"

	"github.com/lib/pq"
)

// UserStateRepository отвечает за работу с таблицей user_states в базе данных.
type UserStateRepository struct {
	// DB — подключение к базе данных.
	DB *sql.DB
}

// Patch обновляет указанные поля user_state для пользователя по userID.
// patch — карта с именами и значениями обновляемых полей.
// Если patch пустой, ничего не происходит.
func (r *UserStateRepository) Patch(userID int64, patch map[string]interface{}) error {
	if len(patch) == 0 {
		return nil
	}
	// Формируем SET-часть запроса и аргументы
	setParts := []string{}
	args := []interface{}{}
	idx := 1
	for k, v := range patch {
		setParts = append(setParts, k+"=$"+strconv.Itoa(idx))
		args = append(args, v)
		idx++
	}
	// Обновляем поле updated_at
	setParts = append(setParts, "updated_at=now()")
	setClause := strings.Join(setParts, ", ")
	args = append(args, userID)
	query := `UPDATE user_states SET ` + setClause + ` WHERE user_id=$` + strconv.Itoa(idx)
	// Выполняем запрос
	_, err := r.DB.Exec(query, args...)
	if err != nil {
		log.Printf("[Patch] Error updating user_state for user_id=%d: %v, query=%s, args=%v", userID, err, query, args)
		return err
	}
	log.Printf("[Patch] Successfully updated user_state for user_id=%d, query=%s, args=%v", userID, query, args)
	return nil
}

// NewUserStateRepository создаёт новый репозиторий для работы с user_states.
func NewUserStateRepository(db *sql.DB) *UserStateRepository {
	return &UserStateRepository{DB: db}
}

// Get возвращает состояние пользователя по userID.
// Возвращает структуру UserState или ошибку, если запись не найдена или возникла ошибка.
func (r *UserStateRepository) Get(userID int64) (*models.UserState, error) {
	// Выполняем SELECT-запрос с COALESCE для nullable-полей
	row := r.DB.QueryRow(`SELECT user_id, COALESCE(last_theory_task_id, 0), COALESCE(last_code_task_id, 0), COALESCE(last_theory_answer, ''), COALESCE(last_code_answer, ''), COALESCE(theory_tags, ARRAY[]::text[]), COALESCE(code_tags, ARRAY[]::text[]), COALESCE(last_action, ''), updated_at FROM user_states WHERE user_id=$1`, userID)
	var state models.UserState
	// Сканируем результат в структуру
	err := row.Scan(&state.UserID, &state.LastTheoryTaskID, &state.LastCodeTaskID, &state.LastTheoryAnswer, &state.LastCodeAnswer, pq.Array(&state.TheoryTags), pq.Array(&state.CodeTags), &state.LastAction, &state.UpdatedAt)
	if err != nil {
		log.Printf("[Get] Error getting user_state for user_id=%d: %v", userID, err)
		return nil, err
	}
	log.Printf("[Get] Successfully got user_state for user_id=%d: %+v", userID, state)
	return &state, nil
}
