# Contributing to Inbox451

## Use Issues for Everything

- For a small change, just send a PR.
- For bigger changes open an issue for discussion before sending a PR.
- PR should have:
  - Test case
  - Documentation
  - Example (If it makes sense)
- You can also contribute by:
  - Reporting issues
  - Suggesting new features or enhancements
  - Improve/fix documentation

## Development Environment

### Requirements

- Go 1.23+
- Node.js 20+
- PNPM 8+
- Docker and Docker Compose
- PostgreSQL 16 (for local development)

## Development Workflow

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests: `make test`
5. Run linters: `make lint`
6. Submit a pull request

## Code Style

We use several tools to ensure consistent code style:

### Backend (Go)
- `golangci-lint` for linting
- `gofmt` for formatting
- Run `make lint` before committing

### Frontend (Vue/TypeScript)
- ESLint with TypeScript config
- Prettier for formatting
- Run `cd frontend && pnpm lint` before committing

## Testing

- Run backend tests: `make test`
- Run frontend tests: `cd frontend && pnpm test`
- Run integration tests: `make test-integration`

## Building

- Development build: `make dev`
- Production build: `make pack-bin`
- Docker image: `make docker-build`

## Release Process

See [RELEASE.md](./RELEASE.md) for details on our release process.
