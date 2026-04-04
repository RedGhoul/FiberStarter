# Contributing to FiberStarter

Thank you for considering contributing to FiberStarter! This guide will help you get started.

## Development Setup

### Prerequisites

- Go 1.24+
- MySQL 8.0+ (or Docker)

### Quick Start

```bash
# Clone and enter the project
git clone <repository-url>
cd FiberStarter

# Copy environment config
cp .env.example .env
# Edit .env with your database credentials

# Option A: Run with Docker (no local MySQL needed)
make docker-up

# Option B: Run locally
make run
```

## Making Changes

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/your-feature`
3. Make your changes
4. Run tests: `make test`
5. Commit with a clear message: `git commit -m 'Add your feature'`
6. Push: `git push origin feature/your-feature`
7. Open a Pull Request

## Code Style

- Follow standard Go conventions (`gofmt`, `go vet`)
- Use the existing project patterns (repository pattern, controller structure)
- Keep controllers thin -- business logic belongs in repos or utils
- Write tests for new functionality

## Project Structure

When adding a new feature, you'll typically touch these layers:

```
models/     -> Define your data struct
migrations/ -> Create Up/Down SQL files
repos/      -> Add database operations
controllers/-> Add HTTP handlers
views/      -> Add templates
routes/     -> Register routes
```

See [docs/ADDING_FEATURES.md](docs/ADDING_FEATURES.md) for a detailed walkthrough.

## Running Tests

```bash
make test
```

## Reporting Issues

- Use the GitHub issue tracker
- Include steps to reproduce the issue
- Include Go version and OS information
