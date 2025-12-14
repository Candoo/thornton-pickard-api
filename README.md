# Thornton Pickard Camera API

Go-based REST API for the Thornton Pickard camera and ephemera catalog. Provides authentication, data management, search functionality, and image uploads for historical camera collections.

## 🎯 About

This is the backend API service for the Thornton Pickard camera catalog application. It provides:

- Comprehensive camera and ephemera data management
- User authentication and authorization with JWT
- Advanced search and filtering capabilities
- Image upload and storage
- RESTful API design with Swagger documentation
- PostgreSQL database with GORM ORM

## 🔗 Related Repositories

This API works with:

- **[my-modern-react-setup](https://github.com/Candoo/my-modern-react-setup)** - React frontend
- **[Pickard-Index](https://github.com/Candoo/Pickard-Index)** - Full-stack deployment setup

**⚡ Quick Start:** For the easiest setup of the entire application, use the [Pickard-Index deployment repository](https://github.com/Candoo/Pickard-Index).

## ✨ Features

- 🚀 **Fast and efficient** Go backend with Gin framework
- 🔐 **JWT authentication** with role-based access control
- 🗄️ **PostgreSQL database** with GORM ORM
- 📚 **Auto-generated Swagger documentation** for API exploration
- 🐳 **Docker support** with Docker Compose
- 🔍 **Advanced search and filtering** with pagination
- 📸 **Image upload functionality** with file validation
- ✅ **Unit tests** with testify framework
- 🌱 **Database seeding** for development and testing
- 🛡️ **CORS middleware** for frontend integration
- 📊 **RESTful API design** following best practices
- 🔄 **Automatic migrations** with GORM
- 📝 **Structured logging** for debugging

## 🚀 Getting Started

### Prerequisites

- Go 1.21 or higher
- PostgreSQL 15 or higher
- Docker & Docker Compose (optional, for containerized setup)

### Local Development Setup

1. **Install and start PostgreSQL** (or use Docker):

   ```bash
   # Using Docker
   docker run -d \
     --name postgres \
     -e POSTGRES_PASSWORD=postgres \
     -e POSTGRES_DB=thornton_pickard \
     -p 5432:5432 \
     postgres:15-alpine
   ```

2. **Clone the repository:**

   ```bash
   git clone https://github.com/Candoo/thornton-pickard-api.git
   cd thornton-pickard-api
   ```

3. **Install dependencies:**

   ```bash
   go mod download
   ```

4. **Configure environment:**

   ```bash
   cp .env.example .env
   ```

   Edit `.env` with your configuration:

   ```env
   ENV=development
   PORT=8080
   
   # Database
   DB_HOST=localhost
   DB_USER=postgres
   DB_PASSWORD=postgres
   DB_NAME=thornton_pickard
   DB_PORT=5432
   
   # Authentication (CHANGE IN PRODUCTION!)
   JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
   
   # Uploads
   UPLOAD_DIR=./uploads
   
   # Seeding
   SEED=true
   ```

5. **Generate Swagger documentation:**

   ```bash
   # Install swag if not already installed
   go install github.com/swaggo/swag/cmd/swag@latest
   
   # Generate docs
   swag init -g cmd/api/main.go
   ```

6. **Run with database seeding:**

   ```bash
   SEED=true go run cmd/api/main.go
   ```

7. **Access the API:**
   - API Base: http://localhost:8080
   - Swagger Docs: http://localhost:8080/swagger/index.html
   - Health Check: http://localhost:8080/health

### Docker Deployment

Run the entire backend stack with Docker Compose:

```bash
# Start PostgreSQL + API
docker-compose up -d

# View logs
docker-compose logs -f api

# Stop services
docker-compose down

# Rebuild after code changes
docker-compose up -d --build
```

### Full-Stack Deployment

For the complete application (frontend + backend + database), use the [Pickard-Index deployment repository](https://github.com/Candoo/Pickard-Index).

## 📚 API Endpoints

### Authentication

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/api/v1/auth/register` | Register new user | No |
| POST | `/api/v1/auth/login` | Login and get JWT token | No |
| GET | `/api/v1/auth/profile` | Get user profile | Yes |

### Cameras

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/api/v1/cameras` | List all cameras (with pagination & search) | No |
| GET | `/api/v1/cameras/:id` | Get camera by ID | No |
| POST | `/api/v1/cameras` | Create camera | Yes |
| PUT | `/api/v1/cameras/:id` | Update camera | Yes |
| DELETE | `/api/v1/cameras/:id` | Delete camera | Admin only |

### Ephemera

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/api/v1/ephemera` | List all ephemera | No |
| GET | `/api/v1/ephemera/:id` | Get ephemera by ID | No |
| POST | `/api/v1/ephemera` | Create ephemera | Yes |
| PUT | `/api/v1/ephemera/:id` | Update ephemera | Yes |
| DELETE | `/api/v1/ephemera/:id` | Delete ephemera | Admin only |

### Manufacturers

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/api/v1/manufacturers` | List all manufacturers | No |
| GET | `/api/v1/manufacturers/:id` | Get manufacturer by ID | No |

### Uploads

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/api/v1/upload` | Upload single image | Yes |
| POST | `/api/v1/upload/multiple` | Upload multiple images | Yes |

## 🔍 Query Parameters

### Pagination

- `page` - Page number (default: 1)
- `page_size` - Items per page (default: 10, max: 100)

### Search & Filters

- `search` - Full-text search across name, manufacturer, description
- `manufacturer` - Filter by manufacturer name
- `year_from` - Filter by year from (inclusive)
- `year_to` - Filter by year to (inclusive)
- `format` - Filter by camera format
- `sort` - Sort by field (name, year_introduced, manufacturer)
- `order` - Sort order (asc, desc)

### Example Queries

```bash
# Get page 2 with 20 items
curl "http://localhost:8080/api/v1/cameras?page=2&page_size=20"

# Search for "ruby" cameras
curl "http://localhost:8080/api/v1/cameras?search=ruby"

# Filter by manufacturer and year range
curl "http://localhost:8080/api/v1/cameras?manufacturer=Thornton-Pickard&year_from=1900&year_to=1920"

# Sort by year descending
curl "http://localhost:8080/api/v1/cameras?sort=year_introduced&order=desc"

# Combined filters
curl "http://localhost:8080/api/v1/cameras?search=reflex&format=Plate&year_from=1900&sort=year_introduced"
```

## 🔐 Authentication

The API uses JWT (JSON Web Tokens) for authentication.

### 1. Register or Login

```bash
# Register a new user
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "securepassword123",
    "name": "John Doe"
  }'

# Login to get token
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "securepassword123"
  }'

# Response includes JWT token:
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "name": "John Doe",
    "role": "user"
  }
}
```

### 2. Use Token in Requests

Include the token in the `Authorization` header:

```bash
curl -X POST http://localhost:8080/api/v1/cameras \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Imperial Triple Extension",
    "manufacturer": "Thornton-Pickard",
    "year_introduced": 1905,
    "format": "Plate"
  }'
