# Inbox451

A simple email server that allows you to create inboxes and rules to filter emails, written in Go.

## Table of Contents

- [Mercury](#mercury)
  - [Table of Contents](#table-of-contents)
  - [Features](#features)
  - [Prerequisites](#prerequisites)
  - [Development Setup](#development-setup)
  - [Configuration](#configuration)
    - [Using Config File](#using-config-file)
    - [Using Environment Variables](#using-environment-variables)
  - [Running the Server](#running-the-server)
  - [API Examples](#api-examples)
    - [Account Management](#account-management)
    - [Inbox Management](#inbox-management)
    - [Rule Management](#rule-management)
  - [Testing Email Reception](#testing-email-reception)
  - [Development](#development)
    - [Project Structure](#project-structure)
    - [API Testing with Bruno](#api-testing-with-bruno)
  - [Architecture](#architecture)

## Features

- HTTP API for managing accounts, inboxes, and rules
- SMTP server for receiving emails
- Rule-based email filtering
- Configurable via YAML and environment variables

## Prerequisites

- Go 1.22 or later
- Docker and Docker Compose
- Make (optional, but recommended)

## Development Setup

1. Clone the repository
2. Start the PostgreSQL database and run the application:
```bash
make dev
```

Additional commands:
```bash
# Start the database
make db-up

# Stop the database
make db-down

# Clean database (remove volume)
make db-clean

# Reset database (down, clean, up)
make db-reset

# Run tests
make test
```

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
    hostname: "localhost"
    username: ""
    password: ""
  imap:
    port: ":1143"
    hostname: "localhost"
database:
  url: "postgres://inbox:inbox@localhost:5432/inbox451?sslmode=disable"
  max_open_conns: 25
  max_idle_conns: 5
  conn_max_lifetime: 5m
logging:
  level: "info"  # Available: debug, info, warn, error, fatal
  format: "json"
```

### Using Environment Variables

Environment variables override config file settings:

```shell
INBOX451_SERVER_HTTP_PORT=":9090"
INBOX451_SERVER_SMTP_PORT=":2025"
INBOX451_DATABASE_URL="postgres://user:pass@host:5432/dbname"
INBOX451_LOGGING_LEVEL="debug"
```

## Running the Server

```shell
# Using make (recommended)
make dev

# Manual start
docker compose up -d db
go run cmd/mercury/main.go
```

## API Examples

### Account Management

Create an Account:
```shell
curl -X POST http://localhost:8080/accounts \
  -H "Content-Type: application/json" \
  -d '{"name": "Test Account"}'
```

List Accounts (with pagination):
```shell
curl "http://localhost:8080/accounts?limit=10&offset=0"
```

### Inbox Management

Create an Inbox:
```shell
curl -X POST http://localhost:8080/accounts/1/inboxes \
  -H "Content-Type: application/json" \
  -d '{"email": "inbox@example.com"}'
```

List Inboxes for Account (with pagination):
```shell
curl "http://localhost:8080/accounts/1/inboxes?limit=10&offset=0"
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
├── Makefile               # Build and development commands
├── cmd/
│   └── main.go           # Main application entry point
├── config/
│   └── default.yaml      # Default configuration
├── bruno/                # API testing collections
│   └── mercury-api/
│       ├── accounts/     # Account-related requests
│       ├── inboxes/      # Inbox-related requests
│       ├── messages/     # Message-related requests
│       ├── rules/        # Rule-related requests
│       └── environments/ # Environment configurations
├── internal/
│   ├── api/             # HTTP API implementation
│   │   ├── handlers.go
│   │   ├── middleware.go
│   │   └── server.go
│   ├── config/          # Configuration management
│   ├── core/            # Business logic
│   │   ├── accounts.go
│   │   ├── inboxes.go
│   │   ├── messages.go
│   │   └── rules.go
│   ├── logger/          # Logging package
│   ├── models/          # Data models and pagination
│   ├── smtp/            # SMTP server implementation
│   ├── imap/           # IMAP server implementation
│   └── storage/         # Database operations
│       ├── queries.sql
```

### API Testing with Bruno

The project includes a comprehensive API test suite using [Bruno](https://www.usebruno.com/). Bruno is a fast and git-friendly API client that allows testing and validating API endpoints.

To use the Bruno collection:
1. Install Bruno from [usebruno.com](https://www.usebruno.com/)
2. Open the `bruno` folder in Bruno
3. Run requests individually or use collections
4. Tests will automatically validate responses

## Architecture

The application follows a layered architecture:
- API Layer: HTTP handlers and SMTP server
- Service Layer: Business logic in core package
- Storage Layer: Database operations

Database schema is inspired by [Archiveopteryx](https://archiveopteryx.org/db/)
