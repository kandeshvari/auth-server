
-- +migrate Up
CREATE TABLE auth (
	id INTEGER AUTO_INCREMENT,
	login VARCHAR(256) NOT NULL UNIQUE,
	password_hash VARCHAR(256) NOT NULL,
	PRIMARY KEY (id)

) CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE UNIQUE INDEX idx_users_login ON auth(login);

CREATE TABLE sessions (
	token VARCHAR(64),
	expire INTEGER NOT NULL,
	user_id INTEGER NOT NULL,
	ip VARCHAR(16) NOT NULL,
	agent VARCHAR(256) NOT NULL,
	os VARCHAR(64) NOT NULL,
	PRIMARY KEY (token),
	FOREIGN KEY (user_id) REFERENCES auth(id)
) CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE UNIQUE INDEX idx_sessions_token ON sessions(token);

-- +migrate Down
DROP TABLE sessions;
DROP TABLE auth;
