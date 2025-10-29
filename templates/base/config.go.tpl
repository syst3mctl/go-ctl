package config

import (
	"os"
	"strconv"
{{if .HasFeature "config"}}
	"github.com/spf13/viper"
{{end}}
{{if .HasFeature "logging"}}
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
{{end}}
)

// Config holds all configuration for the application
type Config struct {
	// Server Configuration
	Server ServerConfig `{{if .HasFeature "config"}}mapstructure:"server"{{else}}json:"server"{{end}}`

	{{if ne .Database.ID ""}}// Database Configuration
	Database DatabaseConfig `{{if .HasFeature "config"}}mapstructure:"database"{{else}}json:"database"{{end}}`
	{{end}}

	// Application Configuration
	App AppConfig `{{if .HasFeature "config"}}mapstructure:"app"{{else}}json:"app"{{end}}`

	{{if .HasFeature "jwt"}}// JWT Configuration
	JWT JWTConfig `{{if .HasFeature "config"}}mapstructure:"jwt"{{else}}json:"jwt"{{end}}`
	{{end}}

	{{if .HasFeature "logging"}}// Logging Configuration
	Logging LoggingConfig `{{if .HasFeature "config"}}mapstructure:"logging"{{else}}json:"logging"{{end}}`
	{{end}}
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Host string `{{if .HasFeature "config"}}mapstructure:"host"{{else}}json:"host"{{end}}`
	Port int    `{{if .HasFeature "config"}}mapstructure:"port"{{else}}json:"port"{{end}}`
}

{{if ne .Database.ID ""}}// DatabaseConfig holds database-related configuration
type DatabaseConfig struct {
	{{if eq .Database.ID "postgres"}}Host     string `{{if .HasFeature "config"}}mapstructure:"host"{{else}}json:"host"{{end}}`
	Port     int    `{{if .HasFeature "config"}}mapstructure:"port"{{else}}json:"port"{{end}}`
	Name     string `{{if .HasFeature "config"}}mapstructure:"name"{{else}}json:"name"{{end}}`
	User     string `{{if .HasFeature "config"}}mapstructure:"user"{{else}}json:"user"{{end}}`
	Password string `{{if .HasFeature "config"}}mapstructure:"password"{{else}}json:"password"{{end}}`
	SSLMode  string `{{if .HasFeature "config"}}mapstructure:"sslmode"{{else}}json:"sslmode"{{end}}`
	URL      string `{{if .HasFeature "config"}}mapstructure:"url"{{else}}json:"url"{{end}}`
	{{else if eq .Database.ID "mysql"}}Host     string `{{if .HasFeature "config"}}mapstructure:"host"{{else}}json:"host"{{end}}`
	Port     int    `{{if .HasFeature "config"}}mapstructure:"port"{{else}}json:"port"{{end}}`
	Name     string `{{if .HasFeature "config"}}mapstructure:"name"{{else}}json:"name"{{end}}`
	User     string `{{if .HasFeature "config"}}mapstructure:"user"{{else}}json:"user"{{end}}`
	Password string `{{if .HasFeature "config"}}mapstructure:"password"{{else}}json:"password"{{end}}`
	URL      string `{{if .HasFeature "config"}}mapstructure:"url"{{else}}json:"url"{{end}}`
	{{else if eq .Database.ID "sqlite"}}Path string `{{if .HasFeature "config"}}mapstructure:"path"{{else}}json:"path"{{end}}`
	{{else if eq .Database.ID "mongodb"}}URI string `{{if .HasFeature "config"}}mapstructure:"uri"{{else}}json:"uri"{{end}}`
	{{else if eq .Database.ID "redis"}}URL      string `{{if .HasFeature "config"}}mapstructure:"url"{{else}}json:"url"{{end}}`
	Password string `{{if .HasFeature "config"}}mapstructure:"password"{{else}}json:"password"{{end}}`
	DB       int    `{{if .HasFeature "config"}}mapstructure:"db"{{else}}json:"db"{{end}}`
	{{else}}URL string `{{if .HasFeature "config"}}mapstructure:"url"{{else}}json:"url"{{end}}`
	{{end}}
}
{{end}}

