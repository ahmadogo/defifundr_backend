# Contributing to DefiFundr

Thank you for your interest in contributing to DefiFundr! This document provides guidelines and instructions for contributing to the project, explains our architecture, and details the data flow within the application.

## Table of Contents

- [Architecture Overview](#architecture-overview)
- [Project Structure](#project-structure)
- [Data Flow](#data-flow)
- [Development Environment Setup](#development-environment-setup)
- [Coding Standards](#coding-standards)
- [Testing Requirements](#testing-requirements)
- [Pull Request Process](#pull-request-process)
- [Issue Guidelines](#issue-guidelines)
- [Documentation](#documentation)
- [License](#license)

## Architecture Overview

DefiFundr follows the **Hexagonal Architecture** (also known as Ports and Adapters) pattern. This architecture emphasizes separation of concerns by dividing the codebase into three main layers:

1. **Core Domain** - Contains the business logic and domain models
2. **Ports** - Interfaces that define how the core domain interacts with external systems
3. **Adapters** - Implementations that connect the core domain to external systems

This approach allows us to:
- Keep the business logic independent of external concerns
- Easily replace or update external dependencies
- Test business logic in isolation
- Maintain a clean and organized codebase

### Key Components

- **API Layer**: RESTful endpoints for client communication
- **Service Layer**: Business logic implementation
- **Repository Layer**: Data access operations
- **Domain Layer**: Core business entities and rules
- **Infrastructure**: Cross-cutting concerns like logging, validation, and authentication

## Project Structure

The project follows a structured organization:

```
├── cmd                      # Application entrypoints
│   ├── api                  # API server
│   │   ├── docs             # Swagger documentation (generated with `make swagger`)
│   │   └── main.go          # API server entry point
│   └── seed                 # Database seeding tool
├── config                   # Application configuration
├── db                       # Database-related code
│   ├── migrations           # SQL migration files (managed with goose)
│   ├── query                # SQL query files (for sqlc)
│   └── sqlc                 # Generated database code (via `make sqlc`)
├── docs                     # Project documentation
│   └── db_diagram           # Database schema visualization
├── infrastructure           # Cross-cutting concerns
│   ├── common               # Shared utilities
│   │   ├── logging          # Logging utilities
│   │   ├── utils            # General utilities
│   │   └── validation       # Input validation
│   ├── hash                 # Password hashing
│   └── middleware           # HTTP middleware
├── internal                 # Application core code
│   ├── adapters             # Implementation of ports
│   │   ├── dto              # Data Transfer Objects
│   │   │   ├── request      # Request models
│   │   │   └── response     # Response models
│   │   ├── handlers         # HTTP handlers
│   │   ├── repositories     # Data access implementations
│   │   └── routers          # HTTP route definitions
│   └── core                 # Core business logic
│       ├── domain           # Domain models
│       ├── ports            # Interface definitions
│       └── services         # Business logic implementation
├── pkg                      # Reusable packages
│   ├── app_errors           # Error handling
│   ├── pagination           # Pagination utilities
│   ├── random               # Random data generation
│   ├── token_maker          # Authentication token handling
│   └── tracing              # Distributed tracing
├── scripts                  # Utility scripts
├── smart-contracts          # Smart contract source code
│   └── ethereum             # Ethereum contracts
├── test                     # Test code
│   ├── e2e                  # End-to-end tests
│   ├── integration          # Integration tests
│   └── unit                 # Unit tests
├── .env                     # Environment variables (loaded by Makefile)
├── docker-compose.yml       # Docker services configuration
├── Dockerfile               # Application container definition
├── go.mod                   # Go module definition
├── Makefile                 # Build and development commands
└── sqlc.yaml                # SQLC configuration
```

## Data Flow

The data flow in DefiFundr follows the principles of hexagonal architecture:

1. **HTTP Request Lifecycle:**
   - Request arrives at the API server (started with `make server` or `make air`)
   - Middleware processes the request (authentication, rate limiting)
   - Router directs the request to the appropriate handler
   - Handler validates input and transforms it into domain objects
   - Service layer applies business logic
   - Repository layer performs data access operations
   - Response flows back through the layers to the client

2. **Authentication Flow:**
   - User submits credentials via auth endpoints
   - `auth_handler.go` validates the request format
   - `auth_service.go` validates credentials using the hash package
   - On successful validation, token_maker creates a PASETO token
   - Session is created and stored via `session_repository.go`
   - Token is returned to the client for subsequent requests
   - Future requests are authenticated via `auth_middleware.go`

3. **Database Access Flow:**
   - Service layer calls repository methods
   - Repository uses generated sqlc code (created via `make sqlc`)
   - SQL queries defined in `db/query/*.sql` files
   - Database migrations in `db/migrations/*.sql` (managed with `make migrate-*` commands)
   - Database can be seeded with test data using `make seed`

4. **Blockchain Interaction Flow:**
   - Service layer calls blockchain adapter methods
   - Adapter uses generated Go bindings (created via `make gencontract`)
   - Smart contracts (Payroll.sol, Invoice.sol) define on-chain logic
   - Transactions are signed and submitted to the blockchain
   - Events are monitored for transaction confirmation

### Visual Data Flow

```
Client Request → Router → Handler → Service → Repository → Database
       ↑                     ↓         ↓          ↓
       └─────────────────────┴─────────┴──────────┘
                        Response
```

## Development Environment Setup

### Prerequisites

- Go 1.21+
- Docker and Docker Compose
- PostgreSQL 14+
- Make
- Migrate CLI

### Getting Started

1. Clone the repository:
   ```
   git clone https://github.com/your-org/defifundr.git
   cd defifundr
   ```

2. Install required development tools:
   ```
   make install-tools
   ```
   This installs:
   - air (for hot reloading)
   - goose (for migrations)
   - sqlc (for SQL code generation)
   - mockgen (for test mocks)
   - golangci-lint (for linting)

3. Start the development environment:
   ```
   make docker-up
   ```
   
4. Set up the database:
   ```
   make postgres
   make createdb
   ```

5. Run database migrations:
   ```
   make migrate-up
   ```

6. Generate SQL code:
   ```
   make sqlc
   ```

7. Start the API server (choose one):
   ```
   make server    # Standard mode
   make air       # Hot reload mode
   ```

8. Seed the database with test data (optional):
   ```
   make seed
   ```

## Coding Standards

DefiFundr follows these coding standards:

1. **Go Guidelines**:
   - Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
   - Use `gofmt` for code formatting
   - Run `make lint` to ensure code quality
   - Maintain 80-100% test coverage for business logic
   - Document all exported functions, types, and packages

2. **Commit Messages**:
   - Follow the [Conventional Commits](https://www.conventionalcommits.org/) specification
   - Link to relevant issues where applicable

3. **Error Handling**:
   - Use the custom `app_errors` package for application errors
   - Provide meaningful error messages
   - Ensure proper error propagation

4. **Testing**:
   - Write unit tests for business logic (run with `make test`)
   - Write integration tests for external dependencies
   - Use table-driven tests where appropriate
   - Use mocks (generate with `make mock`) for dependency isolation

5. **SQL & Database**:
   - Define queries in `db/query/*.sql` files
   - Generate Go code with `make sqlc`
   - Create migrations with `make migrate-create`
   - Test migrations with `make migrate-up` and `make migrate-down`

6. **Solidity & Smart Contracts**:
   - Follow Solidity best practices
   - Regenerate Go bindings with `make gencontract` after contract changes

## Testing Requirements

Before submitting a PR, ensure:

1. All unit tests pass:
   ```
   make test
   ```

2. Generate and update mocks if necessary:
   ```
   make mock
   ```

3. Run linting to check code quality:
   ```
   make lint
   ```

4. Ensure your code builds without errors:
   ```
   make build
   ```

## Pull Request Process

Please follow these steps for submitting contributions:

1. Create a feature branch from `main`
2. Implement your changes with appropriate tests
3. Ensure all tests pass and code coverage requirements are met
4. Update documentation as necessary
5. Submit a pull request following our [PR guidelines](ISSUE_PR_GUIDELINES.md)

## Issue Guidelines

For detailed guidance on creating and managing issues, please refer to our [Issue and PR Guidelines](ISSUE_PR_GUIDELINES.md).

## Documentation

We value comprehensive documentation:

1. **Code Documentation**:
   - Document all exported functions, types, and packages
   - Provide examples for complex functionality
   - Update API documentation when endpoints change

2. **Architecture Documentation**:
   - Update diagrams when architecture changes
   - Document design decisions and trade-offs
   - Maintain up-to-date data flow documentation
   - Update database documentation:
     ```
     make db_docs     # Generate DB documentation
     make db_schema   # Generate SQL schema from DBML
     ```
   - Update API documentation:
     ```
     make swagger     # Generate Swagger documentation
     ```

## License

By contributing to DefiFundr, you agree that your contributions will be licensed under the project's license.

## Command Reference

DefiFundr provides a comprehensive set of Makefile commands to streamline development. You can see all available commands with:

```
make help
```

### Key Commands by Category:

#### Docker Commands
- `make docker-up` - Start all Docker containers
- `make docker-down` - Stop and remove all Docker containers
- `make docker-logs` - View logs from all Docker containers
- `make docker-ps` - List running Docker containers
- `make docker-build` - Rebuild Docker images
- `make docker-restart` - Restart all Docker containers

#### Database Commands
- `make postgres` - Run PostgreSQL container
- `make createdb` - Create database
- `make dropdb` - Drop database
- `make dockerlogs` - View PostgreSQL container logs

#### Migration Commands
- `make migrate-create` - Create a new migration file
- `make migrate-up` - Run all pending migrations
- `make migrate-down` - Revert the last migration
- `make migrate-status` - Show migration status
- `make migrate-reset` - Revert all migrations

#### Development Commands
- `make sqlc` - Generate SQL code with sqlc
- `make mock` - Generate mock code for testing
- `make server` - Run the API server
- `make air` - Run the server with hot reload
- `make test` - Run tests
- `make lint` - Run linters
- `make build` - Build the application
- `make clean` - Clean build artifacts
- `make seed` - Seed the database with test data

#### Smart Contract Commands
- `make gencontract` - Generate Go bindings for smart contracts

#### Documentation Commands
- `make db_docs` - Generate DB documentation with dbdocs
- `make db_schema` - Generate SQL schema from DBML
- `make swagger` - Generate Swagger documentation

## Questions and Support

If you have questions about contributing, please:
1. Run `make help` to see available commands
2. Check existing documentation
3. Open a discussion on GitHub
4. Reach out to project maintainers

Thank you for contributing to DefiFundr!