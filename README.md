# Thornton Pickard Camera API

Go-based REST API for Thornton Pickard camera and ephemera data with authentication, pagination, search, and image uploads.

## Features

- 🚀 Fast and efficient Go backend with Gin
- 🔐 JWT authentication with user roles
- 🗄️ PostgreSQL database with GORM ORM
- 📚 Auto-generated Swagger documentation
- 🐳 Docker support with Docker Compose
- 🔍 Advanced search and filtering
- 📄 Pagination support
- 📸 Image upload functionality
- ✅ Unit tests with testify
- 🌱 Database seeding
- 🛡️ CORS middleware
- 📊 RESTful API design

## Getting Started

### Prerequisites

- Go 1.21+
- PostgreSQL 15+
- Docker & Docker Compose (optional)

### Local Development

1. **Install PostgreSQL** or run with Docker:
```bash
docker run -d \
  --name postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=thornton_pickard \
  -p 5432:5432 \
  postgres:15-alpine
```

2. **Clone and setup:**
```bash
# Clone repository
git clone <your-repo-url>
cd thornton-pickard-api

# Install dependencies
go mod download

# Copy and configure environment
cp .env.example .env
# Edit .env with your settings

# Generate Swagger docs
swag init -g cmd/api/main.go

# Run with database seeding
SEED=true go run cmd/api/main.go
```

3. **Access the API:**
   - API: http://localhost:8080
   - Swagger Docs: http://localhost:8080/swagger/index.html
   - Health Check: http://localhost:8080/health

### Docker Deployment
```bash
# Start everything (PostgreSQL + API)
docker-compose up -d

# View logs
docker-compose logs -f api

# Stop
docker-compose down

# Rebuild after changes
docker-compose up -d --build
```

## API Endpoints

### Authentication
- `POST /api/v1/auth/register` - Register new user
- `POST /api/v1/auth/login` - Login and get JWT token
- `GET /api/v1/auth/profile` - Get user profile (requires auth)

### Cameras
- `GET /api/v1/cameras` - List all cameras (with pagination & search)
- `GET /api/v1/cameras/:id` - Get camera by ID
- `POST /api/v1/cameras` - Create camera (requires auth)
- `PUT /api/v1/cameras/:id` - Update camera (requires auth)
- `DELETE /api/v1/cameras/:id` - Delete camera (admin only)
- `GET /api/v1/cameras/search` - Search cameras (deprecated, use GET /cameras with query params)

### Ephemera
- `GET /api/v1/ephemera` - List all ephemera
- `GET /api/v1/ephemera/:id` - Get ephemera by ID
- `POST /api/v1/ephemera` - Create ephemera (requires auth)

### Manufacturers
- `GET /api/v1/manufacturers` - List manufacturers

### Uploads
- `POST /api/v1/upload` - Upload single image (requires auth)
- `POST /api/v1/upload/multiple` - Upload multiple images (requires auth)

## Query Parameters

### Pagination
- `page` - Page number (default: 1)
- `page_size` - Items per page (default: 10, max: 100)

### Search & Filters
- `search` - Full-text search across name, manufacturer, description
- `manufacturer` - Filter by manufacturer
- `year_from` - Filter by year from
- `year_to` - Filter by year to
- `format` - Filter by camera format
- `sort` - Sort by field (name, year_introduced)
- `order` - Sort order (asc, desc)

### Example Queries
```bash
# Get page 2 with 20 items
GET /api/v1/cameras?page=2&page_size=20

# Search for "ruby" cameras
GET /api/v1/cameras?search=ruby

# Filter by manufacturer and year range
GET /api/v1/cameras?manufacturer=Thornton-Pickard&year_from=1900&year_to=1920

# Sort by year descending
GET /api/v1/cameras?sort=year_introduced&order=desc

# Combined filters
GET /api/v1/cameras?search=reflex&format=Plate&year_from=1900&sort=year_introduced
```

## Authentication

The API uses JWT (JSON Web Tokens) for authentication.

### 1. Register or Login
```bash
# Register
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'

# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'
```

### 2. Use Token in Requests
```bash
curl -X POST http://localhost:8080/api/v1/cameras \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"New Camera","manufacturer":"Thornton-Pickard"}'
```

### Default Admin Account

After seeding, use these credentials:
- Email: `admin@thorntonpickard.com`
- Password: `admin123`

**⚠️ Change this password in production!**

## Testing
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

## Database Seeding

The API automatically seeds the database on startup if `SEED=true` in `.env`.

Seed data includes:
- Admin user account
- Sample manufacturers
- Sample cameras (from `seeds/cameras.json` if available)

### Custom Seed Data

Create `seeds/cameras.json` with your camera data:
```json
[
  {
    "name": "Camera Name",
    "manufacturer": "Manufacturer",
    "year_introduced": 1900,
    "format": "Plate"
  }
]
```

## Image Uploads

### Upload Single Image
```bash
curl -X POST http://localhost:8080/api/v1/upload \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "file=@/path/to/image.jpg"
```

Response:
```json
{
  "url": "/uploads/uuid_timestamp.jpg",
  "filename": "image.jpg",
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
```

### Restrictions
- Allowed formats: JPG, JPEG, PNG, GIF, WebP
- Max file size: 5MB per file
- Max files: 10 per request (multiple upload)

## Project Structure
```
thornton-pickard-api/
├── cmd/
│   └── api/
│       └── main.go              # Application entry point
├── internal/
│   ├── models/                  # Database models
│   │   ├── camera.go
│   │   ├── ephemera.go
│   │   ├── manufacturer.go
│   │   └── user.go
│   ├── handlers/                # HTTP request handlers
│   │   ├── camera.go
│   │   ├── ephemera.go
│   │   ├── auth.go
│   │   └── upload.go
│   ├── middleware/              # Middleware (CORS, auth)
│   │   ├── cors.go
│   │   ├── auth.go
│   │   └── pagination.go
│   ├── services/                # Business logic
│   │   ├── auth.go
│   │   └── storage.go
│   ├── database/                # Database connection & seeding
│   │   ├── database.go
│   │   └── seed.go
│   └── utils/                   # Utilities
│       └── pagination.go
├── tests/                       # Unit tests
│   ├── camera_test.go
│   └── auth_test.go
├── docs/                        # Swagger docs (auto-generated)
├── seeds/                       # Seed data files
│   └── cameras.json
├── uploads/                     # Uploaded files
├── .env                         # Environment variables
├── Dockerfile                   # Docker image definition
├── docker-compose.yml           # Docker Compose configuration
├── go.mod                       # Go dependencies
└── README.md                    # This file
```

## Environment Variables
```env
ENV=development                  # development or production
PORT=8080                       # Server port

# Database
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=thornton_pickard
DB_PORT=5432

# Authentication
JWT_SECRET=your-secret-key      # Change in production!

# Uploads
UPLOAD_DIR=./uploads

# Seeding
SEED=true                       # Set to true to seed on startup
```

## Tech Stack

- **Language:** Go 1.21
- **Framework:** Gin (HTTP framework)
- **Database:** PostgreSQL 15
- **ORM:** GORM
- **