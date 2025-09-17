CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY,
    question TEXT,
    answer TEXT,
    created_at TIMESTAMP DEFAULT now()
);
