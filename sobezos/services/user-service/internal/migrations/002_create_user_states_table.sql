-- Migration for user_states table
CREATE TABLE IF NOT EXISTS user_states (
    user_id BIGINT PRIMARY KEY REFERENCES users(telegram_id) ON DELETE CASCADE,
    last_theory_task_id INT,
    theory_tags TEXT[],
    completed_theory_tasks TEXT[],
    last_action TEXT,
    updated_at TIMESTAMP DEFAULT now()
);
