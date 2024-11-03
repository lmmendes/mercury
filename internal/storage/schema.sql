-- User related types
DROP TYPE IF EXISTS user_status CASCADE;
CREATE TYPE user_status AS ENUM ('enabled', 'disabled');

DROP TYPE IF EXISTS user_kind CASCADE;
CREATE TYPE user_kind AS ENUM ('admin', 'user');

-- Create tables
DROP TABLE IF EXISTS accounts CASCADE;
CREATE TABLE accounts (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_accounts_name ON accounts(name);
CREATE INDEX idx_accounts_created_at ON accounts(created_at);

DROP TABLE IF EXISTS inboxes CASCADE;
CREATE TABLE inboxes (
    id SERIAL PRIMARY KEY,
    account_id INTEGER REFERENCES accounts(id) ON DELETE CASCADE,
    email TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_inboxes_account_id ON inboxes(account_id);
CREATE UNIQUE INDEX idx_inboxes_email ON inboxes(LOWER(email));
CREATE INDEX idx_inboxes_created_at ON inboxes(created_at);

DROP TABLE IF EXISTS rules CASCADE;
CREATE TABLE rules (
    id SERIAL PRIMARY KEY,
    inbox_id INTEGER REFERENCES inboxes(id) ON DELETE CASCADE,
    sender TEXT,
    receiver TEXT,
    subject TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_rules_inbox_id ON rules(inbox_id);
CREATE INDEX idx_rules_sender ON rules(sender);
CREATE INDEX idx_rules_receiver ON rules(receiver);

DROP TABLE IF EXISTS messages CASCADE;
CREATE TABLE messages (
    id SERIAL PRIMARY KEY,
    inbox_id INTEGER REFERENCES inboxes(id) ON DELETE CASCADE,
    sender TEXT NOT NULL,
    receiver TEXT NOT NULL,
    subject TEXT,
    body TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_messages_inbox_id ON messages(inbox_id);
CREATE INDEX idx_messages_sender ON messages(sender);
CREATE INDEX idx_messages_receiver ON messages(receiver);
CREATE INDEX idx_messages_created_at ON messages(created_at);

DROP TABLE IF EXISTS users CASCADE;
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    username TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    status user_status NOT NULL DEFAULT 'enabled',
    kind user_kind NOT NULL DEFAULT 'user',
    password_login BOOLEAN DEFAULT false,
    loggedin_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX idx_users_username ON users(LOWER(username));
CREATE UNIQUE INDEX idx_users_email ON users(LOWER(email));
CREATE INDEX idx_users_status ON users(status);
CREATE INDEX idx_users_created_at ON users(created_at);
