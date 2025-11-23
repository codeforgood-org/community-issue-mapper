# 🗺️ Community Issue Mapper

A production-ready crowdsourced issue reporting platform built with Leaflet.js and Go. Report and track community issues like potholes, accessibility problems, broken streetlights, and more.

![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)
![Docker](https://img.shields.io/badge/Docker-ready-2496ED?logo=docker)

## ✨ Features

### Core Functionality
- 🗺️ **Interactive Map**: Click-to-report issues with Leaflet.js integration
- 📸 **Image Upload**: Attach photos to issue reports (max 10MB)
- 🏷️ **Categorization**: Multiple issue categories (pothole, accessibility, streetlight, graffiti, trash, other)
- 📊 **Real-time Statistics**: Track total, new, and resolved issues
- 👍 **Community Voting**: Upvote important issues to increase visibility
- 🔍 **Advanced Filtering**: Filter by category, status, and location
- 📍 **Geolocation**: Automatic user location detection
- 📱 **Responsive Design**: Works seamlessly on desktop and mobile

### Technical Features
- ⚡ **RESTful API**: Clean, well-documented API endpoints
- 🐘 **PostgreSQL Database**: Robust data persistence with proper indexing
- 🐳 **Docker Ready**: Full containerization with docker-compose
- 🔒 **Production Security**: CORS, rate limiting, input validation
- 📈 **Scalable Architecture**: Clean architecture pattern with repository layer
- 🚀 **CI/CD Pipeline**: Automated testing and deployment with GitHub Actions
- 📝 **Comprehensive Logging**: Request/response logging for debugging
- 🏥 **Health Checks**: Built-in health and readiness endpoints

## 🏗️ Architecture

```
community-issue-mapper/
├── backend/                    # Go backend application
│   ├── cmd/
│   │   └── server/            # Application entry point
│   ├── internal/
│   │   ├── api/               # API layer
│   │   │   ├── handlers/      # HTTP handlers
│   │   │   ├── middleware/    # Middleware (CORS, logging, rate limiting)
│   │   │   └── router.go      # Route definitions
│   │   ├── config/            # Configuration management
│   │   ├── database/          # Database connection
│   │   ├── models/            # Data models
│   │   ├── repository/        # Data access layer
│   │   └── service/           # Business logic layer
│   ├── migrations/            # Database migrations
│   └── go.mod                 # Go dependencies
├── frontend/                  # Frontend application
│   ├── css/                   # Stylesheets
│   ├── js/                    # JavaScript modules
│   │   ├── api.js            # API client
│   │   ├── map.js            # Map management
│   │   └── app.js            # Main application logic
│   └── index.html            # Main HTML file
├── uploads/                   # Uploaded images storage
├── .github/workflows/         # GitHub Actions CI/CD
├── docker-compose.yml         # Docker composition
├── Dockerfile                 # Docker image definition
└── README.md                 # This file
```

## 🚀 Quick Start

### Prerequisites
- Docker and Docker Compose (recommended)
- OR Go 1.21+ and PostgreSQL 14+

### Option 1: Docker (Recommended)

1. **Clone the repository**
   ```bash
   git clone https://github.com/codeforgood-org/community-issue-mapper.git
   cd community-issue-mapper
   ```

2. **Create environment file**
   ```bash
   cp .env.example .env
   # Edit .env with your preferred settings
   ```

3. **Start the application**
   ```bash
   docker-compose up -d
   ```

4. **Run migrations**
   ```bash
   docker-compose run migrate
   ```

5. **Access the application**
   - Web UI: http://localhost:8080
   - API: http://localhost:8080/api/v1
   - Health Check: http://localhost:8080/health

### Option 2: Local Development

1. **Install dependencies**
   ```bash
   cd backend
   go mod download
   ```

2. **Set up PostgreSQL**
   ```bash
   createdb community_issues
   ```

3. **Run migrations**
   ```bash
   # Install migrate tool first
   brew install golang-migrate  # macOS
   # or download from https://github.com/golang-migrate/migrate

   make migrate-up
   ```

4. **Start the server**
   ```bash
   make run
   # or with auto-reload (requires air)
   make dev
   ```

## 📚 API Documentation

### Base URL
```
http://localhost:8080/api/v1
```

### Endpoints

#### Issues

##### List Issues
```http
GET /api/v1/issues
```

Query Parameters:
- `category` (optional): Filter by category (pothole, accessibility, streetlight, graffiti, trash, other)
- `status` (optional): Filter by status (new, in_progress, resolved, closed)
- `min_lat`, `max_lat`, `min_lng`, `max_lng` (optional): Bounding box filter
- `limit` (optional): Maximum number of results (default: 100, max: 1000)
- `offset` (optional): Pagination offset

Response:
```json
[
  {
    "id": 1,
    "title": "Large pothole on Main Street",
    "description": "Dangerous pothole near the intersection",
    "category": "pothole",
    "status": "new",
    "latitude": 37.7749,
    "longitude": -122.4194,
    "address": "123 Main St",
    "image_url": "/uploads/1.jpg",
    "reporter_name": "John Doe",
    "reporter_email": "john@example.com",
    "votes": 5,
    "created_at": "2025-01-01T10:00:00Z",
    "updated_at": "2025-01-01T10:00:00Z"
  }
]
```

##### Create Issue
```http
POST /api/v1/issues
Content-Type: application/json
```

Request Body:
```json
{
  "title": "Broken streetlight",
  "description": "Streetlight not working for 2 weeks",
  "category": "streetlight",
  "latitude": 37.7749,
  "longitude": -122.4194,
  "address": "456 Oak Ave",
  "reporter_name": "Jane Smith",
  "reporter_email": "jane@example.com"
}
```

##### Get Issue
```http
GET /api/v1/issues/{id}
```

##### Update Issue
```http
PATCH /api/v1/issues/{id}
Content-Type: application/json
```

Request Body (all fields optional):
```json
{
  "title": "Updated title",
  "description": "Updated description",
  "status": "in_progress",
  "category": "pothole"
}
```

##### Delete Issue
```http
DELETE /api/v1/issues/{id}
```

##### Upload Image
```http
POST /api/v1/issues/{id}/upload
Content-Type: multipart/form-data
```

Form Data:
- `image`: Image file (max 10MB)

##### Vote for Issue
```http
POST /api/v1/issues/{id}/vote
```

#### Statistics

##### Get Statistics
```http
GET /api/v1/stats
```

Response:
```json
{
  "total_issues": 150,
  "issues_by_status": {
    "new": 45,
    "in_progress": 30,
    "resolved": 60,
    "closed": 15
  },
  "issues_by_category": {
    "pothole": 40,
    "accessibility": 25,
    "streetlight": 20,
    "graffiti": 15,
    "trash": 30,
    "other": 20
  }
}
```

#### Health Checks

##### Health Check
```http
GET /health
```

##### Readiness Check
```http
GET /ready
```

## 🛠️ Development

### Available Make Commands

```bash
make help              # Show all available commands
make build             # Build the application
make run               # Run the application
make test              # Run tests
make test-coverage     # Run tests with coverage report
make clean             # Clean build artifacts
make docker-build      # Build Docker image
make docker-up         # Start Docker containers
make docker-down       # Stop Docker containers
make docker-logs       # View Docker logs
make migrate-up        # Run database migrations up
make migrate-down      # Run database migrations down
make migrate-create    # Create a new migration
make deps              # Download dependencies
make lint              # Run linter
make dev               # Run with auto-reload
```

### Running Tests

```bash
cd backend
go test -v ./...

# With coverage
go test -v -race -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Creating Database Migrations

```bash
make migrate-create name=add_new_field
```

This creates two files in `backend/migrations/`:
- `{timestamp}_add_new_field.up.sql`
- `{timestamp}_add_new_field.down.sql`

## 🌍 Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `8080` |
| `ENV` | Environment (development/production) | `development` |
| `DB_HOST` | Database host | `localhost` |
| `DB_PORT` | Database port | `5432` |
| `DB_USER` | Database user | `issueapp` |
| `DB_PASSWORD` | Database password | `issueapp_password` |
| `DB_NAME` | Database name | `community_issues` |
| `DB_SSL_MODE` | Database SSL mode | `disable` |
| `CORS_ALLOWED_ORIGINS` | Allowed CORS origins | `*` |
| `MAX_UPLOAD_SIZE` | Max upload size in bytes | `10485760` (10MB) |
| `UPLOAD_DIR` | Upload directory | `./uploads` |
| `RATE_LIMIT_REQUESTS` | Rate limit requests | `100` |
| `RATE_LIMIT_DURATION` | Rate limit duration | `1m` |

## 📦 Deployment

### Docker Deployment

1. **Build and push Docker image**
   ```bash
   docker build -t your-registry/community-issue-mapper:latest .
   docker push your-registry/community-issue-mapper:latest
   ```

2. **Deploy with docker-compose**
   ```bash
   docker-compose up -d
   ```

### Cloud Deployment

The application can be easily deployed to:
- **AWS ECS/Fargate**: Use the included Dockerfile
- **Google Cloud Run**: Supports container deployment
- **Azure Container Instances**: Direct Docker deployment
- **Heroku**: Use container registry
- **DigitalOcean App Platform**: Docker-based deployment

### Environment-specific configurations

- Set `ENV=production` for production
- Use proper database credentials
- Configure CORS origins appropriately
- Set up SSL/TLS termination at load balancer
- Use managed PostgreSQL service for production

## 🔒 Security

- ✅ CORS protection with configurable origins
- ✅ Rate limiting to prevent abuse
- ✅ Input validation and sanitization
- ✅ SQL injection prevention (parameterized queries)
- ✅ File upload validation
- ✅ Secure headers
- ⚠️ **TODO**: Add authentication/authorization
- ⚠️ **TODO**: Add HTTPS enforcement
- ⚠️ **TODO**: Add API key management

## 🧪 Testing

The project includes:
- Unit tests for repository layer
- Integration tests for API endpoints
- GitHub Actions CI/CD pipeline
- Code coverage reporting

## 📊 Monitoring

### Health Checks

- `/health`: Database connectivity check
- `/ready`: Application readiness check

### Logging

All requests are logged with:
- HTTP method
- Request URI
- Response status
- Response size
- Request duration

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Leaflet.js](https://leafletjs.com/) - Interactive map library
- [OpenStreetMap](https://www.openstreetmap.org/) - Map data
- [Gorilla Mux](https://github.com/gorilla/mux) - HTTP router
- [PostgreSQL](https://www.postgresql.org/) - Database

## 📧 Support

For support, please open an issue on GitHub or contact the maintainers.

## 🗺️ Roadmap

- [ ] User authentication and authorization
- [ ] Email notifications for issue updates
- [ ] Admin dashboard for issue management
- [ ] Mobile applications (iOS/Android)
- [ ] SMS notifications
- [ ] Integration with city management systems
- [ ] Analytics and reporting dashboard
- [ ] Multi-language support
- [ ] Dark mode
- [ ] Offline support with PWA

---

Built with ❤️ for communities everywhere
