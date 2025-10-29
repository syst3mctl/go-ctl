package config

import (
	"os"
	"strconv"

	"github.com/spf13/viper"


)

// Config holds all configuration for the application
type Config struct {
	// Server Configuration
	Server ServerConfig `mapstructure:"server"`

	// Database Configuration
	Database DatabaseConfig `mapstructure:"database"`
	

	// Application Configuration
	App AppConfig `mapstructure:"app"`

	

	
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

// DatabaseConfig holds database-related configuration
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Name     string `mapstructure:"name"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	SSLMode  string `mapstructure:"sslmode"`
	URL      string `mapstructure:"url"`
	
}


// AppConfig holds application-level configuration
type AppConfig struct {
	Name        string `mapstructure:"name"`
	Version     string `mapstructure:"version"`
	Environment string `mapstructure:"environment"`
	Debug       bool   `mapstructure:"debug"`
}





// Load loads the configuration from environment variables and config files
func Load() (*Config, error) {
	// Initialize Viper
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("$HOME/.test-app")

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
	

	

	return &config, nil
}

// setDefaults sets default configuration values
func setDefaults() {
	// Server defaults
	viper.SetDefault("server.host", "localhost")
	viper.SetDefault("server.port", 8080)

	// Database defaults
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.name", "test-app_db")
	viper.SetDefault("database.user", "postgres")
	viper.SetDefault("database.password", "password")
	viper.SetDefault("database.sslmode", "disable")
	viper.SetDefault("database.url", "postgres://postgres:password@localhost:5432/test-app_db?sslmode=disable")
	
	

	// App defaults
	viper.SetDefault("app.name", "test-app")
	viper.SetDefault("app.version", "1.0.0")
	viper.SetDefault("app.environment", "development")
	viper.SetDefault("app.debug", true)

	

	
}




// Address returns the full server address
func (c *Config) Address() string {
	return c.Server.Host + ":" + strconv.Itoa(c.Server.Port)
}

// PostgresDSN returns the PostgreSQL connection string
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


// IsDevelopment returns true if running in development mode
func (c *Config) IsDevelopment() bool {
	return c.App.Environment == "development"
}

// IsProduction returns true if running in production mode
func (c *Config) IsProduction() bool {
	return c.App.Environment == "production"
}
