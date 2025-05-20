DROP TABLE IF EXISTS users;
DROP INDEX IF EXISTS idx_users_username;

CREATE TABLE IF NOT EXISTS users
(
    id       INTEGER,
    username string UNIQUE,
    password string
);

CREATE INDEX IF NOT EXISTS idx_users_username ON users (username);
