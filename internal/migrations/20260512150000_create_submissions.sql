-- 20260512150000 - create_submissions migration
-- Creates the submissions table for hunt entries.

CREATE TABLE IF NOT EXISTS submissions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    hunt_id INTEGER NOT NULL,
    member_id INTEGER NOT NULL,
    team_id INTEGER NOT NULL,
    text TEXT NOT NULL,
    image_path TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
