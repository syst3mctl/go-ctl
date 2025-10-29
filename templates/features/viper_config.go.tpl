package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
{{if .HasFeature "logging"}}	"github.com/rs/zerolog/log"
{{end}})

// Config holds all configuration for the application
type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Server   ServerConfig   `mapstructure:"server"`
{{if ne .DbDriver.ID ""}}	Database DatabaseConfig `mapstructure:"database"`
{{end}}{{if .HasFeature "jwt"}}	JWT      JWTConfig      `mapstructure:"jwt"`
{{end}}{{if .HasFeature "cors"}}	CORS     CORSConfig     `mapstructure:"cors"`
{{end}}{{if .HasFeature "logging"}}	Logging  LoggingConfig  `mapstructure:"logging"`
{{end}}	External ExternalConfig `mapstructure:"external"`
}

// AppConfig holds application-specific configuration
type AppConfig struct {
	Name        string `mapstructure:"name"`
	Version     string `mapstructure:"version"`
	Env         string `mapstructure:"env"`
	Debug       bool   `mapstructure:"debug"`
	Description string `mapstructure:"description"`
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
	TLS          TLSConfig     `mapstructure:"tls"`
}

// TLSConfig holds TLS configuration
type TLSConfig struct {
	Enabled  bool   `mapstructure:"enabled"`
	CertFile string `mapstructure:"cert_file"`
	KeyFile  string `mapstructure:"key_file"`
}

{{if ne .DbDriver.ID ""}}// DatabaseConfig holds database configuration
type DatabaseConfig struct {
{{if eq .Database.ID "postgres"}}	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	User         string        `mapstructure:"user"`
	Password     string        `mapstructure:"password"`
	Name         string        `mapstructure:"name"`
	SSLMode      string        `mapstructure:"ssl_mode"`
	MaxOpenConns int           `mapstructure:"max_open_conns"`
	MaxIdleConns int           `mapstructure:"max_idle_conns"`
	MaxLifetime  time.Duration `mapstructure:"max_lifetime"`
{{else if eq .Database.ID "mysql"}}	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	User         string        `mapstructure:"user"`
	Password     string        `mapstructure:"password"`
	Name         string        `mapstructure:"name"`
	Charset      string        `mapstructure:"charset"`
	ParseTime    bool          `mapstructure:"parse_time"`
	MaxOpenConns int           `mapstructure:"max_open_conns"`
	MaxIdleConns int           `mapstructure:"max_idle_conns"`
	MaxLifetime  time.Duration `mapstructure:"max_lifetime"`
{{else if eq .Database.ID "sqlite"}}	Name        string        `mapstructure:"name"`
	MaxOpenConns int          `mapstructure:"max_open_conns"`
	MaxIdleConns int          `mapstructure:"max_idle_conns"`
	MaxLifetime  time.Duration `mapstructure:"max_lifetime"`
{{else if eq .Database.ID "mongodb"}}	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	AuthDB   string `mapstructure:"auth_db"`
{{else if eq .Database.ID "redis"}}	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
{{else if eq .Database.ID "bigquery"}}	ProjectID     string `mapstructure:"project_id"`
	DatasetID     string `mapstructure:"dataset_id"`
	CredentialsFile string `mapstructure:"credentials_file"`
{{end}}}
{{end}}

{{if .HasFeature "jwt"}}// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret         string        `mapstructure:"secret"`
	Expiration     time.Duration `mapstructure:"expiration"`
	RefreshExpiration time.Duration `mapstructure:"refresh_expiration"`
	Issuer         string        `mapstructure:"issuer"`
	Algorithm      string        `mapstructure:"algorithm"`
}
{{end}}

{{if .HasFeature "cors"}}// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowedOrigins     []string `mapstructure:"allowed_origins"`
	AllowedMethods     []string `mapstructure:"allowed_methods"`
	AllowedHeaders     []string `mapstructure:"allowed_headers"`
	ExposedHeaders     []string `mapstructure:"exposed_headers"`
	AllowCredentials   bool     `mapstructure:"allow_credentials"`
	MaxAge             int      `mapstructure:"max_age"`
	OptionsPassthrough bool     `mapstructure:"options_passthrough"`
}
{{end}}

{{if .HasFeature "logging"}}// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level      string `mapstructure:"level"`
	Format     string `mapstructure:"format"` // json, console
	Output     string `mapstructure:"output"` // stdout, stderr, file
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`    // MB
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`     // days
	Compress   bool   `mapstructure:"compress"`
}
{{end}}

