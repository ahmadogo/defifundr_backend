# DefiFundr Setup Guide

This guide will help you set up your development environment for contributing to DefiFundr.

## Prerequisites

Before you begin, make sure you have the following installed:

- [Go](https://golang.org/doc/install) (version 1.21 or higher)
- [Docker](https://docs.docker.com/get-docker/) and [Docker Compose](https://docs.docker.com/compose/install/)
- [Git](https://git-scm.com/downloads)
- [Solidity Compiler](https://docs.soliditylang.org/en/v0.8.17/installing-solidity.html) (for smart contract development)
- [Node.js and npm](https://nodejs.org/) (for smart contract tools)

## Initial Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/your-org/defifundr.git
   cd defifundr
   ```

2. Install required development tools:
   ```bash
   make install-tools
   ```
   
   This command installs:
   - `air` - Live reload for Go applications
   - `goose` - Database migration tool
   - `sqlc` - SQL compiler for Go
   - `mockgen` - Mock generation for testing
   - `golangci-lint` - Linting tool for Go

3. Set up environment variables:
   ```bash
   cp .env.example .env
   ```
   
   Edit the `.env` file to configure your development environment. At minimum, set:
   - `DB_SOURCE` - Database connection string
   - `JWT_SECRET` - Secret key for JWT tokens
   - `ETHEREUM_RPC_URL` - Ethereum node URL

## Database Setup

1. Start a PostgreSQL container:
   ```bash
   make postgres
   ```

2. Create the database:
   ```bash
   make createdb
   ```

3. Run migrations to set up the schema:
   ```bash
   make migrate-up
   ```

4. (Optional) Seed the database with test data:
   ```bash
   make seed
   ```

## Smart Contract Setup

1. Generate Go bindings for smart contracts:
   ```bash
   make gencontract
   ```

   This compiles the Solidity contracts and generates Go interfaces.

## Running the Application

1. Start all required Docker containers:
   ```bash
   make docker-up
   ```

2. Generate SQL code:
   ```bash
   make sqlc
   ```

3. Start the API server:
   ```bash
   # Option 1: Standard mode
   make server
   
   # Option 2: Hot reload mode (recommended for development)
   make air
   ```

4. Generate API documentation:
   ```bash
   make swagger
   ```
   
   Access the Swagger UI at http://localhost:8080/swagger/index.html

## Verifying Your Setup

1. Check that the database is running:
   ```bash
   make docker-ps
   ```

2. Verify database migrations are applied:
   ```bash
   make migrate-status
   ```

3. Run tests to ensure everything is working:
   ```bash
   make test
   ```

## Development Workflow

Once your environment is set up, follow this workflow for development:

1. Create a feature branch:
   ```bash
   git checkout -b feat/your-feature-name
   ```

2. Make your changes and write tests.

3. Run linting:
   ```bash
   make lint
   ```

4. Run tests:
   ```bash
   make test
   ```

5. Generate any updated documentation:
   ```bash
   make swagger
   ```

6. Commit your changes using conventional commits format.

7. Submit a pull request following the guidelines in [CONTRIBUTING.md](CONTRIBUTING.md).

## Troubleshooting

### Database Connection Issues

If you encounter database connection errors:

```bash
# Check if PostgreSQL container is running
make docker-ps

# View PostgreSQL logs
make dockerlogs

# Restart the PostgreSQL container
make docker-restart
```

### Smart Contract Compilation Errors

If smart contract generation fails:

1. Ensure solc (Solidity Compiler) is correctly installed
2. Check contract syntax for errors
3. Verify the correct paths in the Makefile

### Go Module Issues

If you encounter Go module issues:

```bash
# Clear Go module cache
go clean -modcache

# Update dependencies
go mod tidy
```

## Additional Resources

- [Go Documentation](https://golang.org/doc/)
- [Docker Documentation](https://docs.docker.com/)
- [Solidity Documentation](https://docs.soliditylang.org/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)

## Getting Help

If you need further assistance, refer to:

- Project-specific documentation in the `docs` directory
- Open an issue on GitHub
- Contact the project maintainers