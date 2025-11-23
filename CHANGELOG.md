# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2025-01-15

### Added
- Initial release of Community Issue Mapper
- Interactive map interface using Leaflet.js
- Issue reporting with geolocation
- Multiple issue categories (pothole, accessibility, streetlight, graffiti, trash, other)
- Image upload functionality (max 10MB)
- Community voting system
- Real-time statistics dashboard
- Advanced filtering by category, status, and location
- RESTful API with comprehensive endpoints
- PostgreSQL database with optimized indexes
- Docker and docker-compose support
- Database migrations
- CORS middleware with configurable origins
- Rate limiting middleware
- Request/response logging
- Health check and readiness endpoints
- Clean architecture pattern
- Comprehensive API documentation
- Deployment guides for AWS, GCP, DigitalOcean, and Heroku
- CI/CD pipeline with GitHub Actions
- Unit tests and code coverage
- Responsive mobile design
- Auto-location detection
- Production-ready security features

### Security
- Input validation and sanitization
- SQL injection prevention (parameterized queries)
- File upload validation
- Rate limiting protection
- CORS protection

## [Unreleased]

### Planned
- User authentication and authorization
- Email notifications for issue updates
- Admin dashboard for issue management
- Mobile applications (iOS/Android)
- SMS notifications
- Integration with city management systems
- Analytics and reporting dashboard
- Multi-language support
- Dark mode
- Offline support with PWA
- WebSocket support for real-time updates
- Issue comments and discussions
- Issue assignment to city officials
- Automated issue status updates
- Export data to CSV/PDF
- Public API with API keys
- Webhooks for integrations
- Issue clustering on map
- Heatmap visualization
- Time-series analysis
- Bulk issue operations
- Advanced search with Elasticsearch
- Redis caching layer
- S3 integration for file storage
- Email verification
- Password reset functionality
- Social login (Google, Facebook)
- Two-factor authentication

---

For detailed changes in each release, see the [GitHub releases page](https://github.com/codeforgood-org/community-issue-mapper/releases).