// ExternalConfig holds external service configurations
type ExternalConfig struct {
	APITimeout time.Duration            `mapstructure:"api_timeout"`
	Services   map[string]ServiceConfig `mapstructure:"services"`
}

// ServiceConfig holds configuration for external services
type ServiceConfig struct {
	BaseURL string            `mapstructure:"base_url"`
	Timeout time.Duration     `mapstructure:"timeout"`
	Headers map[string]string `mapstructure:"headers"`
	Auth    AuthConfig        `mapstructure:"auth"`
}

// AuthConfig holds authentication configuration for external services
type AuthConfig struct {
	Type   string `mapstructure:"type"` // bearer, basic, apikey
	Token  string `mapstructure:"token"`
	User   string `mapstructure:"user"`
	Pass   string `mapstructure:"pass"`
	Header string `mapstructure:"header"`
}

// Load loads configuration from multiple sources
func Load() (*Config, error) {
	return LoadWithOptions(LoadOptions{})
}

// LoadOptions holds options for loading configuration
type LoadOptions struct {
	ConfigFile   string
	ConfigPaths  []string
	ConfigType   string
	EnvPrefix    string
	AutomaticEnv bool
}

// LoadWithOptions loads configuration with custom options
func LoadWithOptions(opts LoadOptions) (*Config, error) {
	v := viper.New()

	// Set defaults
	setDefaults(v)

	// Configure viper
	if opts.ConfigFile != "" {
		v.SetConfigFile(opts.ConfigFile)
	} else {
		v.SetConfigName("config")
		if opts.ConfigType != "" {
			v.SetConfigType(opts.ConfigType)
		} else {
			v.SetConfigType("yaml")
		}
	}

	// Set config paths
	configPaths := opts.ConfigPaths
	if len(configPaths) == 0 {
		configPaths = []string{
			".",
			"./config",
			"./configs",
			"/etc/{{.ProjectName}}",
			"$HOME/.{{.ProjectName}}",
		}
	}
	for _, path := range configPaths {
		v.AddConfigPath(path)
	}

	// Environment variables
	envPrefix := opts.EnvPrefix
	if envPrefix == "" {
		envPrefix = strings.ToUpper(strings.ReplaceAll("{{.ProjectName}}", "-", "_"))
	}
	v.SetEnvPrefix(envPrefix)
	if opts.AutomaticEnv || opts.AutomaticEnv == false {
		v.AutomaticEnv()
	}
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	// Bind specific environment variables
	bindEnvVariables(v)

	// Read config file
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
{{if .HasFeature "logging"}}		log.Warn().Msg("Config file not found, using defaults and environment variables")
{{end}}	} else {
{{if .HasFeature "logging"}}		log.Info().Str("config_file", v.ConfigFileUsed()).Msg("Using config file")
{{end}}	}

	// Unmarshal config
	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

{{if .HasFeature "logging"}}	log.Info().
		Str("app_name", config.App.Name).
		Str("app_version", config.App.Version).
		Str("env", config.App.Env).
		Str("host", config.Server.Host).
		Int("port", config.Server.Port).
		Msg("Configuration loaded successfully")
{{end}}

	return &config, nil
}

// setDefaults sets default values for configuration
func setDefaults(v *viper.Viper) {
	// App defaults
	v.SetDefault("app.name", "{{.ProjectName}}")
	v.SetDefault("app.version", "1.0.0")
	v.SetDefault("app.env", "development")
	v.SetDefault("app.debug", false)
	v.SetDefault("app.description", "{{.ProjectName}} application")

	// Server defaults
	v.SetDefault("server.host", "localhost")
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.read_timeout", "15s")
	v.SetDefault("server.write_timeout", "15s")
	v.SetDefault("server.idle_timeout", "60s")
	v.SetDefault("server.tls.enabled", false)

{{if ne .DbDriver.ID ""}}	// Database defaults
{{if eq .Database.ID "postgres"}}	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 5432)
	v.SetDefault("database.user", "postgres")
	v.SetDefault("database.password", "")
	v.SetDefault("database.name", "{{.ProjectName}}")
	v.SetDefault("database.ssl_mode", "disable")
	v.SetDefault("database.max_open_conns", 25)
	v.SetDefault("database.max_idle_conns", 25)
	v.SetDefault("database.max_lifetime", "5m")
{{else if eq .Database.ID "mysql"}}	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 3306)
	v.SetDefault("database.user", "root")
	v.SetDefault("database.password", "")
	v.SetDefault("database.name", "{{.ProjectName}}")
	v.SetDefault("database.charset", "utf8mb4")
	v.SetDefault("database.parse_time", true)
	v.SetDefault("database.max_open_conns", 25)
	v.SetDefault("database.max_idle_conns", 25)
	v.SetDefault("database.max_lifetime", "5m")
{{else if eq .Database.ID "sqlite"}}	v.SetDefault("database.name", "./{{.ProjectName}}.db")
	v.SetDefault("database.max_open_conns", 1)
	v.SetDefault("database.max_idle_conns", 1)
	v.SetDefault("database.max_lifetime", "1h")
{{else if eq .Database.ID "mongodb"}}	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 27017)
	v.SetDefault("database.name", "{{.ProjectName}}")
	v.SetDefault("database.auth_db", "admin")
{{else if eq .Database.ID "redis"}}	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 6379)
	v.SetDefault("database.db", 0)
{{else if eq .Database.ID "bigquery"}}	v.SetDefault("database.project_id", "")
	v.SetDefault("database.dataset_id", "{{.ProjectName}}")
{{end}}
{{end}}

