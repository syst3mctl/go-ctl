# {{.ProjectName}} Environment Configuration
# Copy this file to .env and update the values as needed

# Server Configuration
HOST=localhost
PORT=8080

{{if ne .Database.ID ""}}# Database Configuration
{{if eq .Database.ID "postgres"}}DB_HOST=localhost
DB_PORT=5432
DB_NAME={{.ProjectName}}_db
DB_USER=postgres
DB_PASSWORD=password
DB_SSLMODE=disable
DATABASE_URL=postgres://postgres:password@localhost:5432/{{.ProjectName}}_db?sslmode=disable
{{else if eq .Database.ID "mysql"}}DB_HOST=localhost
DB_PORT=3306
DB_NAME={{.ProjectName}}_db
DB_USER=root
DB_PASSWORD=password
DATABASE_URL=root:password@tcp(localhost:3306)/{{.ProjectName}}_db?parseTime=true
{{else if eq .Database.ID "sqlite"}}DB_PATH=./{{.ProjectName}}.db
{{else if eq .Database.ID "mongodb"}}MONGO_URI=mongodb://localhost:27017/{{.ProjectName}}_db
{{else if eq .Database.ID "redis"}}REDIS_URL=redis://localhost:6379/0
REDIS_PASSWORD=
REDIS_DB=0
{{else if eq .Database.ID "bigquery"}}GOOGLE_APPLICATION_CREDENTIALS=path/to/service-account-key.json
BIGQUERY_PROJECT_ID=your-project-id
BIGQUERY_DATASET_ID=your-dataset-id
{{end}}

{{end}}# Application Configuration
APP_NAME={{.ProjectName}}
APP_VERSION=1.0.0
APP_ENV=development
APP_DEBUG=true

{{if .HasFeature "jwt"}}# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
JWT_EXPIRES_IN=24  # hours
{{end}}

{{if .HasFeature "logging"}}# Logging Configuration
LOG_LEVEL=info
LOG_FORMAT=json  # json or console
{{end}}

{{if .HasFeature "cors"}}# CORS Configuration
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:8080
CORS_ALLOWED_METHODS=GET,POST,PUT,DELETE,OPTIONS
CORS_ALLOWED_HEADERS=Content-Type,Authorization
CORS_ALLOW_CREDENTIALS=true
{{end}}

# Security Configuration
{{if not .HasFeature "jwt"}}API_SECRET=your-api-secret-key{{end}}
RATE_LIMIT_REQUESTS=100  # requests per minute
RATE_LIMIT_WINDOW=1      # minutes

# External Services
{{if .HasFeature "config"}}# Email Configuration (if needed)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password

# AWS Configuration (if needed)
AWS_REGION=us-east-1
AWS_ACCESS_KEY_ID=your-access-key
AWS_SECRET_ACCESS_KEY=your-secret-key
AWS_S3_BUCKET=your-bucket-name

# Google Cloud Configuration (if needed)
GOOGLE_CLOUD_PROJECT=your-project-id
GOOGLE_APPLICATION_CREDENTIALS=path/to/service-account.json

# Redis Cache (if using Redis for caching)
CACHE_TTL=3600  # seconds
{{end}}

# Development Tools
{{if .HasFeature "air"}}# Air (Hot Reload) Configuration
AIR_ENABLED=true
{{end}}

# Health Check Configuration
HEALTH_CHECK_TIMEOUT=30s
HEALTH_CHECK_INTERVAL=10s

# Monitoring and Observability
METRICS_ENABLED=true
TRACING_ENABLED=false
PROFILING_ENABLED=false

# File Upload Configuration
MAX_FILE_SIZE=10MB
ALLOWED_FILE_TYPES=jpg,jpeg,png,gif,pdf,txt,csv

# Pagination Defaults
DEFAULT_PAGE_SIZE=20
MAX_PAGE_SIZE=100

# Session Configuration (if using sessions)
SESSION_SECRET=your-session-secret-key
SESSION_TIMEOUT=3600  # seconds

# Environment-specific overrides
# Uncomment and modify as needed for different environments

# Production overrides
# APP_ENV=production
# APP_DEBUG=false
# LOG_LEVEL=warn

# Testing overrides
# APP_ENV=testing
# DB_NAME={{.ProjectName}}_test_db
