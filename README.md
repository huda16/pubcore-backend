# PubCore Backend API

A RESTful publishing platform API built with **Go Gin**, **GORM**, **PostgreSQL**, and **JWT authentication**.

## Features

- 🔐 JWT-based authentication (Register, Login, Logout)
- 📚 Full CRUD for **Books**, **Authors**, and **Publishers**
- 🔍 Pagination, filtering, and sorting on all list endpoints
- 🗄️ Auto-migration + demo data seeder
- 📖 Swagger UI at `/swagger/index.html`
- 🐳 Docker & docker-compose support

## Entities & Relationships

```
Author (1) ──── (N) Books (N) ──── (1) Publisher
```

- One Author can have many Books
- One Book must have exactly one Publisher

## Quick Start

### 1. Prerequisites

- Go 1.22+
- External PostgreSQL Database (Online or Local)
- [swag CLI](https://github.com/swaggo/swag): `go install github.com/swaggo/swag/cmd/swag@v1.16.3`

### 2. Configure environment

```bash
cp .env.example .env
# Edit .env with your Online DB credentials and JWT secret
```

### 3. Generate Swagger docs

```bash
swag init
```

### 4. Run the server

```bash
go run main.go
```

Server starts at `http://localhost:8080`  
Swagger UI: `http://localhost:8080/swagger/index.html`

---

## Docker

```bash
# Copy env file
cp .env.example .env
# Ensure .env has your Online DB credentials

# Start the App
docker-compose up --build
```

---

## API Endpoints

### Auth (public)

| Method | Path | Description |
|--------|------|-------------|
| POST | `/auth/register` | Register new user |
| POST | `/auth/login` | Login, returns JWT |
| POST | `/auth/logout` | Logout (bearer token required) |

### Authors (🔐 JWT required)

| Method | Path | Description |
|--------|------|-------------|
| GET | `/authors?page=1&limit=10&sort=name&order=asc&name=rowling&nationality=british` | List authors |
| GET | `/authors/:id` | Get author + books |
| POST | `/authors` | Create author |
| PUT | `/authors/:id` | Update author |
| DELETE | `/authors/:id` | Delete author |

### Publishers (🔐 JWT required)

| Method | Path | Description |
|--------|------|-------------|
| GET | `/publishers?page=1&limit=10&name=penguin` | List publishers |
| GET | `/publishers/:id` | Get publisher + books |
| POST | `/publishers` | Create publisher |
| PUT | `/publishers/:id` | Update publisher |
| DELETE | `/publishers/:id` | Delete publisher |

### Books (🔐 JWT required)

| Method | Path | Description |
|--------|------|-------------|
| GET | `/books?page=1&limit=10&genre=fantasy&author_id=1&year=1997` | List books |
| GET | `/books/:id` | Get book with author + publisher |
| POST | `/books` | Create book |
| PUT | `/books/:id` | Update book |
| DELETE | `/books/:id` | Delete book |

---

## Demo Credentials (seeded)

| Email | Password | Role |
|-------|----------|------|
| admin@pubcore.io | admin123 | Admin |
| editor@pubcore.io | editor123 | Editor |

---

## Tech Stack

- **Framework**: Go Gin
- **ORM**: GORM
- **Database**: PostgreSQL
- **Auth**: JWT (dgrijalva/jwt-go)
- **Validation**: govalidator
- **Docs**: Swagger (swaggo/swag)
---

## Testing

The project includes unit tests for helpers and business logic. To run the tests, use the following command:

```bash
go test ./... -v
```

This will run all tests in the repository and provide detailed output.
