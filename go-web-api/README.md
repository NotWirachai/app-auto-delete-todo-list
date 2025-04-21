# Go Web API

A complete RESTful API built with Go, Gorilla Mux, and PostgreSQL.

## Features

- RESTful API structure with Gorilla Mux router
- PostgreSQL database integration with GORM
- JWT-based authentication and authorization
- Middleware for logging, CORS, and authentication
- Environment-based configuration
- Modular architecture

## Prerequisites

- Go 1.16 or higher
- PostgreSQL 12 or higher
- Git (for version control)

## Installation

1. Clone the repository:

```bash
git clone https://github.com/yourusername/go-web-api.git
cd go-web-api
```

2. Install dependencies:

```bash
go mod download
```

3. Set up the PostgreSQL database:

```bash
createdb go_web_api
```

4. Build the application:

```bash
go build -o api
```

## Usage

### Running the API

```bash
./api
```

By default, the server will start on port 8080. You can change this by setting the `PORT` environment variable.

### Environment Variables

The API can be configured using the following environment variables:

- `PORT`: The port the server will listen on (default: 8080)
- `HOST`: The host the server will bind to (default: 0.0.0.0)
- `READ_TIMEOUT`: HTTP read timeout in seconds (default: 10)
- `WRITE_TIMEOUT`: HTTP write timeout in seconds (default: 10)
- `DEBUG`: Enable debug mode (default: false)

#### Database Configuration

- `DB_HOST`: PostgreSQL host (default: localhost)
- `DB_PORT`: PostgreSQL port (default: 5432)
- `DB_USER`: PostgreSQL user (default: postgres)
- `DB_PASSWORD`: PostgreSQL password (default: postgres)
- `DB_NAME`: PostgreSQL database name (default: go_web_api)
- `DB_SSLMODE`: PostgreSQL SSL mode (default: disable)

#### JWT Configuration

- `JWT_SECRET`: Secret key for JWT signing (default: your-secret-key)
- `JWT_EXPIRY`: JWT token expiry in minutes (default: 60)

### API Endpoints

#### Public Endpoints

| Method | Endpoint           | Description             |
| ------ | ------------------ | ----------------------- |
| GET    | /health            | Health check            |
| POST   | /api/auth/register | Register a new user     |
| POST   | /api/auth/login    | Login and get JWT token |

#### Protected Endpoints (Requires JWT Token)

| Method | Endpoint       | Description                  |
| ------ | -------------- | ---------------------------- |
| GET    | /api/users/me  | Get current user information |
| GET    | /api/items     | Get all items                |
| GET    | /api/items/:id | Get an item by ID            |
| POST   | /api/items     | Create a new item            |
| PUT    | /api/items/:id | Update an existing item      |
| DELETE | /api/items/:id | Delete an item               |

#### Admin Endpoints (Requires Admin Role)

| Method | Endpoint         | Description                |
| ------ | ---------------- | -------------------------- |
| GET    | /api/admin/users | Get all users (admin only) |

### Example Requests

#### Register a new user

```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"user1","email":"user1@example.com","password":"password123","first_name":"John","last_name":"Doe"}'
```

#### Login

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"user1","password":"password123"}'
```

#### Get all items (with JWT token)

```bash
curl -X GET http://localhost:8080/api/items \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### Create a new item (with JWT token)

```bash
curl -X POST http://localhost:8080/api/items \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{"title":"New Item","description":"This is a new item","price":49.99}'
```

## Project Structure

```
go-web-api/
├── api/         # API routes
├── config/      # Application configuration
├── database/    # Database connection and utilities
├── handlers/    # Request handlers
├── middleware/  # Middleware (logging, auth, etc.)
├── models/      # Data models
├── main.go      # Application entry point
├── go.mod       # Go modules file
└── README.md    # This file
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.