{{if .HasFeature "jwt"}}	// JWT defaults
	v.SetDefault("jwt.secret", "change-me-in-production")
	v.SetDefault("jwt.expiration", "24h")
	v.SetDefault("jwt.refresh_expiration", "168h") // 7 days
	v.SetDefault("jwt.issuer", "{{.ProjectName}}")
	v.SetDefault("jwt.algorithm", "HS256")
{{end}}

{{if .HasFeature "cors"}}	// CORS defaults
	v.SetDefault("cors.allowed_origins", []string{"http://localhost:3000"})
	v.SetDefault("cors.allowed_methods", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	v.SetDefault("cors.allowed_headers", []string{"*"})
	v.SetDefault("cors.allow_credentials", true)
	v.SetDefault("cors.max_age", 300)
{{end}}

{{if .HasFeature "logging"}}	// Logging defaults
	v.SetDefault("logging.level", "info")
	v.SetDefault("logging.format", "json")
	v.SetDefault("logging.output", "stdout")
	v.SetDefault("logging.max_size", 100)
	v.SetDefault("logging.max_backups", 3)
	v.SetDefault("logging.max_age", 28)
	v.SetDefault("logging.compress", true)
{{end}}

	// External services defaults
	v.SetDefault("external.api_timeout", "30s")
}

// bindEnvVariables binds specific environment variables
func bindEnvVariables(v *viper.Viper) {
	// Common environment variables
	v.BindEnv("app.env", "APP_ENV", "ENVIRONMENT")
	v.BindEnv("app.debug", "APP_DEBUG", "DEBUG")
	v.BindEnv("server.port", "PORT", "HTTP_PORT")
	v.BindEnv("server.host", "HOST", "HTTP_HOST")

{{if ne .DbDriver.ID ""}}	// Database environment variables
{{if eq .Database.ID "postgres"}}	v.BindEnv("database.host", "POSTGRES_HOST", "DB_HOST")
	v.BindEnv("database.port", "POSTGRES_PORT", "DB_PORT")
	v.BindEnv("database.user", "POSTGRES_USER", "DB_USER")
	v.BindEnv("database.password", "POSTGRES_PASSWORD", "DB_PASSWORD")
	v.BindEnv("database.name", "POSTGRES_DB", "DB_NAME")
	v.BindEnv("database.ssl_mode", "POSTGRES_SSL_MODE", "DB_SSL_MODE")
{{else if eq .Database.ID "mysql"}}	v.BindEnv("database.host", "MYSQL_HOST", "DB_HOST")
	v.BindEnv("database.port", "MYSQL_PORT", "DB_PORT")
	v.BindEnv("database.user", "MYSQL_USER", "DB_USER")
	v.BindEnv("database.password", "MYSQL_PASSWORD", "DB_PASSWORD")
	v.BindEnv("database.name", "MYSQL_DATABASE", "DB_NAME")
{{else if eq .Database.ID "sqlite"}}	v.BindEnv("database.name", "SQLITE_DB", "DB_NAME")
{{else if eq .Database.ID "mongodb"}}	v.BindEnv("database.host", "MONGO_HOST", "DB_HOST")
	v.BindEnv("database.port", "MONGO_PORT", "DB_PORT")
	v.BindEnv("database.user", "MONGO_USER", "DB_USER")
	v.BindEnv("database.password", "MONGO_PASSWORD", "DB_PASSWORD")
	v.BindEnv("database.name", "MONGO_DATABASE", "DB_NAME")
{{else if eq .Database.ID "redis"}}	v.BindEnv("database.host", "REDIS_HOST", "DB_HOST")
	v.BindEnv("database.port", "REDIS_PORT", "DB_PORT")
	v.BindEnv("database.password", "REDIS_PASSWORD", "DB_PASSWORD")
	v.BindEnv("database.db", "REDIS_DB", "DB_NAME")
{{else if eq .Database.ID "bigquery"}}	v.BindEnv("database.project_id", "GCP_PROJECT_ID", "GOOGLE_CLOUD_PROJECT")
	v.BindEnv("database.dataset_id", "BIGQUERY_DATASET")
	v.BindEnv("database.credentials_file", "GOOGLE_APPLICATION_CREDENTIALS")
{{end}}
{{end}}

{{if .HasFeature "jwt"}}	// JWT environment variables
	v.BindEnv("jwt.secret", "JWT_SECRET")
	v.BindEnv("jwt.expiration", "JWT_EXPIRATION")
{{end}}

{{if .HasFeature "logging"}}	// Logging environment variables
	v.BindEnv("logging.level", "LOG_LEVEL")
	v.BindEnv("logging.format", "LOG_FORMAT")
	v.BindEnv("logging.output", "LOG_OUTPUT")
{{end}}
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.App.Name == "" {
		return fmt.Errorf("app.name is required")
	}

	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return fmt.Errorf("server.port must be between 1 and 65535")
	}

