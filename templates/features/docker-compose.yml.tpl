version: '3.8'

services:
  {{.ProjectName}}:
    build:
      context: .
      dockerfile: Dockerfile
    ports:
      - "{{if .HasFeature "config"}}${PORT:-8080}{{else}}8080{{end}}:{{if .HasFeature "config"}}${PORT:-8080}{{else}}8080{{end}}"
    environment:
{{if .HasFeature "env"}}      - APP_ENV=${APP_ENV:-development}
      - APP_PORT=${PORT:-8080}
{{if ne .DbDriver.ID ""}}{{if eq .Database.ID "postgres"}}      - POSTGRES_HOST=postgres
      - POSTGRES_PORT=5432
      - POSTGRES_DB=${POSTGRES_DB:-{{.ProjectName}}}
      - POSTGRES_USER=${POSTGRES_USER:-postgres}
      - POSTGRES_PASSWORD=${POSTGRES_PASSWORD:-password}
{{else if eq .Database.ID "mysql"}}      - MYSQL_HOST=mysql
      - MYSQL_PORT=3306
      - MYSQL_DATABASE=${MYSQL_DATABASE:-{{.ProjectName}}}
      - MYSQL_USER=${MYSQL_USER:-root}
      - MYSQL_PASSWORD=${MYSQL_PASSWORD:-password}
{{else if eq .Database.ID "mongodb"}}      - MONGO_HOST=mongo
      - MONGO_PORT=27017
      - MONGO_DATABASE=${MONGO_DATABASE:-{{.ProjectName}}}
      - MONGO_USER=${MONGO_USER:-}
      - MONGO_PASSWORD=${MONGO_PASSWORD:-}
{{else if eq .Database.ID "redis"}}      - REDIS_HOST=redis
      - REDIS_PORT=6379
      - REDIS_PASSWORD=${REDIS_PASSWORD:-}
      - REDIS_DB=${REDIS_DB:-0}
{{end}}{{end}}{{if .HasFeature "jwt"}}      - JWT_SECRET=${JWT_SECRET:-your-super-secret-jwt-key}
{{end}}{{else}}      - APP_ENV=development
      - APP_PORT=8080
{{if ne .DbDriver.ID ""}}{{if eq .Database.ID "postgres"}}      - POSTGRES_HOST=postgres
      - POSTGRES_PORT=5432
      - POSTGRES_DB={{.ProjectName}}
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=password
{{else if eq .Database.ID "mysql"}}      - MYSQL_HOST=mysql
      - MYSQL_PORT=3306
      - MYSQL_DATABASE={{.ProjectName}}
      - MYSQL_USER=root
      - MYSQL_PASSWORD=password
{{else if eq .Database.ID "mongodb"}}      - MONGO_HOST=mongo
      - MONGO_PORT=27017
      - MONGO_DATABASE={{.ProjectName}}
{{else if eq .Database.ID "redis"}}      - REDIS_HOST=redis
      - REDIS_PORT=6379
      - REDIS_DB=0
{{end}}{{end}}{{if .HasFeature "jwt"}}      - JWT_SECRET=your-super-secret-jwt-key
{{end}}{{end}}
    depends_on:
{{if ne .DbDriver.ID ""}}{{if eq .Database.ID "postgres"}}      - postgres
{{else if eq .Database.ID "mysql"}}      - mysql
{{else if eq .Database.ID "mongodb"}}      - mongo
{{else if eq .Database.ID "redis"}}      - redis
{{end}}{{end}}    volumes:
{{if .HasFeature "env"}}      - ./.env:/root/.env:ro
{{end}}    restart: unless-stopped
    networks:
      - {{.ProjectName}}-network

{{if ne .DbDriver.ID ""}}{{if eq .Database.ID "postgres"}}  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: {{if .HasFeature "env"}}${POSTGRES_DB:-{{.ProjectName}}}{{else}}{{.ProjectName}}{{end}}
      POSTGRES_USER: {{if .HasFeature "env"}}${POSTGRES_USER:-postgres}{{else}}postgres{{end}}
      POSTGRES_PASSWORD: {{if .HasFeature "env"}}${POSTGRES_PASSWORD:-password}{{else}}password{{end}}
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./init.sql:/docker-entrypoint-initdb.d/init.sql:ro
    restart: unless-stopped
    networks:
      - {{.ProjectName}}-network
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${POSTGRES_USER:-postgres}"]
      interval: 30s
      timeout: 10s
      retries: 5

{{else if eq .Database.ID "mysql"}}  mysql:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: {{if .HasFeature "env"}}${MYSQL_PASSWORD:-password}{{else}}password{{end}}
      MYSQL_DATABASE: {{if .HasFeature "env"}}${MYSQL_DATABASE:-{{.ProjectName}}}{{else}}{{.ProjectName}}{{end}}
      MYSQL_USER: {{if .HasFeature "env"}}${MYSQL_USER:-user}{{else}}user{{end}}
      MYSQL_PASSWORD: {{if .HasFeature "env"}}${MYSQL_PASSWORD:-password}{{else}}password{{end}}
    ports:
      - "3306:3306"
    volumes:
      - mysql_data:/var/lib/mysql
      - ./init.sql:/docker-entrypoint-initdb.d/init.sql:ro
    restart: unless-stopped
    networks:
      - {{.ProjectName}}-network
    healthcheck:
      test: ["CMD", "mysqladmin", "ping", "-h", "localhost"]
      interval: 30s
      timeout: 10s
      retries: 5

{{else if eq .Database.ID "mongodb"}}  mongo:
    image: mongo:6.0
    environment:
{{if .HasFeature "env"}}      MONGO_INITDB_ROOT_USERNAME: ${MONGO_USER:-admin}
      MONGO_INITDB_ROOT_PASSWORD: ${MONGO_PASSWORD:-password}
      MONGO_INITDB_DATABASE: ${MONGO_DATABASE:-{{.ProjectName}}}
{{else}}      MONGO_INITDB_ROOT_USERNAME: admin
      MONGO_INITDB_ROOT_PASSWORD: password
      MONGO_INITDB_DATABASE: {{.ProjectName}}
{{end}}    ports:
      - "27017:27017"
    volumes:
      - mongo_data:/data/db
      - ./mongo-init.js:/docker-entrypoint-initdb.d/mongo-init.js:ro
    restart: unless-stopped
    networks:
      - {{.ProjectName}}-network
    healthcheck:
      test: echo 'db.runCommand("ping").ok' | mongosh localhost:27017/test --quiet
      interval: 30s
      timeout: 10s
      retries: 5

{{else if eq .Database.ID "redis"}}  redis:
    image: redis:7-alpine
    command: redis-server{{if .HasFeature "env"}} --requirepass ${REDIS_PASSWORD}{{end}}
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
      - ./redis.conf:/usr/local/etc/redis/redis.conf:ro
    restart: unless-stopped
    networks:
      - {{.ProjectName}}-network
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 30s
      timeout: 10s
      retries: 5

{{end}}{{end}}networks:
  {{.ProjectName}}-network:
    driver: bridge

{{if ne .DbDriver.ID ""}}volumes:
{{if eq .Database.ID "postgres"}}  postgres_data:
{{else if eq .Database.ID "mysql"}}  mysql_data:
{{else if eq .Database.ID "mongodb"}}  mongo_data:
{{else if eq .Database.ID "redis"}}  redis_data:
{{end}}{{end}}
