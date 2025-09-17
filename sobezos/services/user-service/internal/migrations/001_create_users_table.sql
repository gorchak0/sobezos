CREATE TABLE IF NOT EXISTS users (
  id SERIAL PRIMARY KEY,
  telegram_id BIGINT NOT NULL UNIQUE,
  username TEXT,
  role TEXT NOT NULL CHECK (role IN ('user','admin')),
  created_at TIMESTAMP DEFAULT now()
);
