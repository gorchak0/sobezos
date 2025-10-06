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
