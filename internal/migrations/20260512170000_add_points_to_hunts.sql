-- 20260512170000 - add_points_to_hunts
-- Adds a points field to hunts with a default of 10.

ALTER TABLE hunts ADD COLUMN points INTEGER NOT NULL DEFAULT 10;
