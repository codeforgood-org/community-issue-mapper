# Platform Enhancements

This document describes the major enhancements added to the Community Issue Mapper platform.

## Table of Contents

1. [Authentication System](#authentication-system)
2. [Admin Dashboard](#admin-dashboard)
3. [Real-time Updates](#real-time-updates)
4. [Comments & Discussions](#comments--discussions)
5. [Advanced Analytics](#advanced-analytics)
6. [Progressive Web App (PWA)](#progressive-web-app-pwa)
7. [Map Clustering](#map-clustering)
8. [Kubernetes Deployment](#kubernetes-deployment)
9. [Monitoring & Observability](#monitoring--observability)

## Authentication System

### Features

- **JWT-based authentication** with secure token handling
- **User registration and login** with email/password
- **Password hashing** using bcrypt (cost factor: 12)
- **Role-based access control (RBAC)**
  - User: Standard user with basic permissions
  - Moderator: Can manage issues and moderate content
  - Admin: Full system access
- **User profiles** with customizable avatars and bios
- **Token refresh** capability (valid for 7 days after expiry)

### API Endpoints

```
POST /api/v1/auth/register   - Register new user
POST /api/v1/auth/login      - Login and get JWT token
GET  /api/v1/auth/profile    - Get user profile (authenticated)
PATCH /api/v1/auth/profile   - Update user profile (authenticated)
```

### Database Schema

```sql
users table:
- id (serial primary key)
- email (unique, indexed)
- name
- password_hash
- role (user, moderator, admin)
- avatar
- bio
- verified (boolean)
- active (boolean)
- created_at, updated_at
```

### Security Features

- Passwords hashed with bcrypt (12 rounds)
- JWT tokens with expiration
- Role-based middleware for endpoint protection
- Email validation with regex
- Password minimum length: 8 characters

## Admin Dashboard

### Features

- **Comprehensive analytics dashboard**
  - Total issues, users, comments
  - Issues this week/month
  - Resolved issues statistics
  - Average resolution time
- **Issues management table**
  - Search and filter capabilities
  - Status updates
  - Bulk operations
- **Real-time charts** using Chart.js
  - Issues trend (last 30 days)
  - Issues by category (doughnut chart)
  - Category distribution
- **Activity timeline** showing recent actions
- **Top contributors** leaderboard

### Dashboard Sections

1. **Dashboard**: Overview with stats and charts
2. **Manage Issues**: Full CRUD operations on issues
3. **Users**: User management (placeholder)
4. **Analytics**: Advanced analytics and insights
5. **Settings**: System settings (placeholder)

### Access

Navigate to `/admin.html` (requires admin role)

### Technologies

- Chart.js for data visualization
- Responsive grid layout
- Real-time data updates
- Modern UI with gradient sidebar

## Real-time Updates

### WebSocket Implementation

- **Bi-directional communication** between server and clients
- **Hub pattern** for managing connections
- **Automatic reconnection** on connection loss
- **Ping/pong heartbeat** for connection health

### Features

- Live issue updates
- Real-time statistics
- Instant notifications
- Concurrent user tracking

### Connection

```javascript
const ws = new WebSocket('ws://localhost:8080/ws');

ws.onmessage = (event) => {
    const message = JSON.parse(event.data);
    handleRealtimeUpdate(message);
};
```

### Message Types

- `issue_created`: New issue reported
- `issue_updated`: Issue status changed
- `issue_voted`: New vote on issue
- `comment_added`: New comment posted
- `stats_updated`: Statistics changed

## Comments & Discussions

### Features

- **Threaded comments** on issues
- **User attribution** with names and avatars
- **Comment timestamps**
- **Delete own comments**
- **Activity tracking** for comments

### Database Schema

```sql
comments table:
- id
- issue_id (foreign key)
- user_id (foreign key)
- content (text, max 2000 chars)
- created_at, updated_at

Indexes:
- issue_id (for fast retrieval)
- user_id (for user activity)
```

### API Endpoints

```
GET  /api/v1/issues/{id}/comments    - Get all comments for an issue
POST /api/v1/issues/{id}/comments    - Add comment (authenticated)
DELETE /api/v1/comments/{id}          - Delete comment (owner only)
```

## Advanced Analytics

### Metrics

- **Total counts**: Issues, users, comments
- **Time-based metrics**:
  - Issues this week/month
  - Resolved issues this week/month
- **Performance metrics**:
  - Average resolution time (hours)
  - Resolution rate percentage
- **Distribution analytics**:
  - Issues by category
  - Issues by status
- **Trend data**: 30-day time series
- **Top contributors**: Users with most activity

### API Endpoint

```
GET /api/v1/analytics
```

### Response Example

```json
{
  "total_issues": 150,
  "total_users": 45,
  "total_comments": 230,
  "issues_this_week": 12,
  "issues_this_month": 48,
  "resolved_this_week": 8,
  "resolved_this_month": 35,
  "avg_resolution_time_hours": 24.5,
  "issues_by_category": {...},
  "issues_by_status": {...},
  "issues_trend": [...],
  "top_contributors": [...]
}
```

## Progressive Web App (PWA)

### Features

- **Installable** on mobile and desktop
- **Offline support** with service worker
- **Caching strategies**:
  - Cache-first for static assets
  - Network-first for API requests
- **Background sync** for offline submissions
- **Push notifications** (infrastructure ready)
- **App-like experience** with manifest

### Service Worker

Location: `/sw.js`

Caching strategies:
- Precache: HTML, CSS, JS, manifest
- Runtime cache: API responses
- Network-first: API endpoints
- Cache-first: Static assets

### Manifest

Location: `/manifest.json`

Features:
- App name and short name
- Theme colors
- Icons (72x72 to 512x512)
- Display mode: standalone
- Categories: civic, social, utilities

### Installation

Users can install the app:
- **Desktop**: Chrome/Edge - Install button in address bar
- **Mobile**: Add to Home Screen option

## Map Clustering

### Features

- **Marker clustering** using Leaflet.markercluster
- **Automatic grouping** of nearby markers
- **Spiderfy on zoom** for clustered markers
- **Coverage display** on hover
- **Performance optimization** for large datasets
- **Toggleable** clustering on/off

### Configuration

```javascript
markerCluster = L.markerClusterGroup({
    maxClusterRadius: 50,
    spiderfyOnMaxZoom: true,
    showCoverageOnHover: true,
    zoomToBoundsOnClick: true
});
```

### Benefits

- Better performance with 1000+ markers
- Cleaner map visualization
- Improved user experience
- Reduced visual clutter

### Export

GeoJSON export functionality included:
```javascript
const geoJSON = issueMap.exportToGeoJSON();
```

## Kubernetes Deployment

### Components

1. **Application Deployment**
   - 3 replicas (min)
   - Auto-scaling (up to 10 pods)
   - Resource limits and requests
   - Health checks (liveness, readiness)
   - Rolling updates

2. **PostgreSQL StatefulSet**
   - Persistent storage
   - Health probes
   - Resource management

3. **Ingress**
   - NGINX ingress controller
   - SSL/TLS with cert-manager
   - Rate limiting
   - Path-based routing

4. **ConfigMaps & Secrets**
   - Environment configuration
   - Sensitive data management
   - Database credentials

5. **HorizontalPodAutoscaler**
   - CPU-based scaling (70% threshold)
   - Memory-based scaling (80% threshold)
   - Min: 3, Max: 10 replicas

6. **PodDisruptionBudget**
   - Minimum 2 pods available
   - Ensures high availability

### Deployment

```bash
# Apply all Kubernetes resources
kubectl apply -f k8s/

# Check deployment status
kubectl get pods
kubectl get svc
kubectl get ingress

# Scale manually
kubectl scale deployment community-issue-mapper --replicas=5
```

### Monitoring

```bash
# View logs
kubectl logs -f deployment/community-issue-mapper

# Check resource usage
kubectl top pods
kubectl top nodes
```

## Monitoring & Observability

### Prometheus

**Metrics collected:**
- HTTP request rate and latency
- Error rates (4xx, 5xx)
- Resource usage (CPU, memory)
- Database connections
- Custom application metrics

**Alerting rules:**
- High error rate (>5% for 5 min)
- High response time (p95 > 1s)
- Pod down (2 min)
- High memory usage (>90%)
- Database down

### Grafana Dashboard

**Panels:**
1. Total Issues (stat)
2. API Request Rate (graph)
3. Response Time p95 (graph)
4. Error Rate (graph)
5. Memory Usage (graph)
6. CPU Usage (graph)
7. Active Pods (stat)
8. Database Connections (graph)

### Setup

```bash
# Deploy Prometheus
kubectl apply -f monitoring/prometheus-config.yaml

# Access Prometheus
kubectl port-forward svc/prometheus 9090:9090

# Import Grafana dashboard
# Use monitoring/grafana-dashboard.json
```

### Custom Metrics

Add application-specific metrics:
```go
// Example: Track issue creation
issuesCreated := prometheus.NewCounter(
    prometheus.CounterOpts{
        Name: "issues_created_total",
        Help: "Total number of issues created",
    },
)
```

## Additional Features

### Activity Tracking

All user actions tracked:
- Issue creation
- Issue updates
- Comments
- Votes
- Status changes

Stored in `activities` table for analytics.

### Notifications (Infrastructure)

Database schema ready for notifications:
```sql
notifications table:
- id
- user_id
- title
- message
- type (issue_update, comment, mention)
- issue_id
- read (boolean)
- created_at
```

Future integration with:
- Email notifications
- Push notifications
- SMS notifications

### Vote Tracking

Prevents duplicate voting:
```sql
issue_votes table:
- id
- issue_id
- user_id
- created_at
- UNIQUE(issue_id, user_id)
```

## Performance Optimizations

1. **Database Indexes**
   - All foreign keys indexed
   - Common query fields indexed
   - Composite indexes for complex queries

2. **Caching Strategy**
   - Service worker caching
   - Browser caching headers
   - API response caching (ready for Redis)

3. **Map Performance**
   - Marker clustering
   - Lazy loading
   - Viewport-based queries

4. **Bundle Optimization**
   - CDN for libraries
   - Minified assets
   - Gzip compression

## Security Enhancements

1. **Authentication**
   - JWT with expiration
   - Secure password hashing
   - Role-based access control

2. **Input Validation**
   - Email format validation
   - Content length limits
   - SQL injection prevention

3. **Rate Limiting**
   - Per-IP rate limiting
   - Configurable thresholds
   - DDoS protection

4. **CORS**
   - Configurable origins
   - Secure defaults

5. **SSL/TLS**
   - HTTPS enforcement
   - Certificate management

## Future Enhancements

- [ ] Email notification system
- [ ] Redis caching layer
- [ ] Advanced search with Elasticsearch
- [ ] Export to CSV/PDF
- [ ] Issue templates
- [ ] Bulk operations
- [ ] Heatmap visualization
- [ ] Mobile apps (iOS/Android)
- [ ] Social login (OAuth)
- [ ] Two-factor authentication
- [ ] Webhook integrations
- [ ] API rate limiting per user
- [ ] Image optimization
- [ ] Video attachments
- [ ] Multi-language support
- [ ] Dark mode
- [ ] Accessibility improvements

## Migration Guide

### From v1.0.0 to v2.0.0

1. **Run new migrations**
   ```bash
   make migrate-up
   ```

2. **Update environment variables**
   Add JWT_SECRET and JWT_EXPIRY

3. **Deploy new version**
   ```bash
   docker-compose down
   docker-compose up -d --build
   ```

4. **Verify deployment**
   ```bash
   curl http://localhost:8080/health
   ```

## Support

For issues or questions:
- GitHub Issues: https://github.com/codeforgood-org/community-issue-mapper/issues
- Documentation: See `/docs` directory

---

Built with ❤️ for communities everywhere
