# {{.ProjectName}}

{{.ProjectName}} is a Go web application generated using [go-ctl](https://github.com/syst3mctl/go-ctl).

## 🚀 Getting Started

### Prerequisites

- Go {{.GoVersion}} or later
- {{if eq .Database.ID "postgres"}}PostgreSQL database{{else if eq .Database.ID "mysql"}}MySQL database{{else if eq .Database.ID "mongodb"}}MongoDB database{{else if eq .Database.ID "redis"}}Redis server{{else if eq .Database.ID "sqlite"}}SQLite (no additional setup required){{else}}Database setup as configured{{end}}

### Installation

1. Clone this repository:
```bash
git clone <your-repo-url>
cd {{.ProjectName}}
```

2. Install dependencies:
```bash
go mod tidy
```

{{if .HasFeature "env"}}3. Copy the environment file and configure your settings:
```bash
cp .env.example .env
# Edit .env with your configuration
```
{{end}}

### Running the Application

{{if .HasFeature "makefile"}}#### Using Makefile
```bash
# Development mode{{if .HasFeature "air"}} with hot reload{{end}}
make dev

# Build and run
make build
make run

# Run tests
make test
```

#### Manual Commands{{else}}#### Development{{end}}
```bash
{{if .HasFeature "air"}}# Development with hot reload
air

# Or run directly{{else}}# Run the application{{end}}
go run cmd/{{.ProjectName}}/main.go
```

The server will start on http://localhost:8080

## 📁 Project Structure

```
{{.ProjectName}}/
├── cmd/{{.ProjectName}}/
│   └── main.go                 # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go          # Configuration management
│   ├── domain/
│   │   └── model.go           # Domain models and business entities
│   ├── service/
│   │   └── service.go         # Business logic and service interfaces
│   {{if ne .DbDriver.ID ""}}├── storage/
│   │   └── {{.DbDriver.ID}}/
│   │       └── {{.DbDriver.ID}}.go    # Database layer implementation{{end}}
│   └── handler/
│       └── handler.go         # HTTP handlers and routing
{{if .HasFeature "docker"}}├── docker-compose.yml         # Docker services configuration
├── Dockerfile                 # Container build instructions{{end}}
{{if .HasFeature "makefile"}}├── Makefile                   # Build automation{{end}}
{{if .HasFeature "air"}}├── .air.toml                  # Hot reload configuration{{end}}
{{if .HasFeature "env"}}├── .env.example               # Environment variables template{{end}}
└── go.mod                     # Go module definition
```

## 🛠️ Technology Stack

- **Language**: Go {{.GoVersion}}
- **Web Framework**: {{.HttpPackage.Name}}{{if ne .HttpPackage.Description ""}} - {{.HttpPackage.Description}}{{end}}
{{if ne .Database.ID ""}}  - **Database**: {{.Database.Name}}{{if ne .Database.Description ""}} - {{.Database.Description}}{{end}}{{end}}
{{if ne .DbDriver.ID ""}}  - **Database Driver**: {{.DbDriver.Name}}{{if ne .DbDriver.Description ""}} - {{.DbDriver.Description}}{{end}}{{end}}

### Additional Features
{{range .Features}}
- **{{.Name}}**: {{.Description}}
{{end}}

## 📚 API Documentation

{{if eq .HttpPackage.ID "gin"}}### Gin Framework
This project uses the Gin web framework. Key features:
- High performance HTTP router
- Middleware support
- JSON validation
- Error management

Example endpoint:
```go
r.GET("/api/health", func(c *gin.Context) {
    c.JSON(200, gin.H{"status": "ok"})
})
```
{{else if eq .HttpPackage.ID "echo"}}### Echo Framework
This project uses the Echo web framework. Key features:
- High performance and minimalist
- Built-in middleware
- Data binding and validation
- Template rendering

Example endpoint:
```go
e.GET("/api/health", func(c echo.Context) error {
    return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
})
```
{{else if eq .HttpPackage.ID "fiber"}}### Fiber Framework
This project uses the Fiber web framework. Key features:
- Express-inspired API
- Built on top of Fasthttp
- Low memory footprint
- Rapid server-side programming

Example endpoint:
```go
app.Get("/api/health", func(c *fiber.Ctx) error {
    return c.JSON(fiber.Map{"status": "ok"})
})
```
{{else if eq .HttpPackage.ID "chi"}}### Chi Router
This project uses the Chi router. Key features:
- Lightweight and fast
- Fully compatible with net/http
- Composable middleware stack
- Context-aware

Example endpoint:
```go
r.Get("/api/health", func(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
})
```
{{else}}### Standard Library (net/http)
This project uses Go's standard HTTP library with a custom router implementation.
{{end}}

{{if ne .DbDriver.ID ""}}## 🗄️ Database

{{if eq .DbDriver.ID "gorm"}}### GORM
This project uses GORM as the ORM. Key features:
- Full-featured ORM
- Associations (Has One, Has Many, Belongs To, Many To Many)
- Hooks (Before/After Create/Save/Update/Delete/Find)
- Preloading and Joins
- Transactions, Nested Transactions, Save Point, RollbackTo

Example model:
```go
type User struct {
    ID        uint      `gorm:"primarykey"`
    CreatedAt time.Time
    UpdatedAt time.Time
    Name      string    `gorm:"size:255;not null"`
    Email     string    `gorm:"uniqueIndex;size:255;not null"`
}
```
{{else if eq .DbDriver.ID "sqlx"}}### sqlx
This project uses sqlx for database operations. Key features:
- Extensions to database/sql
- Easier result scanning
- Named parameter support
- Get/Select convenience methods

Example query:
```go
var users []User
err := db.Select(&users, "SELECT * FROM users WHERE active = $1", true)
```
{{else if eq .DbDriver.ID "mongo-driver"}}### MongoDB Driver
This project uses the official MongoDB Go driver. Key features:
- Full MongoDB feature support
- Type-safe operations
- Aggregation pipeline support
- Change streams

Example operation:
```go
collection := client.Database("mydb").Collection("users")
result, err := collection.InsertOne(ctx, User{Name: "John", Email: "john@example.com"})
```
{{else}}### Database/SQL
This project uses Go's standard database/sql package with appropriate drivers.
{{end}}
{{end}}

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...
```

## 🔧 Configuration

{{if .HasFeature "env"}}The application uses environment variables for configuration. See `.env.example` for available options:

```env
# Server Configuration
PORT=8080
HOST=localhost

{{if ne .Database.ID ""}}# Database Configuration
{{if eq .Database.ID "postgres"}}DB_HOST=localhost
DB_PORT=5432
DB_NAME={{.ProjectName}}_db
DB_USER=postgres
DB_PASSWORD=password
DATABASE_URL=postgres://postgres:password@localhost:5432/{{.ProjectName}}_db?sslmode=disable
{{else if eq .Database.ID "mysql"}}DB_HOST=localhost
DB_PORT=3306
DB_NAME={{.ProjectName}}_db
DB_USER=root
DB_PASSWORD=password
DATABASE_URL=root:password@tcp(localhost:3306)/{{.ProjectName}}_db?parseTime=true
{{else if eq .Database.ID "mongodb"}}MONGO_URI=mongodb://localhost:27017/{{.ProjectName}}_db
{{else if eq .Database.ID "redis"}}REDIS_URL=redis://localhost:6379/0
{{else if eq .Database.ID "sqlite"}}DB_PATH=./{{.ProjectName}}.db
{{end}}{{end}}

# Application Configuration
APP_ENV=development
LOG_LEVEL=info
{{if .HasFeature "jwt"}}JWT_SECRET=your-super-secret-jwt-key{{end}}
```
{{else}}Configuration is managed through the `internal/config` package. Modify the configuration struct as needed for your application requirements.{{end}}

## 🚀 Deployment

{{if .HasFeature "docker"}}### Docker

1. Build the image:
```bash
docker build -t {{.ProjectName}} .
```

2. Run with Docker Compose:
```bash
docker-compose up -d
```

### Production{{else}}### Building for Production{{end}}

```bash
# Build optimized binary
go build -ldflags="-w -s" -o bin/{{.ProjectName}} cmd/{{.ProjectName}}/main.go

# Run the binary
./bin/{{.ProjectName}}
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Generated with [go-ctl](https://github.com/syst3mctl/go-ctl) - Go Project Initializr
- {{.HttpPackage.Name}} - {{.HttpPackage.Description}}
{{if ne .DbDriver.ID ""}}  - {{.DbDriver.Name}} - {{.DbDriver.Description}}{{end}}

---

**Happy Coding! 🎉**
