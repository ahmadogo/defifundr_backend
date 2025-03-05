# DefiFundr - A decentralized crowdfunding platform for the Ethereum blockchain

[![GitHub Workflow Status (with event)](https://img.shields.io/github/actions/workflow/status/demola234/deFICrowdFunding-Backend/test.yml)](https://github.com/DefiFundr-Labs/defifundr_backend/actions)
![GitHub go.mod Go version (branch & subdirectory of monorepo)](https://img.shields.io/github/go-mod/go-version/demola234/deFICrowdFunding-Backend/main)
[![GitHub issues](https://img.shields.io/github/issues/demola234/deFICrowdFunding-Backend)](https://github.com/DefiFundr-Labs/defifundr_backend/issues?q=is%3Aissue%20state%3Aopen)
[![GitHub Repo stars](https://img.shields.io/github/stars/demola234/deFICrowdFunding-Backend)](https://github.com/DefiFundr-Labs/defifundr_backend/stargazers)

## What is DefiFundr?

DefiFundr is a revolutionary decentralized payroll and invoice management system that bridges the gap between traditional financial systems and blockchain technology. The platform provides a seamless, secure, and transparent solution for businesses to manage employee payments, handle freelancer invoices, and automate salary disbursements across both fiat and cryptocurrency channels.

## Installation

```bash
git clone
cd defifundr_backend
go mod download
```

## Usage

### Using Makefile

```bash
make server
```

### Using Go

```bash
go run main.go
```

### Using Air (Hot Reload)

```bash
air
```

## Testing

```bash
make test
```

### Unit Tests

```bash
go test ./...
```

### Coverage

```bash
go test -v -cover ./...
```
