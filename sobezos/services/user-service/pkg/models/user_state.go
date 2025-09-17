package models

import "time"

type UserState struct {
	UserID           int64      `json:"user_id"`
	LastTheoryTaskID *int       `json:"last_theory_task_id,omitempty"`
	LastCodeTaskID   *int       `json:"last_code_task_id,omitempty"`
	LastTheoryAnswer *string    `json:"last_theory_answer,omitempty"`
	LastCodeAnswer   *string    `json:"last_code_answer,omitempty"`
	TheoryTags       []string   `json:"theory_tags"`
	CodeTags         []string   `json:"code_tags"`
	LastAction       *string    `json:"last_action,omitempty"`
	UpdatedAt        *time.Time `json:"updated_at,omitempty"`
}
