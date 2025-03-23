# DefiFundr API Documentation

This document provides a comprehensive overview of the DefiFundr API, including endpoints, authentication, request/response formats, and error handling.

## Table of Contents

- [DefiFundr API Documentation](#defifundr-api-documentation)
  - [Table of Contents](#table-of-contents)
  - [API Overview](#api-overview)
  - [Authentication](#authentication)
    - [Token Management](#token-management)
    - [Request Format](#request-format)
    - [Response Format](#response-format)
    - [Pagination](#pagination)
  - [Error Handling](#error-handling)
    - [Common Error Codes](#common-error-codes)
  - [Rate Limiting](#rate-limiting)
  - [API Versioning](#api-versioning)
  - [Swagger Documentation](#swagger-documentation)
    - [Swagger Files](#swagger-files)

## API Overview

The DefiFundr API follows RESTful principles and uses JSON for request and response payloads. The API is organized around resources and supports standard HTTP methods:

- `GET` - Retrieve resources
- `POST` - Create resources
- `PUT` - Update resources (full update)
- `PATCH` - Update resources (partial update)
- `DELETE` - Delete resources

The base URL for all API endpoints is:

```
https://api.defifundr.com/v1
```

For local development:

```
http://localhost:8080/v1
```

## Authentication

The API uses PASETO (Platform-Agnostic Security Tokens) for authentication. To access protected endpoints:

1. Obtain an access token by authenticating at `/v1/auth/login`
2. Include the token in the `Authorization` header of subsequent requests:

```
Authorization: Bearer [YOUR_TOKEN]
```

### Token Management

- **Token Expiration**: Access tokens expire after 15 minutes
- **Token Refresh**: Use the refresh token endpoint `/v1/auth/refresh` to obtain a new access token
- **Token Revocation**: Tokens can be invalidated at `/v1/auth/logout`

### Request Format

Requests with a body should use JSON format with the `Content-Type: application/json` header.

Example request:

```json
{
  "email": "user@example.com",
  "password": "SecurePassword123"
}
```

### Response Format

Successful responses return a JSON object with the following structure:

```json
{
  "status": "success",
  "data": {
    // Response data specific to the endpoint
  }
}
```

### Pagination

Endpoints that return collections support pagination with the following query parameters:

- `page`: Page number (1-based)
- `limit`: Items per page (default: 10, max: 100)

Paginated responses include metadata:

```json
{
  "status": "success",
  "data": [...],
  "meta": {
    "page": 1,
    "limit": 10,
    "total": 42,
    "pages": 5
  }
}
```

## Error Handling

Errors return appropriate HTTP status codes and a JSON object with error details:

```json
{
  "status": "error",
  "error": {
    "code": "unauthorized",
    "message": "Invalid authentication credentials"
  }
}
```

### Common Error Codes

| HTTP Status | Error Code       | Description                               |
|-------------|------------------|-------------------------------------------|
| 400         | bad_request      | Invalid request parameters                |
| 401         | unauthorized     | Authentication required or invalid token  |
| 403         | forbidden        | Permission denied                         |
| 404         | not_found        | Resource not found                        |
| 409         | conflict         | Resource conflict                         |
| 422         | validation_error | Validation failed                         |
| 429         | rate_limited     | Too many requests                         |
| 500         | server_error     | Internal server error                     |

## Rate Limiting

The API implements rate limiting to prevent abuse. Limits are applied per API key/IP address.

Rate limit headers are included in all responses:

```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 99
X-RateLimit-Reset: 1620000000
```

When rate limited, you'll receive a 429 Too Many Requests response.

## API Versioning

API versions are specified in the URL path (e.g., `/v1/users`). When breaking changes are necessary, a new API version will be released.

## Swagger Documentation

Interactive API documentation is available through Swagger UI. To generate and access the Swagger documentation:

1. Generate Swagger documentation:
   ```bash
   make swagger
   ```

2. Access the Swagger UI at:
   ```
   http://localhost:8080/swagger/index.html
   ```

The Swagger documentation provides:
- Interactive endpoint testing
- Request/response schemas
- Authentication information
- Example requests

### Swagger Files

- `cmd/api/docs/swagger.json` - Swagger specification in JSON format
- `cmd/api/docs/swagger.yaml` - Swagger specification in YAML format
- `cmd/api/docs/docs.go` - Go code for Swagger integration