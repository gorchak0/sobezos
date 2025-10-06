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

// NewUserStateRepository создаёт новый репозиторий для работы с user_states.
func NewUserStateRepository(db *sql.DB) *UserStateRepository {
	return &UserStateRepository{DB: db}
}

// Patch обновляет указанные поля user_state для пользователя по userID.
func (r *UserStateRepository) Patch(userID int64, patch models.UserState) error {
	setParts := []string{}
	args := []interface{}{}
	idx := 1

	if patch.LastTheoryTaskID != nil {
		setParts = append(setParts, "last_theory_task_id=$"+strconv.Itoa(idx))
		args = append(args, patch.LastTheoryTaskID)
		idx++
	}
	if patch.TheoryTags != nil {
		setParts = append(setParts, "theory_tags=$"+strconv.Itoa(idx))
		args = append(args, pq.Array(patch.TheoryTags))
		idx++
	}
	if patch.CompletedTheoryTasks != nil {
		setParts = append(setParts, "completed_theory_tasks=$"+strconv.Itoa(idx))
		args = append(args, pq.Array(patch.CompletedTheoryTasks))
		idx++
	}
	if patch.LastAction != nil {
		setParts = append(setParts, "last_action=$"+strconv.Itoa(idx))
		args = append(args, patch.LastAction)
		idx++
	}

	if len(setParts) == 0 {
		return nil
	}
	setParts = append(setParts, "updated_at=now()")
	setClause := strings.Join(setParts, ", ")
	args = append(args, userID)
	query := `UPDATE user_states SET ` + setClause + ` WHERE user_id=$` + strconv.Itoa(idx)

	_, err := r.DB.Exec(query, args...)
	if err != nil {
		log.Printf("[Patch] Error updating user_state for user_id=%d: %v, query=%s, args=%v", userID, err, query, args)
		return err
	}
	log.Printf("[Patch] Successfully updated user_state for user_id=%d, query=%s, args=%v", userID, query, args)
	return nil
}

// Get возвращает структуру UserState или ошибку, если запись не найдена или возникла ошибка.
// Get возвращает структуру UserState или ошибку, если запись не найдена или возникла ошибка.
func (r *UserStateRepository) Get(userID int64) (*models.UserState, error) {

	row := r.DB.QueryRow(`
    SELECT user_id, last_theory_task_id, theory_tags, 
           completed_theory_tasks, last_action, updated_at
    FROM user_states
    WHERE user_id=$1
`, userID)

	var state models.UserState
	err := row.Scan(
		&state.UserID,
		&state.LastTheoryTaskID,
		pq.Array(&state.TheoryTags),
		pq.Array(&state.CompletedTheoryTasks),
		&state.LastAction,
		&state.UpdatedAt,
	)

	if err != nil {
		log.Printf("[Get] Error getting user_state for user_id=%d: %v", userID, err)
		return nil, err
	}
	log.Printf("[Get] Successfully got user_state for user_id=%d: %+v", userID, state)
	return &state, nil
}

// AddState создает новую запись user_state для пользователя
func (r *UserStateRepository) AddState(state models.UserState) error {
	query := `INSERT INTO user_states (
		user_id, last_theory_task_id, theory_tags, completed_theory_tasks, last_action, updated_at
	) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.DB.Exec(query,
		state.UserID,
		state.LastTheoryTaskID,
		pq.Array(state.TheoryTags),
		pq.Array(state.CompletedTheoryTasks),
		state.LastAction,
		state.UpdatedAt,
	)
	if err != nil {
		log.Printf("[Add] Error inserting user_state for user_id=%d: %v", state.UserID, err)
		return err
	}
	log.Printf("[Add] Successfully inserted user_state for user_id=%d", state.UserID)
	return nil
}