```

### Default Admin Account

After seeding, use these credentials:

- **Email:** `admin@thorntonpickard.com`
- **Password:** `admin123`

**⚠️ IMPORTANT:** Change this password immediately in production!

## 📸 Image Uploads

### Upload Single Image

```bash
curl -X POST http://localhost:8080/api/v1/upload \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "file=@/path/to/camera-image.jpg"

# Response:
{
  "url": "/uploads/uuid_timestamp.jpg",
  "filename": "camera-image.jpg",
  "size": 245678
}
```

### Upload Multiple Images

```bash
curl -X POST http://localhost:8080/api/v1/upload/multiple \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "files=@image1.jpg" \
  -F "files=@image2.jpg" \
  -F "files=@image3.jpg"

# Response:
{
  "files": [
    {
      "url": "/uploads/uuid1_timestamp.jpg",
      "filename": "image1.jpg",
      "size": 245678
    },
    {
      "url": "/uploads/uuid2_timestamp.jpg",
      "filename": "image2.jpg",
      "size": 198432
    }
  ]
}
```

### Upload Restrictions

- **Allowed formats:** JPG, JPEG, PNG, GIF, WebP
- **Max file size:** 5MB per file
- **Max files:** 10 per request (multiple upload)

## 🧪 Testing

Run the test suite:

```bash
# Run all tests
go test ./tests/... -v

