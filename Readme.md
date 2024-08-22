# DefiFundr - A decentralized crowdfunding platform for the Ethereum blockchain

![GitHub Workflow Status (with event)](https://img.shields.io/github/actions/workflow/status/demola234/deFICrowdFunding-Backend/test.yml)
![GitHub go.mod Go version (branch & subdirectory of monorepo)](https://img.shields.io/github/go-mod/go-version/demola234/deFICrowdFunding-Backend/main)
![GitHub issues](https://img.shields.io/github/issues/demola234/deFICrowdFunding-Backend)
![GitHub Repo stars](https://img.shields.io/github/stars/demola234/deFICrowdFunding-Backend)

## What is DefiFundr?

DefiFundr is a decentralized crowdfunding platform built on the Ethereum blockchain. It allows users to create and contribute to crowdfunding campaigns, and allows campaign creators to set a funding goal and deadline. If the funding goal is met before the deadline, the campaign is successful and the funds are released to the campaign creator. If the funding goal is not met before the deadline, the campaign is unsuccessful and the funds are returned to the contributors.

## Installation

```bash
git clone
cd deFICrowdFunding-Backend
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

## API Documentation

<!-- swagger logo and link to view -->

![Swagger](https://img.shields.io/badge/-Swagger-%23Clojure?style=for-the-badge&logo=swagger&logoColor=white&link=https://defifundr-hyper.koyeb.app/swagger/index.html)
[View API Documentation](https://defifundr-hyper.koyeb.app/swagger/index.html)

## Database Documentation

<!-- dbdiagram logo and link to view -->

[DB Diagram](https://dbdocs.io/kolawoleoluwasegun567/DefiFundr)

## Smart Contract Address

<!-- etherscan logo and link to view -->

![Ethereum](https://img.shields.io/badge/Ethereum-3C3C3D?style=for-the-badge&logo=Ethereum&logoColor=white)(https://sepolia.etherscan.io/address/0x574Bc33136180f0734fc3fa55379e9e28701395E#code)

## API Endpoints

| Endpoint                           |       Functionality        | HTTP method |
| ---------------------------------- | :------------------------: | :---------: |
| /api/v1/campaigns                  |   Create a new campaign    |    POST     |
| /api/v1/campaigns                  |     Get all campaigns      |     GET     |
| /api/v1/campaigns/:id              |    Get a campaign by id    |     GET     |
| /api/v1/campaigns/owner            |  Get a campaign by owner   |     GET     |
| /api/v1/campaigns/donation/:id     |   Get a campaign donors    |     GET     |
| /api/v1/campaigns/donate           |    Donate to a campaign    |    POST     |
| /api/v1/campaigns/withdraw         |  Withdraw from a campaign  |    POST     |
| /api/v1/campaigns/myDonations      |      Get my donations      |     GET     |
| /api/v1/campaigns/categories       |     Get all categories     |     GET     |
| /api/v1/campaigns/categories/:id   | Get campaigns by category  |     GET     |
| /api/v1/campaigns/search           |  Search campaigns by name  |     GET     |
| /api/v1/campaignsTypes             |   Get all campaign types   |     GET     |
| /api/v1/campaigns/latestCampaigns  |    Get latest campaigns    |     GET     |
| /api/v1/user                       |      Get user details      |     GET     |
| /api/v1/user                       |    Update user details     |    POST     |
| /api/v1/userAddress                |    Get user by address     |    POST     |
| /api/v1/user/avatar                |     Set profile avatar     |     GET     |
| /api/v1/user/avatar/set            |       Select avatar        |    POST     |
| /api/v1/user/biometrics            |       Set biometrics       |    POST     |
| /api/v1/user/logout                |        Logout user         |    POST     |
| /api/v1/user/password/change       |      Change password       |    POST     |
| /api/v1/user/password              |      Create password       |    POST     |
| /api/v1/user/password/reset        |       Reset password       |    POST     |
| /api/v1/user/password/reset/verify | Verify password reset code |    POST     |
| /api/v1/user/verify                |        Verify user         |    POST     |
| /api/v1/user/verify/resend         |  Resend verification code  |    POST     |
| /api/v1/user/checkUsername         |  Check if username exists  |    POST     |
| /api/v1/user/privatekey            |    Get user private key    |    POST     |
| /api/v1/user/login                 |         Login user         |    POST     |
| /api/v1/user/renewAccess           |     Renew access token     |    POST     |
| /api/v1/currentPrice               |   Get current ETH price    |     GET     |
