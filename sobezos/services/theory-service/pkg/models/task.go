package models

import "time"

// Task represents a theory task
type Task struct {
	ID        int       `json:"id"`
	Question  string    `json:"question"`
	Answer    string    `json:"answer"`
	Tags      []string  `json:"tags,omitempty"` //только для отдачи клиенту
	CreatedAt time.Time `json:"created_at"`
}
