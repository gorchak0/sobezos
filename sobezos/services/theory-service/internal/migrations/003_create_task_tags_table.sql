CREATE TABLE IF NOT EXISTS  task_tags (
    task_id INT NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    tag_id INT NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (task_id, tag_id)
);