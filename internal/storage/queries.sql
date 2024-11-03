-- name: create-account
INSERT INTO accounts (name, created_at, updated_at)
VALUES ($1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
RETURNING id, created_at, updated_at;

-- name: get-account
SELECT id, name, created_at, updated_at
FROM accounts
WHERE id = $1;

-- name: update-account
UPDATE accounts
SET name = $1, updated_at = CURRENT_TIMESTAMP
WHERE id = $2
RETURNING updated_at;

-- name: delete-account
DELETE FROM accounts WHERE id = $1;

-- name: list-accounts
SELECT id, name, created_at, updated_at
FROM accounts
ORDER BY id
LIMIT $1 OFFSET $2;

-- name: count-accounts
SELECT COUNT(*) FROM accounts;

-- name: create-inbox
INSERT INTO inboxes (account_id, email, created_at, updated_at)
VALUES ($1, $2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
RETURNING id, created_at, updated_at;

-- name: get-inbox
SELECT id, account_id, email, created_at, updated_at
FROM inboxes
WHERE id = $1;

-- name: update-inbox
UPDATE inboxes
SET email = $1
WHERE id = $2;

-- name: delete-inbox
DELETE FROM inboxes WHERE id = $1;

-- name: list-inboxes-by-account
SELECT id, account_id, email, created_at, updated_at
FROM inboxes
WHERE account_id = $1
ORDER BY id
LIMIT $2 OFFSET $3;

-- name: count-inboxes-by-account
SELECT COUNT(*)
FROM inboxes
WHERE account_id = $1;

-- name: get-inbox-by-email
SELECT id, account_id, email, created_at, updated_at
FROM inboxes
WHERE email = $1;

-- name: create-rule
INSERT INTO rules (inbox_id, sender, receiver, subject, created_at, updated_at)
VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
RETURNING id, created_at, updated_at;

-- name: get-rule
SELECT id, inbox_id, sender, receiver, subject, created_at, updated_at
FROM rules
WHERE id = $1;

-- name: update-rule
UPDATE rules
SET sender = $1, receiver = $2, subject = $3
WHERE id = $4;

-- name: delete-rule
DELETE FROM rules WHERE id = $1;

-- name: list-rules-by-inbox
SELECT id, inbox_id, sender, receiver, subject, created_at, updated_at
FROM rules
WHERE inbox_id = $1
ORDER BY id
LIMIT $2 OFFSET $3;

-- name: count-rules-by-inbox
SELECT COUNT(*)
FROM rules
WHERE inbox_id = $1;

-- name: list-rules
SELECT id, inbox_id, sender, receiver, subject, created_at, updated_at
FROM rules
ORDER BY id
LIMIT $1 OFFSET $2;

-- name: count-rules
SELECT COUNT(*) FROM rules;

-- name: create-message
INSERT INTO messages (inbox_id, sender, receiver, subject, body, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
RETURNING id, created_at, updated_at;

-- name: get-message
SELECT id, inbox_id, sender, receiver, subject, body, created_at, updated_at
FROM messages
WHERE id = $1;

-- name: list-messages-by-inbox
SELECT id, inbox_id, sender, receiver, subject, body, created_at, updated_at
FROM messages
WHERE inbox_id = $1
ORDER BY id
LIMIT $2 OFFSET $3;

-- name: count-messages-by-inbox
SELECT COUNT(*)
FROM messages
WHERE inbox_id = $1;

-- name: create-user
INSERT INTO users (
    name, username, password, email, status, kind,
    password_login, created_at, updated_at
) VALUES ($1, $2, $3, $4, $5, $6, $7, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
RETURNING id, created_at, updated_at;

-- name: get-user
SELECT id, name, username, password, email, status, kind,
       password_login, loggedin_at, created_at, updated_at
FROM users
WHERE id = $1;

-- name: update-user
UPDATE users
SET name = $1, username = $2, password = $3, email = $4,
    status = $5, kind = $6, password_login = $7,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $8
RETURNING updated_at;

-- name: delete-user
DELETE FROM users WHERE id = $1;

-- name: get-user-by-username
SELECT id, name, username, password, email, status, kind,
       password_login, loggedin_at, created_at, updated_at
FROM users
WHERE username = $1;

-- name: initialize-tables
CREATE TABLE IF NOT EXISTS accounts (
    id SERIAL PRIMARY KEY,
    name TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS inboxes (
    id SERIAL PRIMARY KEY,
    account_id INTEGER REFERENCES accounts(id) ON DELETE CASCADE,
    email TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS rules (
    id SERIAL PRIMARY KEY,
    inbox_id INTEGER REFERENCES inboxes(id) ON DELETE CASCADE,
    sender TEXT,
    receiver TEXT,
    subject TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS messages (
    id SERIAL PRIMARY KEY,
    inbox_id INTEGER REFERENCES inboxes(id) ON DELETE CASCADE,
    sender TEXT,
    receiver TEXT,
    subject TEXT,
    body TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name TEXT,
    username TEXT UNIQUE,
    password TEXT,
    email TEXT,
    status TEXT,
    kind TEXT,
    password_login BOOLEAN DEFAULT false,
    loggedin_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