// AppConfig holds application-level configuration
type AppConfig struct {
	Name        string `{{if .HasFeature "config"}}mapstructure:"name"{{else}}json:"name"{{end}}`
	Version     string `{{if .HasFeature "config"}}mapstructure:"version"{{else}}json:"version"{{end}}`
	Environment string `{{if .HasFeature "config"}}mapstructure:"environment"{{else}}json:"environment"{{end}}`
	Debug       bool   `{{if .HasFeature "config"}}mapstructure:"debug"{{else}}json:"debug"{{end}}`
}

{{if .HasFeature "jwt"}}// JWTConfig holds JWT-related configuration
type JWTConfig struct {
	Secret    string `{{if .HasFeature "config"}}mapstructure:"secret"{{else}}json:"secret"{{end}}`
	ExpiresIn int    `{{if .HasFeature "config"}}mapstructure:"expires_in"{{else}}json:"expires_in"{{end}}` // in hours
}
{{end}}

{{if .HasFeature "logging"}}// LoggingConfig holds logging-related configuration
type LoggingConfig struct {
	Level  string `{{if .HasFeature "config"}}mapstructure:"level"{{else}}json:"level"{{end}}`
	Format string `{{if .HasFeature "config"}}mapstructure:"format"{{else}}json:"format"{{end}}`
}
{{end}}

// Load loads the configuration from environment variables{{if .HasFeature "config"}} and config files{{end}}
func Load() (*Config, error) {
	{{if .HasFeature "config"}}// Initialize Viper
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("$HOME/.{{.ProjectName}}")

	// Enable reading from environment variables
	viper.AutomaticEnv()

	// Set default values
	setDefaults()

	// Read config file (optional)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
		// Config file not found; ignore error and rely on defaults and env vars
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}
	{{else}}config := &Config{
		Server: ServerConfig{
			Host: getEnv("HOST", "localhost"),
			Port: getEnvAsInt("PORT", 8080),
		},
		{{if ne .Database.ID ""}}Database: DatabaseConfig{
			{{if eq .Database.ID "postgres"}}Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 5432),
			Name:     getEnv("DB_NAME", "{{.ProjectName}}_db"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "password"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
			URL:      getEnv("DATABASE_URL", "postgres://postgres:password@localhost:5432/{{.ProjectName}}_db?sslmode=disable"),
			{{else if eq .Database.ID "mysql"}}Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 3306),
			Name:     getEnv("DB_NAME", "{{.ProjectName}}_db"),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", "password"),
			URL:      getEnv("DATABASE_URL", "root:password@tcp(localhost:3306)/{{.ProjectName}}_db?parseTime=true"),
			{{else if eq .Database.ID "sqlite"}}Path: getEnv("DB_PATH", "./{{.ProjectName}}.db"),
			{{else if eq .Database.ID "mongodb"}}URI: getEnv("MONGO_URI", "mongodb://localhost:27017/{{.ProjectName}}_db"),
			{{else if eq .Database.ID "redis"}}URL:      getEnv("REDIS_URL", "redis://localhost:6379/0"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
			{{else}}URL: getEnv("DATABASE_URL", ""),
			{{end}}
		},
		{{end}}App: AppConfig{
			Name:        getEnv("APP_NAME", "{{.ProjectName}}"),
			Version:     getEnv("APP_VERSION", "1.0.0"),
			Environment: getEnv("APP_ENV", "development"),
			Debug:       getEnvAsBool("APP_DEBUG", true),
		},
		{{if .HasFeature "jwt"}}JWT: JWTConfig{
			Secret:    getEnv("JWT_SECRET", "your-super-secret-jwt-key"),
			ExpiresIn: getEnvAsInt("JWT_EXPIRES_IN", 24), // 24 hours
		},
		{{end}}
		{{if .HasFeature "logging"}}Logging: LoggingConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "json"),
		},
		{{end}}
	}
	{{end}}

	{{if .HasFeature "logging"}}// Initialize logger based on config
	initLogger(&config{{if not .HasFeature "config"}}.Logging{{end}})
	{{end}}

	return {{if .HasFeature "config"}}&config{{else}}config{{end}}, nil
}

