# Tracking Diet Backend

A comprehensive Go-based diet tracking backend with AI-powered recommendations, built with clean architecture principles.

## Features

- **User Management**: Registration, authentication, and profile management
- **Body Measurements**: Track weight, BMI, body composition, and body measurements
- **Nutrition Tracking**: Monitor calorie intake, macronutrients, and meal timing
- **Metabolic Health**: Track blood pressure, glucose, cholesterol, and other metabolic markers
- **Fitness Monitoring**: Record workout data, VO2 max, strength metrics, and activity levels
- **Wellbeing Tracking**: Monitor sleep quality, mood, stress levels, and energy
- **Lab Tests**: Store and analyze medical lab results
- **AI Recommendations**: Generate personalized meal plans and health insights using Google Genkit
- **Background Processing**: Scheduled jobs for automated tasks using Redis queue and cron scheduler
- **RESTful API**: Clean JSON API with proper authentication and validation

## Technology Stack

- **Language**: Go 1.27.0
- **Database**: PostgreSQL with GORM ORM
- **Authentication**: JWT with bcrypt password hashing
- **HTTP Server**: Chi router with middleware
- **AI Integration**: Google Genkit with Gemini Flash model
- **Background Processing**: Redis queue with robfig/cron scheduler
- **Validation**: go-playground/validator
- **Testing**: Go testing framework

## Project Structure

```
tracking-diet-backend/
├── cmd/
│   ├── api/              # HTTP API server
│   └── worker/           # Background worker
├── internal/
│   ├── ai/               # AI integration with Genkit
│   ├── auth/             # JWT and password hashing
│   ├── config/           # Configuration management
│   ├── delivery/         # HTTP handlers and DTOs
│   │   ├── dto/          # Data transfer objects
│   │   ├── handlers/     # HTTP request handlers
│   │   └── middleware/   # HTTP middleware
│   ├── domain/           # Domain models and interfaces
│   ├── repository/       # Data access layer
│   ├── usecase/          # Business logic layer
│   └── worker/           # Background job processing
├── migrations/           # Database migrations
├── pkg/                  # Shared packages
│   └── logger/           # Structured logging
├── .env.example          # Environment variables template
├── docker-compose.yml    # Docker services
├── Dockerfile            # API server Docker image
├── Dockerfile.worker     # Worker Docker image
└── go.mod                # Go module definition
```

## Getting Started

### Prerequisites

- Go 1.27.0 or higher
- PostgreSQL 12+
- Redis 6+
- Google AI API key (for AI features)

### Installation

1. Clone the repository:
```bash
git clone https://github.com/irsyadjpp/tracking-diet-backend.git
cd tracking-diet-backend
```

2. Install dependencies:
```bash
go mod download
```

3. Set up environment variables:
```bash
cp .env.example .env
# Edit .env with your configuration
```

4. Run database migrations:
```bash
# Apply migrations
psql -U your_user -d tracking_diet -f migrations/001_init_schema.up.sql
```

### Configuration

The following environment variables are required:

```env
# Application
ENV=development
PORT=8080

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_user
DB_PASSWORD=your_password
DB_NAME=tracking_diet

# JWT Authentication
JWT_SECRET=your-super-secret-jwt-key-minimum-32-characters
JWT_EXPIRATION=24

# Redis (for background jobs)
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# Google AI (for AI recommendations)
GOOGLE_GENKIT_API_KEY=your-google-ai-api-key

# Database Connection Pool
DB_MAX_IDLE_CONNS=10
DB_MAX_OPEN_CONNS=100
DB_CONN_MAX_LIFETIME=1h

# Server Timeouts
READ_TIMEOUT=15s
WRITE_TIMEOUT=15s
IDLE_TIMEOUT=60s
```

### Running the Application

#### Using Docker Compose (Recommended)

```bash
docker-compose up -d
```

This will start:
- PostgreSQL database
- Redis server
- API server on port 8080
- Background worker

#### Running Locally

1. Start PostgreSQL and Redis:
```bash
# Using Docker
docker run -d -p 5432:5432 -e POSTGRES_PASSWORD=your_password postgres:12
docker run -d -p 6379:6379 redis:6
```

2. Run the API server:
```bash
go run cmd/api/main.go
```

3. Run the background worker (in a separate terminal):
```bash
go run cmd/worker/main.go
```

## API Documentation

### Authentication Endpoints

#### Register
```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "full_name": "John Doe",
  "email": "john@example.com",
  "password": "password123",
  "gender": "M",
  "birth_date": "1990-01-01",
  "height_cm": 175.5
}
```

#### Login
```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "password123"
}
```

#### Refresh Token
```http
POST /api/v1/auth/refresh
Content-Type: application/json

{
  "refresh_token": "your_refresh_token"
}
```

#### Get Current User
```http
GET /api/v1/auth/me
Authorization: Bearer your_jwt_token
```

### Body Measurement Endpoints

#### Create Body Measurement
```http
POST /api/v1/measurements/body
Authorization: Bearer your_jwt_token
Content-Type: application/json

{
  "weight_kg": 75.5,
  "waist_cm": 85.0,
  "hip_cm": 95.0,
  "measured_at": "2024-01-15T10:00:00Z"
}
```

#### List Body Measurements
```http
GET /api/v1/measurements/body?start_date=2024-01-01&end_date=2024-01-31
Authorization: Bearer your_jwt_token
```

### Nutrition Endpoints

#### Create Nutrition Entry
```http
POST /api/v1/measurements/nutrition
Authorization: Bearer your_jwt_token
Content-Type: application/json

{
  "calories_in_kcal": 2000,
  "carbs_g": 250.0,
  "protein_g": 150.0,
  "fat_g": 65.0,
  "fiber_g": 25.0,
  "water_intake_l": 2.5,
  "measured_at": "2024-01-15"
}
```

#### Get Daily Summary
```http
GET /api/v1/measurements/nutrition/daily?date=2024-01-15
Authorization: Bearer your_jwt_token
```

## Background Jobs

The system includes several scheduled background jobs:

- **Daily AI Recommendations** (9:00 AM): Generates personalized meal plans for active users
- **Weekly Health Score** (Sunday 8:00 AM): Calculates comprehensive health scores
- **Monthly Data Cleanup** (1st of month 2:00 AM): Cleans up old data based on retention policy
- **Nutrition Reminders** (8:00 AM, 12:00 PM, 6:00 PM): Sends meal timing reminders

## Development

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/auth
```

### Code Style

- Follow Go standard project layout
- Use clean architecture principles
- Add godoc comments for exported functions
- Implement proper error handling with context
- Use structured logging

## Deployment

### Docker Deployment

Build and run using Docker:

```bash
# Build API image
docker build -t tracking-diet-api -f Dockerfile .

# Build worker image
docker build -t tracking-diet-worker -f Dockerfile.worker .

# Run with docker-compose
docker-compose up -d
```

### Environment-Specific Configuration

- **Development**: Set `ENV=development` for detailed logging
- **Production**: Set `ENV=production` for optimized logging and security

## Security Considerations

- JWT secrets must be at least 32 characters
- Passwords are hashed with bcrypt (cost factor 12)
- SQL injection prevention via GORM parameterization
- Rate limiting should be configured on authentication endpoints
- Input validation on all public endpoints
- CORS configuration for frontend integration

## Monitoring

The application includes:
- Structured logging with request ID tracking
- Health check endpoint at `/health`
- Error stack trace logging in development mode
- Worker pool statistics endpoint

## License

This project is licensed under the MIT License.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## Support

For support, please open an issue in the GitHub repository.