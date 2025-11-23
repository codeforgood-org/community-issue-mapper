# API Documentation

## Overview

The Community Issue Mapper API is a RESTful API that allows you to manage community issues, upload images, and track statistics.

## Base URL

```
http://localhost:8080/api/v1
```

## Authentication

Currently, the API does not require authentication. Future versions will include JWT-based authentication.

## Rate Limiting

The API implements rate limiting to prevent abuse:
- Default: 100 requests per minute per IP address
- Configurable via environment variables

When rate limit is exceeded:
```http
HTTP/1.1 429 Too Many Requests
Content-Type: application/json

{
  "error": "Rate limit exceeded"
}
```

## Common Response Codes

- `200 OK`: Successful request
- `201 Created`: Resource successfully created
- `204 No Content`: Successful request with no response body
- `400 Bad Request`: Invalid request parameters
- `404 Not Found`: Resource not found
- `429 Too Many Requests`: Rate limit exceeded
- `500 Internal Server Error`: Server error

## Error Response Format

```json
{
  "error": "Error message describing what went wrong"
}
```

## Endpoints

### Issues

#### List Issues

Retrieve a list of issues with optional filtering.

```http
GET /api/v1/issues
```

**Query Parameters:**

| Parameter | Type | Description | Example |
|-----------|------|-------------|---------|
| `category` | string | Filter by category | `pothole` |
| `status` | string | Filter by status | `new` |
| `min_lat` | float | Minimum latitude for bounding box | `37.7` |
| `max_lat` | float | Maximum latitude for bounding box | `37.8` |
| `min_lng` | float | Minimum longitude for bounding box | `-122.5` |
| `max_lng` | float | Maximum longitude for bounding box | `-122.4` |
| `limit` | integer | Maximum results (1-1000) | `50` |
| `offset` | integer | Pagination offset | `0` |

**Categories:**
- `pothole`
- `accessibility`
- `streetlight`
- `graffiti`
- `trash`
- `other`

**Statuses:**
- `new`
- `in_progress`
- `resolved`
- `closed`

**Example Request:**
```bash
curl "http://localhost:8080/api/v1/issues?category=pothole&status=new&limit=10"
```

**Example Response:**
```json
[
  {
    "id": 1,
    "title": "Large pothole on Main Street",
    "description": "Dangerous pothole near the intersection with 5th Ave",
    "category": "pothole",
    "status": "new",
    "latitude": 37.7749,
    "longitude": -122.4194,
    "address": "123 Main St, San Francisco, CA",
    "image_url": "/uploads/1.jpg",
    "reporter_name": "John Doe",
    "reporter_email": "john@example.com",
    "votes": 5,
    "created_at": "2025-01-15T10:30:00Z",
    "updated_at": "2025-01-15T10:30:00Z"
  }
]
```

#### Create Issue

Create a new issue report.

```http
POST /api/v1/issues
Content-Type: application/json
```

**Request Body:**
```json
{
  "title": "Broken streetlight",
  "description": "Streetlight has been out for 2 weeks",
  "category": "streetlight",
  "latitude": 37.7749,
  "longitude": -122.4194,
  "address": "456 Oak Ave, San Francisco, CA",
  "reporter_name": "Jane Smith",
  "reporter_email": "jane@example.com"
}
```

**Required Fields:**
- `title` (string, max 255 chars)
- `description` (string)
- `category` (enum: pothole, accessibility, streetlight, graffiti, trash, other)
- `latitude` (float, -90 to 90)
- `longitude` (float, -180 to 180)

**Optional Fields:**
- `address` (string, max 500 chars)
- `reporter_name` (string, max 255 chars)
- `reporter_email` (string, max 255 chars)

**Example Response:**
```json
{
  "id": 2,
  "title": "Broken streetlight",
  "description": "Streetlight has been out for 2 weeks",
  "category": "streetlight",
  "status": "new",
  "latitude": 37.7749,
  "longitude": -122.4194,
  "address": "456 Oak Ave, San Francisco, CA",
  "image_url": "",
  "reporter_name": "Jane Smith",
  "reporter_email": "jane@example.com",
  "votes": 0,
  "created_at": "2025-01-15T11:00:00Z",
  "updated_at": "2025-01-15T11:00:00Z"
}
```

#### Get Issue

Retrieve a specific issue by ID.

```http
GET /api/v1/issues/{id}
```

**Path Parameters:**
- `id` (integer): Issue ID

**Example Request:**
```bash
curl http://localhost:8080/api/v1/issues/1
```

**Example Response:**
```json
{
  "id": 1,
  "title": "Large pothole on Main Street",
  "description": "Dangerous pothole near the intersection with 5th Ave",
  "category": "pothole",
  "status": "new",
  "latitude": 37.7749,
  "longitude": -122.4194,
  "address": "123 Main St, San Francisco, CA",
  "image_url": "/uploads/1.jpg",
  "reporter_name": "John Doe",
  "reporter_email": "john@example.com",
  "votes": 5,
  "created_at": "2025-01-15T10:30:00Z",
  "updated_at": "2025-01-15T10:30:00Z"
}
```

#### Update Issue

Update an existing issue.

```http
PATCH /api/v1/issues/{id}
Content-Type: application/json
```

**Path Parameters:**
- `id` (integer): Issue ID

**Request Body (all fields optional):**
```json
{
  "title": "Updated title",
  "description": "Updated description",
  "category": "pothole",
  "status": "in_progress",
  "address": "Updated address"
}
```

