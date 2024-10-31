# Mercury

A simple email server that allows you to create inboxes and rules to filter emails, written in Go.

## Table of Contents

- [Mercury](#mercury)
  - [Table of Contents](#table-of-contents)
  - [Features](#features)
  - [Configuration](#configuration)
    - [Using Config File](#using-config-file)
    - [Using Environment Variables](#using-environment-variables)
  - [Running the Server](#running-the-server)
    - [Standard Run](#standard-run)
    - [With Custom Config](#with-custom-config)
  - [API Examples](#api-examples)
    - [Account Management](#account-management)
    - [Inbox Management](#inbox-management)
    - [Rule Management](#rule-management)
  - [Testing Email Reception](#testing-email-reception)
  - [Development](#development)
    - [Project Structure](#project-structure)
    - [Adding New Features](#adding-new-features)
  - [Architecture](#architecture)

## Features

- HTTP API for managing accounts, inboxes, and rules
- SMTP server for receiving emails
- Rule-based email filtering
- Configurable via YAML and environment variables
- Structured logging with multiple levels
- SQLite database storage

## Configuration

The application can be configured using:
- YAML configuration file (default: `config/default.yaml`)
- Environment variables
- Command line flags

### Using Config File

```yaml
server:
  http:
    port: ":8080"
  smtp:
    port: ":1025"
    domain: "localhost"
    username: ""
    password: ""
database:
  driver: "sqlite3"
  url: "./email.db"
logging:
  level: "info"  # Available: debug, info, warn, error, fatal
  format: "json"
```

### Using Environment Variables

Environment variables override config file settings:

```shell
MERCURY_SERVER_HTTP_PORT=":9090"
MERCURY_SERVER_SMTP_PORT=":2025"
MERCURY_LOGGING_LEVEL="debug"
```

## Running the Server

### Standard Run

```shell
go run main.go
```

### With Custom Config

```shell
go run main.go -config=/path/to/config.yaml
```

## API Examples

### Account Management

Create an Account:
```shell
curl -X POST http://localhost:8080/accounts \
  -H "Content-Type: application/json" \
  -d '{"name": "Test Account"}'
```

List Accounts:
```shell
curl http://localhost:8080/accounts
```

### Inbox Management

Create an Inbox:
```shell
curl -X POST http://localhost:8080/accounts/1/inboxes \
  -H "Content-Type: application/json" \
  -d '{"email": "inbox@example.com"}'
```

List Inboxes for Account:
```shell
curl http://localhost:8080/accounts/1/inboxes
```

### Rule Management

Create a Rule:
```shell
curl -X POST http://localhost:8080/accounts/1/inboxes/1/rules \
  -H "Content-Type: application/json" \
  -d '{
    "sender": "sender@example.com",
    "receiver": "inbox@example.com",
    "subject": "Test Subject"
  }'
```

List Rules for Inbox:
```shell
curl http://localhost:8080/accounts/1/inboxes/1/rules
```

## Testing Email Reception

Using SWAKS:
```shell
swaks --to inbox@example.com \
      --from sender@example.com \
      --server localhost:1025 \
      --header "Subject: Test Subject" \
      --body "This is a test email."
```

Using Telnet:
```shell
telnet localhost 1025
HELO localhost
MAIL FROM:<sender@example.com>
RCPT TO:<inbox@example.com>
DATA
Subject: Test Subject

This is a test email.
.
QUIT
```

## Development

### Project Structure
```
.
├── config/
│   └── default.yaml     # Default configuration
├── internal/
│   ├── api/            # HTTP API handlers
│   ├── config/         # Configuration management
│   ├── core/           # Core business logic
│   ├── logger/         # Logging package
│   ├── models/         # Data models
│   ├── smtpserver/     # SMTP server implementation
│   └── storage/        # Database operations
└── main.go             # Application entry point
```

### Adding New Features

1. Define models in `internal/models`
2. Implement storage operations in `internal/storage`
3. Add business logic in `internal/core`
4. Create API endpoints in `internal/api`

## Architecture

The application follows a layered architecture:
- API Layer: HTTP handlers and SMTP server
- Service Layer: Business logic in core package
- Storage Layer: Database operations

Database schema is inspired by [Archiveopteryx](https://archiveopteryx.org/db/)
