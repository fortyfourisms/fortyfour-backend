ALTER TABLE users
    DROP INDEX idx_users_status,
    DROP COLUMN login_attempts,
    DROP COLUMN password_changed_at,
    DROP COLUMN status;