# Run specific test
go test ./tests/... -run TestGetCameras -v

# Run with coverage
go test ./tests/... -cover

# Generate coverage report
go test ./tests/... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Test Coverage

The test suite covers:
- API endpoint handlers
- Authentication middleware
- Database operations
- Search and filtering
- Pagination
- Error handling

## 🌱 Database Seeding

The API automatically seeds the database on startup if `SEED=true` in `.env`.

Seed data includes:

- **Admin user** (email: `admin@thorntonpickard.com`, password: `admin123`)
- **Sample manufacturers** (Thornton-Pickard, Kodak, etc.)
- **Sample cameras** (from `seeds/cameras.json` if available)

### Custom Seed Data

Create `seeds/cameras.json` with your camera data:

```json
[
  {
    "name": "Ruby Reflex",
    "manufacturer": "Thornton-Pickard",
    "year_introduced": 1913,
    "year_discontinued": 1926,
    "format": "Plate",
    "film_size": "Quarter Plate",
    "description": "SLR camera with focal plane shutter"
  }
]
```

### Manual Seeding

```bash
# Run seeding command
SEED=true go run cmd/api/main.go

# Or use database CLI
psql -U postgres -d thornton_pickard -f seeds/initial_data.sql
```

## 🏗️ Project Structure

```
thornton-pickard-api/
├── cmd/
│   └── api/
│       └── main.go             # Application entry point
├── internal/
│   ├── models/                 # Database models
│   │   ├── camera.go           # Camera model
│   │   ├── ephemera.go         # Ephemera model
│   │   ├── manufacturer.go     # Manufacturer model
│   │   └── user.go             # User model with auth
│   ├── handlers/               # HTTP request handlers
│   │   ├── camera.go           # Camera CRUD operations
│   │   ├── ephemera.go         # Ephemera CRUD operations
│   │   ├── auth.go             # Authentication handlers
│   │   ├── manufacturer.go     # Manufacturer handlers
│   │   └── upload.go           # File upload handlers
│   ├── middleware/             # HTTP middleware
│   │   ├── cors.go             # CORS configuration
│   │   ├── auth.go             # JWT authentication
│   │   ├── pagination.go       # Pagination middleware
│   │   └── logger.go           # Request logging
│   ├── services/               # Business logic layer
│   │   ├── auth.go             # Authentication service
│   │   ├── camera.go           # Camera service
│   │   └── storage.go          # File storage service
│   ├── database/               # Database operations
│   │   ├── database.go         # DB connection & setup
│   │   └── seed.go             # Database seeding
│   └── utils/                  # Utility functions
│       ├── pagination.go       # Pagination helpers
│       ├── response.go         # API response helpers
│       └── validation.go       # Input validation
├── tests/                      # Unit tests
│   ├── camera_test.go          # Camera endpoint tests
│   ├── auth_test.go            # Auth endpoint tests
│   └── upload_test.go          # Upload endpoint tests
├── docs/                       # Swagger documentation (auto-generated)
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── seeds/                      # Seed data files
│   ├── cameras.json            # Camera seed data
│   └── manufacturers.json      # Manufacturer seed data
├── uploads/                    # Uploaded files directory
├── .env.example                # Environment template
├── .env                        # Your environment (git-ignored)
├── .gitignore
├── Dockerfile                  # Production Docker image
├── docker-compose.yml          # Docker Compose config
├── dockerignore                # Docker ignore file
├── go.mod                      # Go module definition
├── go.sum                      # Go dependencies
└── README.md                   # This file
```

## ⚙️ Configuration

### Environment Variables