{{if ne .DbDriver.ID ""}}	// Database validation
	if err := c.validateDatabase(); err != nil {
		return fmt.Errorf("database config validation failed: %w", err)
	}
{{end}}

{{if .HasFeature "jwt"}}	// JWT validation
	if c.JWT.Secret == "" || c.JWT.Secret == "change-me-in-production" {
		if c.IsProduction() {
			return fmt.Errorf("jwt.secret must be set in production")
		}
{{if .HasFeature "logging"}}		log.Warn().Msg("Using default JWT secret in non-production environment")
{{end}}	}
{{end}}

	return nil
}

{{if ne .DbDriver.ID ""}}// validateDatabase validates database-specific configuration
func (c *Config) validateDatabase() error {
{{if eq .Database.ID "postgres"}}	if c.Database.Host == "" {
		return fmt.Errorf("database.host is required")
	}
	if c.Database.Port <= 0 {
		return fmt.Errorf("database.port is required")
	}
	if c.Database.Name == "" {
		return fmt.Errorf("database.name is required")
	}
{{else if eq .Database.ID "mysql"}}	if c.Database.Host == "" {
		return fmt.Errorf("database.host is required")
	}
	if c.Database.Port <= 0 {
		return fmt.Errorf("database.port is required")
	}
	if c.Database.Name == "" {
		return fmt.Errorf("database.name is required")
	}
{{else if eq .Database.ID "sqlite"}}	if c.Database.Name == "" {
		return fmt.Errorf("database.name is required")
	}
{{else if eq .Database.ID "mongodb"}}	if c.Database.Host == "" {
		return fmt.Errorf("database.host is required")
	}
	if c.Database.Name == "" {
		return fmt.Errorf("database.name is required")
	}
{{else if eq .Database.ID "redis"}}	if c.Database.Host == "" {
		return fmt.Errorf("database.host is required")
	}
{{else if eq .Database.ID "bigquery"}}	if c.Database.ProjectID == "" {
		return fmt.Errorf("database.project_id is required")
	}
{{end}}	return nil
}
{{end}}

// IsProduction returns true if the environment is production
func (c *Config) IsProduction() bool {
	return strings.ToLower(c.App.Env) == "production" || strings.ToLower(c.App.Env) == "prod"
}

// IsDevelopment returns true if the environment is development
func (c *Config) IsDevelopment() bool {
	return strings.ToLower(c.App.Env) == "development" || strings.ToLower(c.App.Env) == "dev"
}

// IsTest returns true if the environment is test
func (c *Config) IsTest() bool {
	return strings.ToLower(c.App.Env) == "test" || strings.ToLower(c.App.Env) == "testing"
}

// Address returns the server address
func (c *Config) Address() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}

