# Showroom Management System - Backend

🚀 Complete Phase 1 implementation of the Showroom Management System backend, built with Go and PostgreSQL.

## 📋 Overview

This is a comprehensive backend system for managing a vehicle showroom, implementing user management and authentication in Phase 1. The system follows clean architecture principles and provides a robust foundation for the complete showroom management solution.

## ✨ Features Implemented (Phase 1)

### 🔐 Authentication & Security
- JWT-based authentication with secure session management
- Password hashing using bcrypt
- Role-based access control (admin, sales, cashier, mechanic, manager)
- Session tracking and management
- Token refresh functionality
- Logout with session invalidation

### 👥 User Management
- Complete CRUD operations for users
- User creation with comprehensive validation
- Profile management and updates
- Role-based permissions and filtering
- User searching and pagination
- Admin-only user management endpoints

### 🏗️ Technical Excellence
- Clean architecture with repository pattern
- Comprehensive error handling and validation
- Request/response DTOs with field validation
- Database migrations and seeding
- Docker containerization ready
- Complete testing suite
- API documentation

## 🗄️ Database Schema

### USERS Table
```sql
CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL, 
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    phone VARCHAR(20) NOT NULL,
    address VARCHAR(500),
    role VARCHAR(20) NOT NULL CHECK (role IN ('admin','sales','cashier','mechanic','manager')),
    salary DECIMAL(15,2),
    hire_date TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    created_by INTEGER NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    profile_image VARCHAR(500),
    notes TEXT,
    FOREIGN KEY (created_by) REFERENCES users(user_id)
);
```

### USER_SESSIONS Table
```sql
CREATE TABLE user_sessions (
    session_id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    session_token VARCHAR(500) UNIQUE NOT NULL,
    login_at TIMESTAMP DEFAULT NOW(),
    logout_at TIMESTAMP,
    ip_address VARCHAR(45),
    user_agent VARCHAR(500),
    is_active BOOLEAN DEFAULT TRUE
);
```

## 🚀 Quick Start

### Prerequisites
- Go 1.21+
- PostgreSQL 15+
- Docker & Docker Compose (optional)
- Make

### 1. Clone and Setup
```bash
git clone <repository>
cd showroom-backend

# Setup environment
make dev-setup
cp .env.example .env
# Edit .env with your configuration
```

### 2. Start with Docker (Recommended)
```bash
# Start services (PostgreSQL + Application)
make docker-run

# Check logs
make docker-logs
```

### 3. Manual Setup
```bash
# Install dependencies
make deps

# Start PostgreSQL (if not using Docker)
# Update .env with your database configuration

# Run migrations and seed data
make migration
make seed

# Start the application
make run
```

### 4. Verify Installation
```bash
# Health check
curl http://localhost:8080/api/v1/health

# Test login with default admin user
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "admin123"}'
```

## 📚 API Documentation

### Base URL
```
http://localhost:8080/api/v1
```

### Authentication Endpoints

#### POST /auth/login
Login user and receive JWT token.

**Request:**
```json
{
  "username": "admin",
  "password": "admin123"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "token_type": "Bearer",
    "expires_in": 86400,
    "expires_at": "2025-01-23T12:35:00Z",
    "user": {
      "user_id": 1,
      "username": "admin",
      "email": "admin@showroom.com",
      "full_name": "System Administrator",
      "role": "admin",
      "is_active": true
    },
    "session_id": 123
  }
}
```

#### POST /auth/logout
**Headers:** `Authorization: Bearer <token>`

#### GET /auth/me
Get current user information.
**Headers:** `Authorization: Bearer <token>`

#### GET /auth/profile
Get user profile with recent sessions.
**Headers:** `Authorization: Bearer <token>`

#### POST /auth/change-password
Change current user's password.
**Headers:** `Authorization: Bearer <token>`
```json
{
  "current_password": "oldpassword",
  "new_password": "newpassword",
  "confirm_password": "newpassword"
}
```

#### POST /auth/refresh
Refresh JWT token.
**Headers:** `Authorization: Bearer <token>`

### Admin User Management Endpoints
**Note:** All admin endpoints require `Authorization: Bearer <token>` header and admin role.

#### POST /admin/users
Create new user.
```json
{
  "username": "newuser",
  "email": "user@example.com",
  "password": "password123",
  "full_name": "New User",
  "phone": "081234567890",
  "role": "cashier"
}
```