```env
# Server
ENV=development                 # development or production
PORT=8080                       # Server port

# Database
DB_HOST=localhost               # Database host
DB_USER=postgres                # Database user
DB_PASSWORD=postgres            # Database password
DB_NAME=thornton_pickard        # Database name
DB_PORT=5432                    # Database port

# Authentication
JWT_SECRET=your-secret-key      # JWT signing key (CHANGE IN PRODUCTION!)
JWT_EXPIRY=24h                  # Token expiry duration

# File Uploads
UPLOAD_DIR=./uploads            # Upload directory path
MAX_UPLOAD_SIZE=5242880         # Max file size in bytes (5MB)

# CORS
ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173  # Allowed frontend origins

# Seeding
SEED=true                       # Enable/disable auto-seeding
```

### CORS Configuration

For frontend integration, configure allowed origins in `.env`:

```env
# Allow multiple origins (comma-separated)
ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173,https://yourdomain.com
```

Or edit in code (`internal/middleware/cors.go`):

```go
config := cors.DefaultConfig()
config.AllowOrigins = []string{
    "http://localhost:3000",
    "http://localhost:5173",
    "https://yourdomain.com",
}
```

## 🚢 Deployment

### Production Checklist

Before deploying to production:

1. **Change default credentials**
   ```env
   JWT_SECRET=generate-a-long-random-secret-key-here
   DB_PASSWORD=your-secure-database-password
   ```

2. **Set production environment**
   ```env
   ENV=production
   ```

3. **Configure CORS properly**
   ```env
   ALLOWED_ORIGINS=https://yourdomain.com
   ```

4. **Set up HTTPS** with SSL/TLS certificates

5. **Configure database backups** and monitoring

6. **Enable rate limiting** (consider adding middleware)

7. **Review security headers** and add as needed

8. **Set up logging** and error tracking

### Deployment Options

This API can be deployed to:

- **Docker platforms:** AWS ECS, Google Cloud Run, Azure Container Instances
- **Kubernetes** with manifests or Helm charts
- **Traditional servers** with systemd service
- **PaaS platforms:** Heroku, Railway, Render

### Example Production Dockerfile

The included `Dockerfile` uses multi-stage builds for optimized images:

```dockerfile
# Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o api cmd/api/main.go

# Run stage
FROM alpine:latest
COPY --from=builder /app/api /api
CMD ["/api"]
```

## 🐛 Troubleshooting

### Database Connection Issues

```bash
# Check PostgreSQL is running
docker ps | grep postgres

# Test connection
psql -h localhost -U postgres -d thornton_pickard

# View database logs
docker logs postgres
```

### Port Already in Use

```bash
# Find process using port 8080
lsof -i :8080

# Kill the process or change PORT in .env
```

### Swagger Documentation Not Working

```bash
# Regenerate Swagger docs
swag init -g cmd/api/main.go

# Ensure docs package is imported in main.go
import _ "github.com/Candoo/thornton-pickard-api/docs"
```

### CORS Errors

1. Check `ALLOWED_ORIGINS` in `.env`
2. Verify frontend URL matches exactly (including port)
3. Ensure CORS middleware is registered before routes
4. Check browser console for specific CORS error

## 📦 Technology Stack

- **Language:** Go 1.21
- **Web Framework:** Gin
- **Database:** PostgreSQL 15
- **ORM:** GORM
- **Authentication:** JWT (golang-jwt)
- **Documentation:** Swagger (swaggo)
- **Testing:** testify
- **Configuration:** godotenv
- **CORS:** gin-cors
- **Containerization:** Docker

### Coding Standards

- Follow Go best practices and idioms
- Write unit tests for new features
- Update Swagger documentation
- Format code with `gofmt`
- Run `go vet` before committing

## 📄 License

MIT License - see LICENSE file for details.

## 🔗 Related Links

- [Frontend Repository](https://github.com/Candoo/my-modern-react-setup)
- [Deployment Repository](https://github.com/Candoo/Pickard-Index)
- [Gin Framework](https://gin-gonic.com/)
- [GORM Documentation](https://gorm.io/)
- [Swagger/OpenAPI](https://swagger.io/)

---

**Built with Go 🚀 for the Thornton Pickard camera community**
