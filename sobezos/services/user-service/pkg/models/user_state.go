package models

import "time"

type UserState struct {
	UserID           int64     `json:"user_id"`
	LastTheoryTaskID int       `json:"last_theory_task_id"`
	LastCodeTaskID   int       `json:"last_code_task_id"`
	LastTheoryAnswer string    `json:"last_theory_answer"`
	LastCodeAnswer   string    `json:"last_code_answer"`
	TheoryTags       []string  `json:"theory_tags"`
	CodeTags         []string  `json:"code_tags"`
	LastAction       string    `json:"last_action"`
	UpdatedAt        time.Time `json:"updated_at"`
}
