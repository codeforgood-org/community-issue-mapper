# Deployment Guide

This guide covers various deployment options for the Community Issue Mapper application.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Docker Deployment](#docker-deployment)
- [Cloud Deployments](#cloud-deployments)
  - [AWS ECS/Fargate](#aws-ecsfargate)
  - [Google Cloud Run](#google-cloud-run)
  - [DigitalOcean App Platform](#digitalocean-app-platform)
  - [Heroku](#heroku)
- [Database Setup](#database-setup)
- [Environment Configuration](#environment-configuration)
- [SSL/TLS Configuration](#ssltls-configuration)
- [Monitoring and Logging](#monitoring-and-logging)

## Prerequisites

- Docker and Docker Compose (for containerized deployment)
- PostgreSQL 14+ database
- Domain name (for production deployment)
- SSL certificate (for HTTPS)

## Docker Deployment

### Local Docker Deployment

1. **Clone the repository**
   ```bash
   git clone https://github.com/codeforgood-org/community-issue-mapper.git
   cd community-issue-mapper
   ```

2. **Configure environment**
   ```bash
   cp .env.example .env
   # Edit .env with your settings
   ```

3. **Build and start services**
   ```bash
   docker-compose up -d
   ```

4. **Run migrations**
   ```bash
   docker-compose run migrate
   ```

5. **Verify deployment**
   ```bash
   curl http://localhost:8080/health
   ```

### Production Docker Deployment

For production, use a separate docker-compose.prod.yml:

```yaml
version: '3.8'

services:
  app:
    image: your-registry/community-issue-mapper:latest
    environment:
      - ENV=production
      - DB_HOST=your-db-host
      - DB_PASSWORD=${DB_PASSWORD}
      - CORS_ALLOWED_ORIGINS=https://yourdomain.com
    restart: always
    ports:
      - "8080:8080"
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
```

Deploy:
```bash
docker-compose -f docker-compose.prod.yml up -d
```

## Cloud Deployments

### AWS ECS/Fargate

#### 1. Build and Push Docker Image

```bash
# Authenticate with ECR
aws ecr get-login-password --region us-east-1 | \
  docker login --username AWS --password-stdin YOUR_AWS_ACCOUNT.dkr.ecr.us-east-1.amazonaws.com

# Build image
docker build -t community-issue-mapper:latest .

# Tag image
docker tag community-issue-mapper:latest \
  YOUR_AWS_ACCOUNT.dkr.ecr.us-east-1.amazonaws.com/community-issue-mapper:latest

# Push image
docker push YOUR_AWS_ACCOUNT.dkr.ecr.us-east-1.amazonaws.com/community-issue-mapper:latest
```

#### 2. Create Task Definition

Create `task-definition.json`:

```json
{
  "family": "community-issue-mapper",
  "networkMode": "awsvpc",
  "requiresCompatibilities": ["FARGATE"],
  "cpu": "256",
  "memory": "512",
  "containerDefinitions": [
    {
      "name": "app",
      "image": "YOUR_AWS_ACCOUNT.dkr.ecr.us-east-1.amazonaws.com/community-issue-mapper:latest",
      "portMappings": [
        {
          "containerPort": 8080,
          "protocol": "tcp"
        }
      ],
      "environment": [
        {"name": "ENV", "value": "production"},
        {"name": "PORT", "value": "8080"},
        {"name": "DB_HOST", "value": "your-rds-endpoint"},
        {"name": "DB_NAME", "value": "community_issues"}
      ],
      "secrets": [
        {
          "name": "DB_PASSWORD",
          "valueFrom": "arn:aws:secretsmanager:region:account:secret:db-password"
        }
      ],
      "logConfiguration": {
        "logDriver": "awslogs",
        "options": {
          "awslogs-group": "/ecs/community-issue-mapper",
          "awslogs-region": "us-east-1",
          "awslogs-stream-prefix": "ecs"
        }
      }
    }
  ]
}
```

#### 3. Create ECS Service

```bash
aws ecs create-service \
  --cluster your-cluster \
  --service-name community-issue-mapper \
  --task-definition community-issue-mapper \
  --desired-count 2 \
  --launch-type FARGATE \
  --network-configuration "awsvpcConfiguration={subnets=[subnet-xxx],securityGroups=[sg-xxx],assignPublicIp=ENABLED}" \
  --load-balancers "targetGroupArn=arn:aws:elasticloadbalancing:...,containerName=app,containerPort=8080"
```

### Google Cloud Run

#### 1. Build and Push to Container Registry

```bash
# Set project
gcloud config set project YOUR_PROJECT_ID

# Build image
gcloud builds submit --tag gcr.io/YOUR_PROJECT_ID/community-issue-mapper

# Or use Docker
docker build -t gcr.io/YOUR_PROJECT_ID/community-issue-mapper .
docker push gcr.io/YOUR_PROJECT_ID/community-issue-mapper
```

#### 2. Deploy to Cloud Run

```bash
gcloud run deploy community-issue-mapper \
  --image gcr.io/YOUR_PROJECT_ID/community-issue-mapper \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated \
  --set-env-vars "ENV=production,DB_HOST=your-cloud-sql-ip" \
  --set-secrets "DB_PASSWORD=db-password:latest" \
  --vpc-connector your-vpc-connector \
  --max-instances 10 \
  --memory 512Mi
```

#### 3. Connect to Cloud SQL

```bash
gcloud run services update community-issue-mapper \
  --add-cloudsql-instances YOUR_PROJECT_ID:REGION:INSTANCE_NAME
```

### DigitalOcean App Platform

#### 1. Create app.yaml

```yaml
name: community-issue-mapper
services:
- name: web
  github:
    repo: your-username/community-issue-mapper
    branch: main
    deploy_on_push: true
  dockerfile_path: Dockerfile
  http_port: 8080
  instance_count: 2
  instance_size_slug: basic-xs
  envs:
  - key: ENV
    value: production
  - key: DB_HOST
    value: ${db.HOSTNAME}
  - key: DB_PORT
    value: ${db.PORT}
  - key: DB_NAME
    value: ${db.DATABASE}
  - key: DB_USER
    value: ${db.USERNAME}
  - key: DB_PASSWORD
    value: ${db.PASSWORD}
    type: SECRET

databases:
- name: db
  engine: PG
  production: true
  version: "14"
```

#### 2. Deploy

```bash
# Install doctl
brew install doctl  # macOS

# Authenticate
doctl auth init

# Create app
doctl apps create --spec app.yaml

# Or deploy via UI
# Visit https://cloud.digitalocean.com/apps
```

### Heroku

#### 1. Create Heroku App

```bash
# Install Heroku CLI
brew install heroku  # macOS

# Login
heroku login

# Create app
heroku create your-app-name

# Add PostgreSQL
heroku addons:create heroku-postgresql:standard-0
```

#### 2. Configure Container Deployment

Create `heroku.yml`:

```yaml
build:
  docker:
    web: Dockerfile
run:
  web: ./main
```

#### 3. Deploy

```bash
# Set stack to container
heroku stack:set container

# Push to Heroku
git push heroku main

# Run migrations
heroku run migrate -path migrations -database $DATABASE_URL up

# Scale dynos
heroku ps:scale web=2

# View logs
heroku logs --tail
```

## Database Setup

### PostgreSQL on AWS RDS

1. **Create RDS Instance**
   ```bash
   aws rds create-db-instance \
     --db-instance-identifier community-issues-db \
     --db-instance-class db.t3.micro \
     --engine postgres \
     --engine-version 14.7 \
     --master-username issueapp \
     --master-user-password YOUR_PASSWORD \
     --allocated-storage 20 \
     --vpc-security-group-ids sg-xxx \
     --db-subnet-group-name your-subnet-group \
     --backup-retention-period 7 \
     --preferred-backup-window 03:00-04:00 \
     --preferred-maintenance-window sun:04:00-sun:05:00
   ```

2. **Run Migrations**
   ```bash
   migrate -path backend/migrations \
     -database "postgres://issueapp:PASSWORD@your-rds-endpoint:5432/community_issues?sslmode=require" up
   ```

### Google Cloud SQL

```bash
gcloud sql instances create community-issues-db \
  --database-version=POSTGRES_14 \
  --tier=db-f1-micro \
  --region=us-central1 \
  --backup \
  --backup-start-time=03:00

gcloud sql databases create community_issues \
  --instance=community-issues-db

gcloud sql users create issueapp \
  --instance=community-issues-db \
  --password=YOUR_PASSWORD
```

## Environment Configuration

### Production Environment Variables

```bash
# Server
ENV=production
PORT=8080

# Database (use managed service)
DB_HOST=your-db-host
DB_PORT=5432
DB_USER=issueapp
DB_PASSWORD=strong-password-here
DB_NAME=community_issues
DB_SSL_MODE=require

# CORS (set to your domain)
CORS_ALLOWED_ORIGINS=https://yourdomain.com,https://www.yourdomain.com

# Uploads
MAX_UPLOAD_SIZE=10485760
UPLOAD_DIR=/app/uploads

# Rate Limiting
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_DURATION=1m
```

### Secrets Management

#### AWS Secrets Manager

```bash
# Create secret
aws secretsmanager create-secret \
  --name community-issues/db-password \
  --secret-string "your-password"

# Reference in ECS task definition
"secrets": [
  {
    "name": "DB_PASSWORD",
    "valueFrom": "arn:aws:secretsmanager:region:account:secret:community-issues/db-password"
  }
]
```

#### Google Secret Manager

```bash
# Create secret
echo -n "your-password" | gcloud secrets create db-password --data-file=-

# Grant access
gcloud secrets add-iam-policy-binding db-password \
  --member=serviceAccount:YOUR_SERVICE_ACCOUNT \
  --role=roles/secretmanager.secretAccessor
```

## SSL/TLS Configuration

### Using NGINX as Reverse Proxy

Create `nginx.conf`:

```nginx
upstream app {
    server app:8080;
}

server {
    listen 80;
    server_name yourdomain.com;
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name yourdomain.com;

    ssl_certificate /etc/nginx/ssl/cert.pem;
    ssl_certificate_key /etc/nginx/ssl/key.pem;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;

    client_max_body_size 10M;

    location / {
        proxy_pass http://app;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### Using Let's Encrypt

```bash
# Install certbot
apt-get install certbot python3-certbot-nginx

# Get certificate
certbot --nginx -d yourdomain.com -d www.yourdomain.com

# Auto-renewal
certbot renew --dry-run
```

## Monitoring and Logging

### Health Check Monitoring

Use a service like UptimeRobot or Pingdom to monitor:
- `https://yourdomain.com/health`
- `https://yourdomain.com/ready`

### Application Monitoring

#### Prometheus + Grafana

Add metrics endpoint to your app and scrape with Prometheus.

#### CloudWatch (AWS)

```bash
# Create log group
aws logs create-log-group --log-group-name /app/community-issue-mapper

# Stream logs from container
# (Configured in ECS task definition)
```

#### Google Cloud Logging

Logs are automatically collected in Cloud Run and viewable in Cloud Logging.

### Database Monitoring

- AWS RDS: Use CloudWatch metrics
- Google Cloud SQL: Use Cloud Monitoring
- Set up alerts for:
  - High CPU usage
  - Low storage space
  - Connection count
  - Slow queries

## Backup Strategy

### Database Backups

```bash
# Manual backup
pg_dump -h your-db-host -U issueapp -d community_issues > backup_$(date +%Y%m%d).sql

# Automated backups
# Configure in your managed database service
# AWS RDS: 7-35 day retention
# Google Cloud SQL: 7-365 day retention
```

### File Storage Backups

For uploads, consider using cloud storage:
- AWS S3
- Google Cloud Storage
- DigitalOcean Spaces

## Scaling

### Horizontal Scaling

- Add more container instances
- Use load balancer to distribute traffic
- Scale database read replicas for read-heavy workloads

### Vertical Scaling

- Increase container memory/CPU
- Upgrade database instance size

### Caching

Consider adding Redis for:
- API response caching
- Session storage (future auth)
- Rate limiting data

## Troubleshooting

### Common Issues

1. **Database Connection Failed**
   - Check security groups/firewall rules
   - Verify database credentials
   - Ensure SSL mode matches database configuration

2. **High Memory Usage**
   - Increase container memory limit
   - Check for memory leaks
   - Monitor with profiling tools

3. **Slow Response Times**
   - Add database indexes
   - Enable query caching
   - Use CDN for static assets

### Debug Mode

Set `ENV=development` to enable detailed error messages (not for production).

## Rollback Strategy

```bash
# AWS ECS
aws ecs update-service \
  --cluster your-cluster \
  --service community-issue-mapper \
  --task-definition community-issue-mapper:PREVIOUS_VERSION

# Google Cloud Run
gcloud run services update-traffic community-issue-mapper \
  --to-revisions=PREVIOUS_REVISION=100

# Heroku
heroku releases:rollback v123
```

## Maintenance Windows

Schedule maintenance during low-traffic periods:
1. Announce maintenance in advance
2. Set up maintenance page
3. Perform updates/migrations
4. Test thoroughly
5. Resume normal operations
6. Monitor for issues

---

For questions or issues, please open an issue on GitHub.
