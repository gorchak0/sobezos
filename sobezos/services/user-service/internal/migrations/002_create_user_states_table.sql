-- Migration for user_states table
CREATE TABLE IF NOT EXISTS user_states (
    user_id BIGINT PRIMARY KEY REFERENCES users(telegram_id) ON DELETE CASCADE,
    last_theory_task_id INT,
    last_code_task_id INT,
    last_theory_answer TEXT,
    last_code_answer TEXT,
    theory_tags TEXT[],
    code_tags TEXT[],
    last_action TEXT,
    updated_at TIMESTAMP DEFAULT now()
);
