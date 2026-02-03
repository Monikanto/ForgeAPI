# Go REST Backend

A robust REST API backend built with Go, Gin, and PostgreSQL.

## Features
- **User Authentication**: Secure registration and login using JWT and Bcrypt.
- **Product Management**: Complete CRUD operations for products.
- **Order Processing**: Transactional order creation with stock checks.
- **Infrastructure**: Docker and Docker Compose ready.
- **Middleware**: Auth middleware for protected routes.

## Quick Start

### Prerequisites
- Docker & Docker Compose
- Go 1.25+ (for local development)

### Running with Docker (Recommended)
```bash
docker-compose up --build
```
The API will be available at `http://localhost:8000`.

### Running Locally
1. **Database Setup**: Ensure PostgreSQL is running and create a database named `gorest`.
2. **Configuration**: Create a `.env` file (or set environment variables) with:
   ```env
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=password
   DB_NAME=gorest
   DB_SSLMODE=disable
   JWT_SECRET=your_secret_key
   PORT=8000
   ```
3. **Run Migrations**: Execute the SQL files in `migrations/` folder against your database.
4. **Install Dependencies**:
   ```bash
   go mod tidy
   ```
5. **Start Server**:
   ```bash
   go run cmd/api/main.go
   ```

## API Endpoints

### Auth
- `POST /api/register` - Register a new user
- `POST /api/login` - Login and receive JWT

### Products
- `GET /api/products` - List all products
- `GET /api/products/:id` - Get product details
- `POST /api/products` - Create product (Auth required)
- `PUT /api/products/:id` - Update product (Auth required)
- `DELETE /api/products/:id` - Delete product (Auth required)

### Orders
- `POST /api/orders` - Create a new order (Auth required)
- `GET /api/orders` - List user orders (Auth required)

## Project Structure
- `cmd/api`: Entry point
- `internal/auth`: JWT logic
- `internal/config`: Configuration loading
- `internal/db`: Database connection
- `internal/handler`: HTTP handlers
- `internal/middleware`: Auth middleware
- `internal/model`: Data models
- `internal/repository`: Database access layer
- `internal/routes`: Route definitions
- `internal/service`: Business logic