{{if .HasFeature "config"}}// setDefaults sets default configuration values
func setDefaults() {
	// Server defaults
	viper.SetDefault("server.host", "localhost")
	viper.SetDefault("server.port", 8080)

	{{if ne .Database.ID ""}}// Database defaults
	{{if eq .Database.ID "postgres"}}viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.name", "{{.ProjectName}}_db")
	viper.SetDefault("database.user", "postgres")
	viper.SetDefault("database.password", "password")
	viper.SetDefault("database.sslmode", "disable")
	viper.SetDefault("database.url", "postgres://postgres:password@localhost:5432/{{.ProjectName}}_db?sslmode=disable")
	{{else if eq .Database.ID "mysql"}}viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 3306)
	viper.SetDefault("database.name", "{{.ProjectName}}_db")
	viper.SetDefault("database.user", "root")
	viper.SetDefault("database.password", "password")
	viper.SetDefault("database.url", "root:password@tcp(localhost:3306)/{{.ProjectName}}_db?parseTime=true")
	{{else if eq .Database.ID "sqlite"}}viper.SetDefault("database.path", "./{{.ProjectName}}.db")
	{{else if eq .Database.ID "mongodb"}}viper.SetDefault("database.uri", "mongodb://localhost:27017/{{.ProjectName}}_db")
	{{else if eq .Database.ID "redis"}}viper.SetDefault("database.url", "redis://localhost:6379/0")
	viper.SetDefault("database.password", "")
	viper.SetDefault("database.db", 0)
	{{end}}
	{{end}}

	// App defaults
	viper.SetDefault("app.name", "{{.ProjectName}}")
	viper.SetDefault("app.version", "1.0.0")
	viper.SetDefault("app.environment", "development")
	viper.SetDefault("app.debug", true)

	{{if .HasFeature "jwt"}}// JWT defaults
	viper.SetDefault("jwt.secret", "your-super-secret-jwt-key")
	viper.SetDefault("jwt.expires_in", 24)
	{{end}}

	{{if .HasFeature "logging"}}// Logging defaults
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.format", "json")
	{{end}}
}
{{else}}// getEnv gets an environment variable with a fallback value
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// getEnvAsInt gets an environment variable as integer with a fallback value
func getEnvAsInt(name string, fallback int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return fallback
}

// getEnvAsBool gets an environment variable as boolean with a fallback value
func getEnvAsBool(name string, fallback bool) bool {
	valueStr := getEnv(name, "")
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return fallback
}
{{end}}

{{if .HasFeature "logging"}}// initLogger initializes the global logger based on configuration
func initLogger({{if .HasFeature "config"}}config *Config{{else}}config *LoggingConfig{{end}}) {
	// Set log level
	{{if .HasFeature "config"}}level, err := zerolog.ParseLevel(config.Logging.Level){{else}}level, err := zerolog.ParseLevel(config.Level){{end}}
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	// Set log format
	{{if .HasFeature "config"}}if config.Logging.Format == "console" {
{{else}}if config.Format == "console" {
{{end}}
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	}
}
{{end}}

// Address returns the full server address
func (c *Config) Address() string {
	return c.Server.Host + ":" + strconv.Itoa(c.Server.Port)
}

{{if ne .Database.ID ""}}{{if eq .Database.ID "postgres"}}// PostgresDSN returns the PostgreSQL connection string
func (c *Config) PostgresDSN() string {
	if c.Database.URL != "" {
		return c.Database.URL
	}
	return "host=" + c.Database.Host +
		   " port=" + strconv.Itoa(c.Database.Port) +
		   " user=" + c.Database.User +
		   " password=" + c.Database.Password +
		   " dbname=" + c.Database.Name +
		   " sslmode=" + c.Database.SSLMode
}
{{else if eq .Database.ID "mysql"}}// MySQLDSN returns the MySQL connection string
func (c *Config) MySQLDSN() string {
	if c.Database.URL != "" {
		return c.Database.URL
	}
	return c.Database.User + ":" + c.Database.Password +
		   "@tcp(" + c.Database.Host + ":" + strconv.Itoa(c.Database.Port) + ")/" +
		   c.Database.Name + "?parseTime=true"
}
{{else if eq .Database.ID "sqlite"}}// SQLiteDSN returns the SQLite connection string
func (c *Config) SQLiteDSN() string {
	return c.Database.Path
}
{{else if eq .Database.ID "mongodb"}}// MongoURI returns the MongoDB connection string
func (c *Config) MongoURI() string {
	return c.Database.URI
}
{{else if eq .Database.ID "redis"}}// RedisURL returns the Redis connection string
func (c *Config) RedisURL() string {
	return c.Database.URL
}
{{end}}{{end}}

// IsDevelopment returns true if running in development mode
func (c *Config) IsDevelopment() bool {
	return c.App.Environment == "development"
}

// IsProduction returns true if running in production mode
func (c *Config) IsProduction() bool {
	return c.App.Environment == "production"
}
