BEGIN;

DROP INDEX IF EXISTS idx_users_line_user_id;

DROP TABLE IF EXISTS users;

DROP TYPE IF EXISTS UserStatus;

END;