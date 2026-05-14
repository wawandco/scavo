-- 20260512160000 - add_is_admin_to_members
-- Adds an is_admin flag to members.

ALTER TABLE members ADD COLUMN is_admin INTEGER DEFAULT 0;