{{if ne .DbDriver.ID ""}}// DatabaseDSN returns the database connection string
func (c *Config) DatabaseDSN() string {
{{if eq .Database.ID "postgres"}}	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host, c.Database.Port, c.Database.User, c.Database.Password, c.Database.Name, c.Database.SSLMode)
{{else if eq .Database.ID "mysql"}}	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		c.Database.User, c.Database.Password, c.Database.Host, c.Database.Port, c.Database.Name)
	if c.Database.Charset != "" {
		dsn += "?charset=" + c.Database.Charset
	}
	if c.Database.ParseTime {
		if strings.Contains(dsn, "?") {
			dsn += "&parseTime=true"
		} else {
			dsn += "?parseTime=true"
		}
	}
	return dsn
{{else if eq .Database.ID "sqlite"}}	return c.Database.Name
{{else if eq .Database.ID "mongodb"}}	if c.Database.User != "" && c.Database.Password != "" {
		return fmt.Sprintf("mongodb://%s:%s@%s:%d/%s?authSource=%s",
			c.Database.User, c.Database.Password, c.Database.Host, c.Database.Port, c.Database.Name, c.Database.AuthDB)
	}
	return fmt.Sprintf("mongodb://%s:%d/%s", c.Database.Host, c.Database.Port, c.Database.Name)
{{else if eq .Database.ID "redis"}}	return fmt.Sprintf("%s:%d", c.Database.Host, c.Database.Port)
{{else if eq .Database.ID "bigquery"}}	return fmt.Sprintf("%s.%s", c.Database.ProjectID, c.Database.DatasetID)
{{else}}	return ""
{{end}}
}
{{end}}

// GetString returns a string configuration value with fallback
func (c *Config) GetString(key, fallback string) string {
	v := viper.GetString(key)
	if v == "" {
		return fallback
	}
	return v
}

// GetInt returns an integer configuration value with fallback
func (c *Config) GetInt(key string, fallback int) int {
	if !viper.IsSet(key) {
		return fallback
	}
	return viper.GetInt(key)
}

// GetBool returns a boolean configuration value with fallback
func (c *Config) GetBool(key string, fallback bool) bool {
	if !viper.IsSet(key) {
		return fallback
	}
	return viper.GetBool(key)
}

// GetDuration returns a duration configuration value with fallback
func (c *Config) GetDuration(key string, fallback time.Duration) time.Duration {
	if !viper.IsSet(key) {
		return fallback
	}
	return viper.GetDuration(key)
}

// GetStringSlice returns a string slice configuration value with fallback
func (c *Config) GetStringSlice(key string, fallback []string) []string {
	if !viper.IsSet(key) {
		return fallback
	}
	return viper.GetStringSlice(key)
}

// Print prints the configuration (without sensitive data)
func (c *Config) Print() {
	fmt.Printf("Configuration:\n")
	fmt.Printf("  App:\n")
	fmt.Printf("    Name: %s\n", c.App.Name)
	fmt.Printf("    Version: %s\n", c.App.Version)
	fmt.Printf("    Environment: %s\n", c.App.Env)
	fmt.Printf("    Debug: %t\n", c.App.Debug)
	fmt.Printf("  Server:\n")
	fmt.Printf("    Address: %s\n", c.Address())
	fmt.Printf("    Read Timeout: %s\n", c.Server.ReadTimeout)
	fmt.Printf("    Write Timeout: %s\n", c.Server.WriteTimeout)
	fmt.Printf("    TLS Enabled: %t\n", c.Server.TLS.Enabled)
{{if ne .DbDriver.ID ""}}	fmt.Printf("  Database:\n")
	fmt.Printf("    Type: {{.Database.ID}}\n")
{{if or (eq .Database.ID "postgres") (eq .Database.ID "mysql") (eq .Database.ID "mongodb")}}	fmt.Printf("    Host: %s\n", c.Database.Host)
	fmt.Printf("    Port: %d\n", c.Database.Port)
{{end}}	fmt.Printf("    Name: %s\n", c.Database.Name)
{{end}}{{if .HasFeature "jwt"}}	fmt.Printf("  JWT:\n")
	fmt.Printf("    Expiration: %s\n", c.JWT.Expiration)
	fmt.Printf("    Algorithm: %s\n", c.JWT.Algorithm)
{{end}}{{if .HasFeature "logging"}}	fmt.Printf("  Logging:\n")
	fmt.Printf("    Level: %s\n", c.Logging.Level)
	fmt.Printf("    Format: %s\n", c.Logging.Format)
	fmt.Printf("    Output: %s\n", c.Logging.Output)
{{end}}	fmt.Printf("  External:\n")
	fmt.Printf("    API Timeout: %s\n", c.External.APITimeout)
}