**Example Response:**
```json
{
  "id": 1,
  "title": "Updated title",
  "description": "Updated description",
  "category": "pothole",
  "status": "in_progress",
  "latitude": 37.7749,
  "longitude": -122.4194,
  "address": "Updated address",
  "image_url": "/uploads/1.jpg",
  "reporter_name": "John Doe",
  "reporter_email": "john@example.com",
  "votes": 5,
  "created_at": "2025-01-15T10:30:00Z",
  "updated_at": "2025-01-15T12:00:00Z"
}
```

#### Delete Issue

Delete an issue.

```http
DELETE /api/v1/issues/{id}
```

**Path Parameters:**
- `id` (integer): Issue ID

**Example Request:**
```bash
curl -X DELETE http://localhost:8080/api/v1/issues/1
```

**Response:**
```
HTTP/1.1 204 No Content
```

#### Upload Image

Upload an image for an issue.

```http
POST /api/v1/issues/{id}/upload
Content-Type: multipart/form-data
```

**Path Parameters:**
- `id` (integer): Issue ID

**Form Data:**
- `image` (file): Image file (max 10MB)

**Supported formats:**
- JPEG
- PNG
- GIF
- WebP

**Example Request:**
```bash
curl -X POST \
  -F "image=@pothole.jpg" \
  http://localhost:8080/api/v1/issues/1/upload
```

**Example Response:**
```json
{
  "image_url": "/uploads/1.jpg"
}
```

#### Vote for Issue

Upvote an issue to increase its visibility.

```http
POST /api/v1/issues/{id}/vote
```

**Path Parameters:**
- `id` (integer): Issue ID

**Example Request:**
```bash
curl -X POST http://localhost:8080/api/v1/issues/1/vote
```

**Example Response:**
```json
{
  "id": 1,
  "title": "Large pothole on Main Street",
  "description": "Dangerous pothole near the intersection with 5th Ave",
  "category": "pothole",
  "status": "new",
  "latitude": 37.7749,
  "longitude": -122.4194,
  "address": "123 Main St, San Francisco, CA",
  "image_url": "/uploads/1.jpg",
  "reporter_name": "John Doe",
  "reporter_email": "john@example.com",
  "votes": 6,
  "created_at": "2025-01-15T10:30:00Z",
  "updated_at": "2025-01-15T12:30:00Z"
}
```

### Statistics

#### Get Statistics

Retrieve platform-wide statistics.

```http
GET /api/v1/stats
```

**Example Request:**
```bash
curl http://localhost:8080/api/v1/stats
```

**Example Response:**
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

### Health Checks

#### Health Check

Check if the application and database are healthy.

```http
GET /health
```

**Example Response:**
```json
{
  "status": "healthy",
  "database": "connected"
}
```

#### Readiness Check

Check if the application is ready to serve traffic.

```http
GET /ready
```

**Example Response:**
```json
{
  "status": "ready"
}
```

## Code Examples

### JavaScript/TypeScript

```javascript
// Create an issue
async function createIssue(issueData) {
  const response = await fetch('http://localhost:8080/api/v1/issues', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(issueData),
  });

  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error);
  }

  return await response.json();
}

// Upload image
async function uploadImage(issueId, imageFile) {
  const formData = new FormData();
  formData.append('image', imageFile);

  const response = await fetch(`http://localhost:8080/api/v1/issues/${issueId}/upload`, {
    method: 'POST',
    body: formData,
  });

  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error);
  }

  return await response.json();
}
```

### Python

```python
import requests

# Create an issue
def create_issue(issue_data):
    response = requests.post(
        'http://localhost:8080/api/v1/issues',
        json=issue_data
    )
    response.raise_for_status()
    return response.json()

# Upload image
def upload_image(issue_id, image_path):
    with open(image_path, 'rb') as f:
        files = {'image': f}
        response = requests.post(
            f'http://localhost:8080/api/v1/issues/{issue_id}/upload',
            files=files
        )
        response.raise_for_status()
        return response.json()
```

### cURL

```bash
# List issues
curl "http://localhost:8080/api/v1/issues?category=pothole&limit=10"

# Create issue
curl -X POST http://localhost:8080/api/v1/issues \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Broken sidewalk",
    "description": "Cracked sidewalk poses tripping hazard",
    "category": "accessibility",
    "latitude": 37.7749,
    "longitude": -122.4194
  }'

# Upload image
curl -X POST http://localhost:8080/api/v1/issues/1/upload \
  -F "image=@photo.jpg"

# Vote for issue
curl -X POST http://localhost:8080/api/v1/issues/1/vote

# Get statistics
curl http://localhost:8080/api/v1/stats
```

## WebSocket Support

WebSocket support for real-time updates is planned for a future release.

## Pagination

For endpoints that return lists (like `/issues`), use the `limit` and `offset` parameters:

```bash
# Get first page (10 items)
curl "http://localhost:8080/api/v1/issues?limit=10&offset=0"

# Get second page (10 items)
curl "http://localhost:8080/api/v1/issues?limit=10&offset=10"

# Get third page (10 items)
curl "http://localhost:8080/api/v1/issues?limit=10&offset=20"
```

## CORS

The API supports Cross-Origin Resource Sharing (CORS). Configure allowed origins via the `CORS_ALLOWED_ORIGINS` environment variable.

Default headers:
- `Access-Control-Allow-Origin`: Configurable
- `Access-Control-Allow-Methods`: GET, POST, PUT, PATCH, DELETE, OPTIONS
- `Access-Control-Allow-Headers`: Accept, Authorization, Content-Type, X-CSRF-Token
- `Access-Control-Max-Age`: 300

## Changelog

### v1.0.0 (Current)
- Initial release
- CRUD operations for issues
- Image upload support
- Voting system
- Statistics endpoint
- Health checks
- Rate limiting
- CORS support
