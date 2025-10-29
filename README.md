# go-ctl - Go Project Initializr

<div align="center">

![go-ctl Logo](https://img.shields.io/badge/go--ctl-Go%20Project%20Generator-blue?style=for-the-badge&logo=go)

**A modern, web-based Go project generator inspired by Spring Boot Initializr**

[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat-square&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green?style=flat-square)](LICENSE)
[![Made by](https://img.shields.io/badge/Made%20by-systemctl-purple?style=flat-square)](https://github.com/syst3mctl)

[🚀 Live Demo](#getting-started) • [📚 Documentation](#documentation) • [🎯 Features](#features) • [🏗️ Architecture](#architecture)

</div>

## 🎯 Overview

`go-ctl` is a sophisticated web application that generates production-ready Go projects with clean architecture. Simply configure your project requirements through an intuitive web interface and download a complete, runnable Go application with best practices built-in.

### ✨ Key Highlights

- 🎨 **Beautiful Web Interface** - Modern, responsive UI with interactive file explorer
- 🏗️ **Clean Architecture** - Enforced separation of concerns in generated projects  
- 🚀 **Multiple Frameworks** - Support for Gin, Echo, Fiber, Chi, and net/http
- 💾 **Database Ready** - GORM, sqlx, MongoDB, Redis integrations
- 📦 **Package Discovery** - Real-time search and dependency management
- 🔍 **Project Preview** - Interactive file browser with syntax highlighting
- ⚡ **Instant Download** - In-memory ZIP generation for fast delivery

## 🚀 Getting Started

### Prerequisites

- **Go 1.23+** installed on your system
- **Git** for cloning the repository

### Quick Start

1. **Clone the repository**
   ```bash
   git clone https://github.com/syst3mctl/go-ctl.git
   cd go-ctl
   ```

2. **Install dependencies**
   ```bash
   go mod tidy
   ```

3. **Run the server**
   ```bash
   go run cmd/server/*.go
   ```

4. **Open your browser**
   ```
   http://localhost:8080
   ```

🎉 **That's it!** Start generating Go projects through the web interface.

### Alternative Ports

If port 8080 is busy:
```bash
PORT=8081 go run cmd/server/*.go
```

## 📱 User Interface

### Main Generator Interface
![Main Interface](docs/images/main-interface.png)

### Interactive File Explorer
![File Explorer](docs/images/file-explorer-modal.png)

### Package Search & Management
![Package Search](docs/images/package-search.png)

## 🎯 Features

### 📋 Project Configuration
- **Project Naming** - Custom project names with Go module support
- **Go Versions** - Support for Go 1.20, 1.21, 1.22, 1.23
- **Framework Selection** - Choose your preferred web framework
- **Database Integration** - Multiple database and driver options

### 🌐 Web Frameworks
| Framework | Description | Features |
|-----------|-------------|----------|
| **Gin** | High-performance HTTP framework | Middleware, JSON validation, error handling |
| **Echo** | Minimalist web framework | Built-in middleware, data binding |
| **Fiber** | Express-inspired framework | Fast HTTP, low memory footprint |
| **Chi** | Lightweight router | Composable middleware, context-aware |
| **net/http** | Standard library | Pure Go implementation |

### 💾 Database Support
**Databases:**
- PostgreSQL
- MySQL  
- SQLite
- MongoDB
- Redis
- BigQuery

**Drivers/ORMs:**
- **GORM** - Full-featured ORM with associations and migrations
- **sqlx** - Enhanced database/sql with easier scanning
- **Ent** - Schema-first entity framework
- **MongoDB Driver** - Official MongoDB Go driver
- **Redis Client** - Advanced Redis client with clustering
- **database/sql** - Standard library SQL interface

### ⚙️ Additional Features
- **Development Tools**
  - `.gitignore` - Comprehensive Go ignore patterns
  - `Makefile` - Build automation and common tasks
  - `Air` - Hot reload configuration for development
  - `.env.example` - Environment variable templates

- **Production Features**
  - `Docker` - Multi-stage Dockerfile and docker-compose
  - `JWT` - JSON Web Token authentication
  - `CORS` - Cross-Origin Resource Sharing middleware
  - `Logging` - Structured logging with zerolog
  - `Config` - Advanced configuration with Viper
  - `Testing` - Test setup with testify framework

### 🔍 Interactive Project Explorer

Our **standout feature** - a modal-based file explorer that lets you:

- **📁 Browse Structure** - Navigate through the generated project hierarchy
- **👁️ Preview Files** - Click files to view syntax-highlighted content
- **📋 Copy Content** - One-click copy to clipboard
- **🎨 Syntax Highlighting** - Beautiful code preview for Go, JSON, YAML, etc.
- **💾 Direct Download** - Generate and download from the modal

### 📦 Package Management
- **🔍 Real-time Search** - Find packages from pkg.go.dev
- **➕ Easy Selection** - Click to add dependencies
- **❌ Visual Removal** - Remove packages with a click
- **✅ Duplicate Prevention** - Automatic validation

## 🏗️ Architecture

### Application Stack
```
┌─────────────────┐
│   Web Browser   │ ← Tailwind CSS + HTMX
├─────────────────┤
│   Go HTTP Server│ ← net/http + html/template  
├─────────────────┤
│ Generation Engine│ ← Template System
├─────────────────┤
│  Metadata Layer │ ← JSON Configuration
└─────────────────┘
```

### Generated Project Structure
```
my-go-app/
├── cmd/my-go-app/
│   └── main.go              # Application entry point
├── internal/                # Private application code
│   ├── config/
│   │   └── config.go        # Configuration management
│   ├── domain/
│   │   └── model.go         # Business entities
│   ├── service/
│   │   └── service.go       # Business logic
│   ├── handler/
│   │   └── handler.go       # HTTP handlers
│   └── storage/             # Data layer
│       └── gorm/
│           └── gorm.go      # Database implementation
├── .env.example             # Environment template
├── .gitignore              # Git ignore patterns
├── Makefile                # Build automation
├── Dockerfile              # Container build
├── docker-compose.yml      # Service orchestration
└── go.mod                  # Go module definition
```

### Design Principles
- **🏛️ Clean Architecture** - Separation of concerns enforced
- **🔌 Dependency Injection** - Interfaces over implementations
- **⚡ Performance** - In-memory operations, no temporary files
- **🛡️ Security** - Input validation and sanitization
- **📊 Observability** - Structured logging and error handling

## 🛠️ Development

### Project Structure
```
go-ctl/
├── cmd/server/              # Web application entry point
│   ├── main.go             # Server setup and routing
│   ├── handlers.go         # HTTP request handlers
│   └── templates.go        # HTML templates
├── internal/
│   ├── generator/          # Core generation engine
│   │   └── generator.go    # Template processing and ZIP creation
│   └── metadata/           # Configuration management
│       └── options.go      # Project options and validation
├── templates/              # Project generation templates
│   ├── base/              # Core files (go.mod, README, config)
│   ├── features/          # Optional features (Docker, Makefile)
│   ├── http/              # Framework-specific implementations  
│   └── database/          # Database layer templates
├── static/                # Static web assets
├── options.json           # Available project options
└── go.mod                 # Module dependencies
```

### Building from Source

1. **Clone and setup**
   ```bash
   git clone https://github.com/syst3mctl/go-ctl.git
   cd go-ctl
   go mod tidy
   ```

2. **Run tests**
   ```bash
   go test ./...
   ```

3. **Build binary**
   ```bash
   go build -o bin/go-ctl cmd/server/*.go
   ```

4. **Run production build**
   ```bash
   ./bin/go-ctl
   ```

### Adding New Features

1. **Add to options.json** - Define new framework/database/feature
2. **Create template** - Add template file in appropriate directory
3. **Update generator** - Modify generation logic if needed
4. **Test thoroughly** - Ensure generated projects compile and run

### Template Development

Templates use Go's `text/template` with custom functions:
```go
// Example template usage
{{.ProjectName}}              // User's project name
{{.GoVersion}}               // Selected Go version  
{{if .HasFeature "docker"}}  // Conditional generation
{{range .GetAllImports}}     // Iterate over imports
{{end}}
```

## 🎨 UI/UX Design

### Design Philosophy
- **🎯 Simplicity** - Complex functionality made simple
- **⚡ Speed** - Fast interactions with immediate feedback  
- **📱 Responsive** - Works beautifully on all screen sizes
- **♿ Accessible** - Keyboard navigation and screen reader friendly

### Technology Choices
- **Tailwind CSS** - Utility-first styling for rapid development
- **HTMX** - HTML-over-the-wire for dynamic interactions
- **Font Awesome** - Consistent iconography throughout
- **Prism.js** - Beautiful syntax highlighting in file preview

## 📚 Documentation

### API Endpoints
| Endpoint | Method | Purpose |
|----------|--------|---------|
| `/` | GET | Main project generator interface |
| `/generate` | POST | Generate and download project ZIP |
| `/explore` | POST | Get project structure for preview |
| `/search-packages` | GET | Search pkg.go.dev for packages |
| `/add-package` | POST | Add package to selection |
| `/file-content` | GET | Get file content for modal preview |

### Configuration Reference
See [`options.json`](options.json) for complete configuration schema.

### Generated Project Usage
Every generated project includes:
- **README.md** - Complete setup and usage instructions
- **Makefile** - Common development tasks
- **Configuration** - Environment variable setup
- **Examples** - Working endpoint implementations

## 🤝 Contributing

We welcome contributions! Here's how you can help:

### Ways to Contribute
- 🐛 **Report Bugs** - Found an issue? Let us know!
- 💡 **Suggest Features** - Ideas for new frameworks/databases/features
- 📖 **Improve Documentation** - Help make our docs clearer
- 🎨 **UI/UX Improvements** - Make the interface even better
- 🧪 **Add Tests** - Help us maintain quality

### Development Process
1. **Fork** the repository
2. **Create** a feature branch (`git checkout -b feature/amazing-feature`)
3. **Commit** your changes (`git commit -m 'Add amazing feature'`)
4. **Push** to the branch (`git push origin feature/amazing-feature`)
5. **Open** a Pull Request

### Code Standards
- Follow Go conventions and `gofmt`
- Write tests for new functionality
- Update documentation for changes
- Keep commits atomic and descriptive

## 📄 License

This project is licensed under the **MIT License** - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- **Spring Boot Initializr** - Inspiration for the concept and UI design
- **Go Community** - For the amazing ecosystem of packages and tools
- **HTMX** - Making dynamic web interfaces simple and elegant
- **Contributors** - Everyone who has helped improve this project

## 🌟 Support

If you find `go-ctl` helpful:

- ⭐ **Star** this repository
- 🐛 **Report issues** you encounter  
- 💡 **Share ideas** for new features
- 📢 **Spread the word** to other Go developers

## 📊 Project Stats

- **Languages**: Go, HTML, CSS, JavaScript
- **Architecture**: Clean Architecture, Template-driven Generation
- **Dependencies**: Minimal, standard library focused
- **Performance**: Sub-second project generation
- **Compatibility**: Go 1.20+ on all platforms

---

<div align="center">

**Built with ❤️ by [systemctl](https://github.com/syst3mctl)**

*Accelerating Go development, one project at a time*

</div>