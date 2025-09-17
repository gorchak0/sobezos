package repository

import (
	"database/sql"
	"sobezos/services/theory-service/pkg/models"

	"github.com/lib/pq"
)

// TaskRepository handles DB operations for tasks
type TaskRepository struct {
	DB *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{DB: db}
}

// GetRandomTaskRow возвращает случайную задачу (без тегов)
func (r *TaskRepository) GetRandomTaskRow(tags []string) (*models.Task, error) {
	task := &models.Task{}
	if len(tags) == 0 {
		row := r.DB.QueryRow(`SELECT t.id, t.question, t.answer, t.created_at FROM tasks t ORDER BY RANDOM() LIMIT 1`)
		if err := row.Scan(&task.ID, &task.Question, &task.Answer, &task.CreatedAt); err != nil {
			return nil, err
		}
	} else {
		query := `SELECT t.id, t.question, t.answer, t.created_at FROM tasks t JOIN task_tags tt ON t.id = tt.task_id JOIN tags tg ON tt.tag_id = tg.id WHERE tg.name = ANY($1) GROUP BY t.id HAVING COUNT(DISTINCT tg.name) = $2 ORDER BY RANDOM() LIMIT 1;`
		row := r.DB.QueryRow(query, pq.Array(tags), len(tags))
		if err := row.Scan(&task.ID, &task.Question, &task.Answer, &task.CreatedAt); err != nil {
			return nil, err
		}
	}
	return task, nil
}

func (r *TaskRepository) CreateTask(task *models.Task) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Сначала создаём задачу
	var taskID int
	err = tx.QueryRow(
		`INSERT INTO tasks (question, answer, created_at) VALUES ($1, $2, $3) RETURNING id`,
		task.Question, task.Answer, task.CreatedAt,
	).Scan(&taskID)
	if err != nil {
		return err
	}

	// Добавляем связи с тегами
	for _, tagName := range task.Tags {
		var tagID int
		// Создаём тег если не существует
		err := tx.QueryRow(
			`INSERT INTO tags (name) VALUES ($1) ON CONFLICT (name) DO UPDATE SET name=EXCLUDED.name RETURNING id`,
			tagName,
		).Scan(&tagID)
		if err != nil {
			return err
		}

		// Создаём запись в task_tags
		_, err = tx.Exec(`INSERT INTO task_tags (task_id, tag_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`, taskID, tagID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *TaskRepository) UpdateTask(task *models.Task) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Обновляем вопрос и ответ
	_, err = tx.Exec(
		`UPDATE tasks SET question=$1, answer=$2 WHERE id=$3`,
		task.Question, task.Answer, task.ID,
	)
	if err != nil {
		return err
	}

	// Удаляем старые связи
	_, err = tx.Exec(`DELETE FROM task_tags WHERE task_id=$1`, task.ID)
	if err != nil {
		return err
	}

	// Добавляем новые связи с тегами
	for _, tagName := range task.Tags {
		var tagID int
		// Создаём тег если не существует
		err := tx.QueryRow(
			`INSERT INTO tags (name) VALUES ($1) ON CONFLICT (name) DO UPDATE SET name=EXCLUDED.name RETURNING id`,
			tagName,
		).Scan(&tagID)
		if err != nil {
			return err
		}

		_, err = tx.Exec(`INSERT INTO task_tags (task_id, tag_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`, task.ID, tagID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *TaskRepository) GetTaskByID(id int) (*models.Task, error) {
	task := &models.Task{}
	err := r.DB.QueryRow(`SELECT id, question, answer, created_at FROM tasks WHERE id=$1`, id).
		Scan(&task.ID, &task.Question, &task.Answer, &task.CreatedAt)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (r *TaskRepository) GetAllTags() ([]*models.Tag, error) {
	rows, err := r.DB.Query(`SELECT id, name, description FROM tags ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []*models.Tag
	for rows.Next() {
		tag := &models.Tag{}
		if err := rows.Scan(&tag.ID, &tag.Name, &tag.Description); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

func (r *TaskRepository) GetTagsByTaskID(taskID int) ([]*models.Tag, error) {
	rows, err := r.DB.Query(`
        SELECT tg.id, tg.name, tg.description
        FROM tags tg
        JOIN task_tags tt ON tg.id = tt.tag_id
        WHERE tt.task_id = $1
        ORDER BY tg.name
    `, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []*models.Tag
	for rows.Next() {
		tag := &models.Tag{}
		if err := rows.Scan(&tag.ID, &tag.Name, &tag.Description); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, nil
}
