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
func (r *UserStateRepository) Patch(userID int64, patch models.UserState) error {
	setParts := []string{}
	args := []interface{}{}
	idx := 1

	//fmt.Printf("\n\n user-service - user_state_repository - patch.TheoryTags - %v\n", patch.TheoryTags)

	if patch.TheoryTags != nil {
		setParts = append(setParts, "theory_tags=$"+strconv.Itoa(idx))
		args = append(args, pq.Array(patch.TheoryTags))
		idx++
		//fmt.Println("--- не равно nil ---")
	} else {
		//fmt.Println("--- равно nil ---")
	}

	if len(setParts) == 0 {
		return nil
	}
	setParts = append(setParts, "updated_at=now()")
	setClause := strings.Join(setParts, ", ")
	args = append(args, userID)
	query := `UPDATE user_states SET ` + setClause + ` WHERE user_id=$` + strconv.Itoa(idx)

	//fmt.Printf("\n\nquery=%s, args=%v\n\n\n", query, args)
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

	row := r.DB.QueryRow(`
    SELECT user_id, last_theory_task_id, last_code_task_id,
           last_theory_answer, last_code_answer,
           theory_tags, code_tags, last_action, updated_at
    FROM user_states
    WHERE user_id=$1
`, userID)

	var state models.UserState
	err := row.Scan(
		&state.UserID,
		&state.LastTheoryTaskID,
		&state.LastCodeTaskID,
		&state.LastTheoryAnswer,
		&state.LastCodeAnswer,
		pq.Array(&state.TheoryTags),
		pq.Array(&state.CodeTags),
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
