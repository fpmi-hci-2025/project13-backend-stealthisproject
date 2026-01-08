# Railway Ticket System - Backend

Backend API for the Railway Ticket System built with Go (Golang).

## Tech Stack

- **Language**: Go 1.21+
- **Framework**: Gin
- **Database**: PostgreSQL 15
- **Authentication**: JWT
- **Documentation**: Swagger/OpenAPI
- **CI/CD**: GitHub Actions
- **Code Quality**: golangci-lint

## Project Structure

```
project13-backend-stealthisproject/
├── cmd/
│   └── server/
│       └── main.go          # Application entry point
├── internal/
│   ├── config/              # Configuration management
│   ├── database/            # Database connection and migrations
│   ├── handlers/            # HTTP handlers
│   ├── middleware/          # HTTP middleware (auth, CORS, etc.)
│   ├── models/              # Data models
│   └── repository/          # Data access layer
├── pkg/
│   └── auth/                # Authentication service
├── .github/
│   └── workflows/
│       └── ci.yml           # CI/CD pipeline
├── docker-compose.yml       # Local development environment
├── Dockerfile              # Container image definition
├── go.mod                  # Go dependencies
└── Makefile               # Build automation

```

## Features

- ✅ User registration and authentication (JWT)
- ✅ User profile management
- ✅ Route search functionality
- ✅ Order creation and management
- ✅ Payment processing
- ✅ Admin CRUD operations for routes and trains
- ✅ Swagger API documentation
- ✅ Docker support for local development
- ✅ CI/CD pipeline with GitHub Actions

## Getting Started

### Prerequisites

- Go 1.21 or higher
- PostgreSQL 15 or higher
- Docker and Docker Compose (optional)

### Local Development

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd project13-backend-stealthisproject
   ```

2. **Set up environment variables**
   ```bash
   export DATABASE_URL="postgres://postgres:postgres@localhost:5432/railway_tickets?sslmode=disable"
   export JWT_SECRET="your-secret-key-change-in-production"
   export ENVIRONMENT="development"
   export PORT="8080"
   ```

3. **Install dependencies**
   ```bash
   go mod download
   ```

4. **Run database migrations**
   The migrations run automatically on application startup.

5. **Run the application**
   ```bash
   make run
   # or
   go run ./cmd/server
   ```

   The API will be available at `http://localhost:8080`

### Using Docker Compose

1. **Start services**
   ```bash
   docker-compose up -d
   ```

2. **View logs**
   ```bash
   docker-compose logs -f app
   ```

3. **Stop services**
   ```bash
   docker-compose down
   ```

## API Documentation

**Important:** Before running the server, generate Swagger documentation:

```bash
make swagger
# or
go install github.com/swaggo/swag/cmd/swag@latest
swag init -g cmd/server/main.go -o docs
```

Once the server is running, Swagger documentation is available at:
- `http://localhost:8080/swagger/index.html`

## API Endpoints

### Authentication
- `POST /api/v1/auth/register` - Register a new user
- `POST /api/v1/auth/login` - Login and get JWT token

### Users
- `GET /api/v1/users/me` - Get current user profile (protected)
- `PUT /api/v1/users/me` - Update current user profile (protected)

### Routes
- `GET /api/v1/routes/search` - Search routes by cities and date
- `GET /api/v1/routes/:id` - Get route details

### Orders
- `POST /api/v1/orders` - Create a new order (protected)
- `GET /api/v1/orders` - Get user's orders (protected)
- `GET /api/v1/orders/:id` - Get order details (protected)
- `POST /api/v1/orders/:id/pay` - Pay for an order (protected)

### Admin (Admin only)
- `POST /api/v1/admin/routes` - Create a route
- `PUT /api/v1/admin/routes/:id` - Update a route
- `DELETE /api/v1/admin/routes/:id` - Delete a route
- `POST /api/v1/admin/trains` - Create a train
- `PUT /api/v1/admin/trains/:id` - Update a train
- `DELETE /api/v1/admin/trains/:id` - Delete a train
- `GET /api/v1/admin/orders` - Get all orders

## Testing

Run tests:
```bash
make test
# or
go test -v ./...
```

Run tests with coverage:
```bash
go test -v -race -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Linting

Run linter:
```bash
make lint
# or
golangci-lint run
```

## Building

Build the application:
```bash
make build
# or
go build -o bin/server ./cmd/server
```

## Database Schema

The database schema includes the following tables:
- `users` - User accounts
- `passengers` - Passenger profiles
- `trains` - Train information
- `carriages` - Carriage information
- `seats` - Seat information
- `stations` - Station information
- `routes` - Route information
- `route_stations` - Route-station relationships
- `orders` - Order information
- `tickets` - Ticket information

Migrations run automatically on application startup.

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DATABASE_URL` | PostgreSQL connection string | `postgres://postgres:postgres@localhost:5432/railway_tickets?sslmode=disable` |
| `JWT_SECRET` | Secret key for JWT tokens | `your-secret-key-change-in-production` |
| `ENVIRONMENT` | Environment (development/production) | `development` |
| `PORT` | Server port | `8080` |

## CI/CD

The project includes a GitHub Actions workflow (`.github/workflows/ci.yml`) that:
- Runs tests
- Runs linter (golangci-lint)
- Builds the application
