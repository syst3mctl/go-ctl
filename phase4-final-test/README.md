# phase4-final-test

phase4-final-test is a Go web application generated using [go-ctl](https://github.com/syst3mctl/go-ctl).

## 🚀 Getting Started

### Prerequisites

- Go 1.23 or later
- PostgreSQL database

### Installation

1. Clone this repository:
```bash
git clone <your-repo-url>
cd phase4-final-test
```

2. Install dependencies:
```bash
go mod tidy
```



### Running the Application

#### Using Makefile
```bash
# Development mode
make dev

# Build and run
make build
make run

# Run tests
make test
```

#### Manual Commands
```bash
# Run the application
go run cmd/phase4-final-test/main.go
```

The server will start on http://localhost:8080

## 📁 Project Structure

```
phase4-final-test/
├── cmd/phase4-final-test/
│   └── main.go                 # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go          # Configuration management
│   ├── domain/
│   │   └── model.go           # Domain models and business entities
│   ├── service/
│   │   └── service.go         # Business logic and service interfaces
│   ├── storage/
│   │   └── gorm/
│   │       └── gorm.go    # Database layer implementation
│   └── handler/
│       └── handler.go         # HTTP handlers and routing
├── docker-compose.yml         # Docker services configuration
├── Dockerfile                 # Container build instructions
├── Makefile                   # Build automation


└── go.mod                     # Go module definition
```

## 🛠️ Technology Stack

- **Language**: Go 1.23
- **Web Framework**: Gin - A high-performance HTTP web framework written in Go.
  - **Database**: PostgreSQL - A powerful, open source object-relational database system.
  - **Database Driver**: GORM - The fantastic ORM library for Golang, aims to be developer friendly.

### Additional Features

- **Docker**: Add Dockerfile and docker-compose.yml for containerization.

- **Makefile**: Add a Makefile with common build, test, and run targets.


## 📚 API Documentation

### Gin Framework
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


## 🗄️ Database

### GORM
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

Configuration is managed through the `internal/config` package. Modify the configuration struct as needed for your application requirements.

## 🚀 Deployment

### Docker

1. Build the image:
```bash
docker build -t phase4-final-test .
```

2. Run with Docker Compose:
```bash
docker-compose up -d
```

### Production

```bash
# Build optimized binary
go build -ldflags="-w -s" -o bin/phase4-final-test cmd/phase4-final-test/main.go

# Run the binary
./bin/phase4-final-test
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
- Gin - A high-performance HTTP web framework written in Go.
  - GORM - The fantastic ORM library for Golang, aims to be developer friendly.

---

**Happy Coding! 🎉**