#### GET /admin/users
List users with filtering and pagination.
**Query Parameters:**
- `role` - Filter by role
- `is_active` - Filter by status
- `search` - Search by username, email, name, phone
- `page` - Page number (default: 1)
- `limit` - Items per page (default: 10, max: 100)

#### GET /admin/users/{id}
Get user by ID.

#### PUT /admin/users/{id}
Update user information.

#### DELETE /admin/users/{id}
Soft delete user.

#### GET /admin/users/role/{role}
Get users by specific role.

#### GET /admin/users/{id}/sessions
Get user's sessions.

#### DELETE /admin/users/{id}/sessions
Revoke all user sessions.

### Health Check

#### GET /health
Check service health.

## 🧪 Testing

### Run Tests
```bash
# All tests
make test

# With coverage
make test-coverage

# Unit tests only
go test -v ./tests/unit/...
```

### Test Users (Seeded)
```
Admin User:
- Username: admin
- Password: admin123
- Role: admin

Other Test Users:
- kasir1 / kasir123 (cashier)
- mekanik1 / mekanik123 (mechanic)  
- sales1 / sales123 (sales)
- manager1 / manager123 (manager)
```

## 🛠️ Development

### Available Make Commands
```bash
make help              # Show all available commands
make build             # Build the application
make run               # Build and run
make test              # Run tests
make test-coverage     # Run tests with coverage
make lint              # Lint code
make fmt               # Format code
make deps              # Download dependencies
make migration         # Run database migrations
make seed              # Seed database with test data
make docker-build      # Build Docker image
make docker-run        # Run with Docker Compose
make docker-stop       # Stop Docker services
```

### Project Structure
```
showroom-backend/
├── cmd/server/                 # Application entry point
├── internal/
│   ├── config/                # Configuration management
│   ├── database/              # Database connection & migrations
│   ├── middleware/            # HTTP middleware
│   ├── models/                # Domain models
│   ├── dto/                   # Data transfer objects
│   ├── repositories/          # Data access layer
│   ├── services/              # Business logic layer
│   ├── handlers/              # HTTP handlers
│   ├── routes/                # Route definitions
│   └── utils/                 # Utility functions
├── tests/                     # Test files
├── scripts/                   # Development scripts
├── configs/                   # Configuration files
└── storage/                   # File storage
```

## 🔧 Configuration

### Environment Variables
```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=showroom_db

# Server
SERVER_PORT=8080
SERVER_HOST=0.0.0.0
APP_ENV=development

# JWT
JWT_SECRET_KEY=your-secret-key
JWT_EXPIRATION_HOUR=24

# Application
APP_NAME=Showroom Management System
APP_VERSION=1.0.0
```

## 🔒 Security Features

- **Password Security**: Bcrypt hashing with salt
- **JWT Security**: Signed tokens with expiration
- **Session Management**: Track and invalidate sessions
- **Role-Based Access**: Granular permission control
- **Input Validation**: Comprehensive request validation
- **SQL Injection Prevention**: Parameterized queries
- **CORS Configuration**: Cross-origin request handling

## 📈 Performance

- **Database Indexing**: Optimized query performance
- **Connection Pooling**: Efficient database connections
- **Pagination**: Handle large datasets efficiently
- **Middleware Pipeline**: Optimized request processing

## 🐳 Docker Support

### Docker Compose Services
- **PostgreSQL**: Database service
- **Application**: Go backend service
- **Redis**: Ready for future caching needs

### Commands
```bash
docker-compose up -d          # Start all services
docker-compose logs -f        # Follow logs
docker-compose down           # Stop services
docker-compose down -v        # Stop and remove volumes
```

## 🚀 Next Steps (Phase 2)

The system is architected to easily support Phase 2 features:
- **Master Data Management**: Customers, Suppliers, Vehicle Brands/Models
- **Product Categories**: Spare parts categorization  
- **Basic Inventory**: Foundation for parts management
- **Reporting**: Analytics and insights

## 🤝 Contributing

1. Follow Go best practices and conventions
2. Write comprehensive tests for new features
3. Update documentation for API changes
4. Use the provided development tools (linting, formatting)
5. Follow the established project structure

## 📝 License

This project is part of the Showroom Management System implementation.

---

**Built with ❤️ using Go, Gin, PostgreSQL, and Docker**