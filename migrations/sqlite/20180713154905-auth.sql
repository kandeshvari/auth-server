-- +migrate Up
CREATE TABLE auth (
  id            INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
  login         VARCHAR UNIQUE NOT NULL,
  password_hash VARCHAR        NOT NULL
);
CREATE UNIQUE INDEX idx_users_login
  ON auth (login);

CREATE TABLE sessions (
  token   VARCHAR PRIMARY KEY,
  expire  INTEGER NOT NULL,
  user_id INTEGER NOT NULL,
  ip      VARCHAR NOT NULL,
  agent   VARCHAR NOT NULL,
  os      VARCHAR NOT NULL,
  FOREIGN KEY (user_id) REFERENCES auth (id)
);
CREATE UNIQUE INDEX idx_sessions_token
  ON sessions (token);

-- +migrate Down
DROP TABLE sessions;
DROP TABLE auth;
