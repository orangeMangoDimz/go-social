# 🚀 Gopher Social API

**NOTE**: This project is for learning purposes only

A modern social media API built with Go, implementing **Clean Architecture** principles and **Domain-Driven Design (DDD)** patterns. This REST API provides comprehensive functionality for user management, posts, comments, and social interactions.

## 📋 Table of Contents

- [🏗️ Architecture Overview](#️-architecture-overview)
- [📁 Project Structure](#-project-structure)
- [🔧 Getting Started](#-getting-started)
- [📚 API Documentation](#-api-documentation)
- [🎯 Core Features](#-core-features)
- [🧪 Testing](#-testing)
- [🚀 Deployment](#-deployment)
- [🤝 Contributing](#-contributing)

## 🏗️ Architecture Overview

This project follows **Clean Architecture** principles with **Domain-Driven Design (DDD)**, ensuring:

- **Separation of Concerns**: Each layer has a single responsibility
- **Dependency Inversion**: Dependencies point inward toward the domain
- **Testability**: Easy to unit test business logic
- **Maintainability**: Clear boundaries between layers
- **Scalability**: Easy to extend and modify

### 🔄 Layered Architecture

```
┌─────────────────────────────────────────────────────────┐
│                  🌐 Presentation Layer                   │
│              (HTTP Handlers & Routes)                   │
├─────────────────────────────────────────────────────────┤
│                  📋 Application Layer                    │
│               (Use Cases & Services)                    │
├─────────────────────────────────────────────────────────┤
│                    🎯 Domain Layer                       │
│              (Entities & Business Rules)                │
├─────────────────────────────────────────────────────────┤
│                 🗄️ Infrastructure Layer                  │
│           (Database, Cache, External Services)          │
└─────────────────────────────────────────────────────────┘
```

### 🎯 Domain-Driven Design Concepts

- **Entities**: Core business objects with identity (`User`, `Post`, `Comment`)
- **Value Objects**: Immutable objects without identity (`Password`, `Email`)
- **Aggregates**: Cluster of entities and value objects
- **Services**: Domain operations that don't belong to specific entities
- **Repositories**: Abstract data access layer

## 📁 Project Structure

```
go-social/
├── 🚀 cmd/                          # Application entry points
│   ├── api/                         # Main API application
│   │   ├── main.go                  # Application bootstrap
│   │   └── *_test.go               # Integration tests
│   └── migrate/                     # Database migration tool
│       └── migrations/              # SQL migration files
│
├── 🏗️ internal/                     # Private application code
│   ├── 🔐 auth/                     # Authentication infrastructure
│   │   ├── auth.go                 # Auth interfaces
│   │   ├── jwt.go                  # JWT implementation
│   │   └── mocks.go                # Auth mocks for testing
│   │
│   ├── ⚙️ config/                   # Configuration management
│   │   └── config.go               # App configuration structs
│   │
│   ├── 🗄️ db/                       # Database setup and utilities
│   │   ├── db.go                   # Database connection
│   │   └── seed.go                 # Database seeding
│   │
│   ├── 🎯 entities/                 # 🔥 DOMAIN LAYER
│   │   ├── comments/               # Comment domain entity
│   │   ├── pagination/             # Pagination value objects
│   │   ├── payload/                # Request/Response DTOs
│   │   ├── posts/                  # Post domain entities
│   │   └── users/                  # User domain entities
│   │
│   ├── 🌐 server/                   # 🔥 PRESENTATION LAYER
│   │   └── http/                   # HTTP transport layer
│   │       ├── handler/            # HTTP request handlers
│   │       │   ├── auth/           # Authentication endpoints
│   │       │   ├── health/         # Health check endpoints
│   │       │   ├── posts/          # Post management endpoints
│   │       │   └── users/          # User management endpoints
│   │       ├── middleware/         # HTTP middlewares
│   │       ├── protocol/           # HTTP utilities (JSON, errors)
│   │       ├── route.go            # Route definitions
│   │       └── server.go           # HTTP server setup
│   │
│   ├── 📋 service/                  # 🔥 APPLICATION LAYER
│   │   ├── auth/                   # Authentication services
│   │   ├── comments/               # Comment business logic
│   │   ├── domain/                 # Domain services
│   │   │   ├── followers/          # Follower domain service
│   │   │   ├── posts/              # Post domain service
│   │   │   ├── roles/              # Role domain service
│   │   │   ├── users/              # User domain service
│   │   │   └── service.go          # Service interfaces
│   │   └── service.go              # Service layer interfaces
│   │
│   ├── 🗄️ storage/                  # 🔥 INFRASTRUCTURE LAYER
│   │   ├── cache/                  # Caching layer (Redis)
│   │   │   ├── redis.go            # Redis implementation
│   │   │   ├── storage.go          # Cache interfaces
│   │   │   ├── mocks.go            # Cache mocks
│   │   │   └── users/              # User-specific cache
│   │   ├── postgres/               # Database persistence
│   │   │   ├── comments/           # Comment repository
│   │   │   ├── followers/          # Follower repository
│   │   │   ├── pagination/         # Pagination helpers
│   │   │   ├── posts/              # Post repository
│   │   │   ├── roles/              # Role repository
│   │   │   ├── users/              # User repository
│   │   │   ├── storage.go          # Repository interfaces
│   │   │   └── mocks.go            # Repository mocks
│   │   └── storage.go              # Storage layer interfaces
│   │
│   ├── 📧 mailer/                   # Email service infrastructure
│   ├── 🛡️ ratelimiter/             # Rate limiting infrastructure
│   └── 🌍 env/                     # Environment variable utilities
│
├── 📖 docs/                         # API Documentation
│   ├── docs.go                     # Generated Swagger docs
│   ├── swagger.json                # Swagger JSON specification
│   └── swagger.yaml                # Swagger YAML specification
│
├── 🔧 Makefile                      # Build and development commands
├── 🐳 docker-compose.yml            # Local development environment
├── 📦 go.mod                       # Go module definition
└── 📦 go.sum                       # Go module checksums
```

### 🎯 Layer Responsibilities

#### 🌐 **Presentation Layer** (`internal/server/http/handler/`)
- **Purpose**: Handle HTTP requests and responses
- **Responsibilities**:
  - Request validation and parsing
  - Response formatting
  - HTTP status code management
  - Authentication middleware
  - Route definition

#### 📋 **Application Layer** (`internal/service/`)
- **Purpose**: Orchestrate business operations and use cases
- **Responsibilities**:
  - Business workflow coordination
  - Transaction management
  - Cross-cutting concerns (logging, caching)
  - External service integration

#### 🎯 **Domain Layer** (`internal/entities/`)
- **Purpose**: Core business logic and rules
- **Responsibilities**:
  - Business entities and value objects
  - Domain business rules
  - Domain events
  - Pure business logic (no dependencies)

#### 🗄️ **Infrastructure Layer** (`internal/storage/`, `internal/auth/`, etc.)
- **Purpose**: External concerns and technical implementations
- **Responsibilities**:
  - Database operations
  - External API calls
  - File system operations
  - Framework-specific code

## 🔧 Getting Started

### Prerequisites

- **Go 1.21+**
- **PostgreSQL 14+**
- **Redis 6+** (optional, for caching)
- **Make** (for build commands)

### 🚀 Quick Start

1. **Clone the repository**
   ```bash
   git clone https://github.com/orangeMangoDimz/go-social.git
   cd go-social
   ```

2. **Set up environment variables**
   ```bash
   touch .envrc
   # Edit .env with your configuration
   # Then use it
   source .envrc
   ```

3. **Run Dev Mode**
   ```bash
   air
   ```

4. **Run database migrations**
   ```bash
   make migrate-up
   ```

5. **Start the API server**
   ```bash
   make run
   ```

6. **Visit the API documentation**
   ```
   http://localhost:8080/v1/swagger/
   ```

### 🔨 Available Make Commands

```bash
make migrate-create <migration_name>      # Generate new migration file (up and down)
make migrate-up                           # Run database migrations
make migrate-down                         # Rollback database migrations
make gen-docs                             # Generate swagger docs
```

## 📚 API Documentation

### 🔗 Endpoints Overview

| Category | Endpoint | Method | Description |
|----------|----------|---------|-------------|
| **Authentication** | `/v1/authentication/user` | POST | Register new user |
| | `/v1/authentication/token` | POST | Login and get JWT token |
| | `/v1/users/activate/{token}` | PUT | Activate user account |
| **Posts** | `/v1/posts` | POST | Create new post |
| | `/v1/posts/{id}` | GET | Get post by ID |
| | `/v1/posts/{id}` | PATCH | Update post |
| | `/v1/posts/{id}` | DELETE | Delete post |
| | `/v1/posts/feed` | GET | Get user's personalized feed |
| **Users** | `/v1/users/{id}` | GET | Get user profile |
| | `/v1/users/{id}/follow` | PUT | Follow user |
| | `/v1/users/{id}/unfollow` | PUT | Unfollow user |
| **System** | `/v1/health` | GET | Health check |

### 🔐 Authentication

The API uses **JWT Bearer tokens** for authentication:

```bash
# Get a token
curl -X POST http://localhost:8080/v1/authentication/token \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password"}'

# Use the token
curl -X GET http://localhost:8080/v1/posts/feed \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 📊 Interactive Documentation

Access the full interactive API documentation with request/response examples:

**Swagger UI**: `http://localhost:8080/v1/swagger/`

## 🎯 Core Features

### 👤 **User Management**
- User registration with email verification
- JWT-based authentication
- Role-based access control
- User profiles and social connections

### 📝 **Content Management**
- Create, read, update, delete posts
- Rich content with tags and metadata
- Comment system
- Personalized user feeds

### 🤝 **Social Features**
- Follow/unfollow users
- Personalized content feeds
- User discovery
- Social interactions

### 🛡️ **Security & Performance**
- JWT authentication with secure token handling
- Rate limiting to prevent abuse
- Input validation and sanitization
- Redis caching for improved performance
- Comprehensive error handling

### 📈 **Scalability Features**
- Pagination for large datasets
- Database indexing for performance
- Connection pooling
- Graceful shutdown handling
- Health check endpoints


## 🚀 Deployment

### 🐳 Docker Deployment

Deployed on Cloud Run (google cloud console)

### 🌐 Production Configuration

Key environment variables for production:

```env
ENV=production
ADDR=0.0.0.0:8080
DB_ADDR=postgres://user:pass@localhost/db?sslmode=require
FROM_EMAIL=hello@demomailtrap.com
MAILTRAP_API_KEY=mykey
REDIS_ADDR=localhost:6379
REDIS_DB=0
REDIS_ENABLED=false
JWT_SECRET=your-super-secure-secret
```

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
