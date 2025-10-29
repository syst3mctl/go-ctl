package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
{{if .HasFeature "config"}}	"{{.ProjectName}}/internal/config"
{{end}})

// Logger wraps zerolog.Logger with additional functionality
type Logger struct {
	logger zerolog.Logger
}

// Config holds logging configuration
type Config struct {
	Level      string `json:"level" yaml:"level"`
	Format     string `json:"format" yaml:"format"` // json, console, pretty
	Output     string `json:"output" yaml:"output"` // stdout, stderr, file
	Filename   string `json:"filename" yaml:"filename"`
	MaxSize    int    `json:"max_size" yaml:"max_size"`       // MB
	MaxBackups int    `json:"max_backups" yaml:"max_backups"` // number of backups
	MaxAge     int    `json:"max_age" yaml:"max_age"`         // days
	Compress   bool   `json:"compress" yaml:"compress"`
	Caller     bool   `json:"caller" yaml:"caller"`     // include caller info
	Timestamp  bool   `json:"timestamp" yaml:"timestamp"` // include timestamp
}

// DefaultConfig returns default logging configuration
func DefaultConfig() Config {
	return Config{
		Level:      "info",
		Format:     "json",
		Output:     "stdout",
		MaxSize:    100,
		MaxBackups: 3,
		MaxAge:     28,
		Compress:   true,
		Caller:     false,
		Timestamp:  true,
	}
}

{{if .HasFeature "config"}}// NewFromConfig creates a new logger from Viper configuration
func NewFromConfig(cfg *config.Config) (*Logger, error) {
	logConfig := Config{
		Level:      cfg.Logging.Level,
		Format:     cfg.Logging.Format,
		Output:     cfg.Logging.Output,
		Filename:   cfg.Logging.Filename,
		MaxSize:    cfg.Logging.MaxSize,
		MaxBackups: cfg.Logging.MaxBackups,
		MaxAge:     cfg.Logging.MaxAge,
		Compress:   cfg.Logging.Compress,
		Caller:     cfg.App.Debug,
		Timestamp:  true,
	}
	return New(logConfig)
}
{{end}}

// New creates a new logger with the given configuration
func New(cfg Config) (*Logger, error) {
	// Set global settings
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.TimeFieldFormat = time.RFC3339

	// Parse log level
	level, err := zerolog.ParseLevel(strings.ToLower(cfg.Level))
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	// Configure output
	var writers []io.Writer

	switch strings.ToLower(cfg.Output) {
	case "stdout":
		writers = append(writers, os.Stdout)
	case "stderr":
		writers = append(writers, os.Stderr)
	case "file":
		if cfg.Filename == "" {
			cfg.Filename = "{{.ProjectName}}.log"
		}

		// Ensure log directory exists
		if dir := filepath.Dir(cfg.Filename); dir != "." {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return nil, fmt.Errorf("failed to create log directory: %w", err)
			}
		}

		// Use lumberjack for log rotation
		fileWriter := &lumberjackWriter{
			Filename:   cfg.Filename,
			MaxSize:    cfg.MaxSize,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAge,
			Compress:   cfg.Compress,
		}
		writers = append(writers, fileWriter)

	case "both":
		writers = append(writers, os.Stdout)
		if cfg.Filename != "" {
			fileWriter := &lumberjackWriter{
				Filename:   cfg.Filename,
				MaxSize:    cfg.MaxSize,
				MaxBackups: cfg.MaxBackups,
				MaxAge:     cfg.MaxAge,
				Compress:   cfg.Compress,
			}
			writers = append(writers, fileWriter)
		}
	default:
		writers = append(writers, os.Stdout)
	}

	// Create multi-writer
	var output io.Writer
	if len(writers) == 1 {
		output = writers[0]
	} else {
		output = zerolog.MultiLevelWriter(writers...)
	}

	// Configure format
	switch strings.ToLower(cfg.Format) {
	case "console", "text":
		output = zerolog.ConsoleWriter{
			Out:        output,
			TimeFormat: "15:04:05",
			NoColor:    false,
		}
	case "pretty":
		output = zerolog.ConsoleWriter{
			Out:        output,
			TimeFormat: "2006-01-02 15:04:05",
			NoColor:    false,
		}
	case "json":
		// Default JSON format, no changes needed
	default:
		// Default to JSON
	}

	// Create logger
	logger := zerolog.New(output)

	// Configure optional fields
	if cfg.Timestamp {
		logger = logger.With().Timestamp().Logger()
	}

	if cfg.Caller {
		logger = logger.With().Caller().Logger()
	}

	// Add service context
	logger = logger.With().
		Str("service", "{{.ProjectName}}").
		Str("version", getVersion()).
		Logger()

	// Set as global logger
	log.Logger = logger

	return &Logger{logger: logger}, nil
}

// Simple version function - in real apps, this would come from build info
func getVersion() string {
	return "1.0.0"
}

// lumberjackWriter is a simple log rotation writer
// In production, you'd use gopkg.in/natefinch/lumberjack.v2
type lumberjackWriter struct {
	Filename   string
	MaxSize    int  // MB
	MaxBackups int
	MaxAge     int  // days
	Compress   bool

	file *os.File
	size int64
}

func (w *lumberjackWriter) Write(p []byte) (n int, err error) {
	if w.file == nil {
		if err := w.openFile(); err != nil {
			return 0, err
		}
	}

	// Simple size check - rotate if needed
	if w.size+int64(len(p)) > int64(w.MaxSize)*1024*1024 {
		w.rotate()
	}

	n, err = w.file.Write(p)
	w.size += int64(n)
	return n, err
}

func (w *lumberjackWriter) openFile() error {
	file, err := os.OpenFile(w.Filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}

	w.file = file
	w.size = info.Size()
	return nil
}

func (w *lumberjackWriter) rotate() error {
	if w.file != nil {
		w.file.Close()
	}

	// Simple rotation - just rename current file
	timestamp := time.Now().Format("2006-01-02-15-04-05")
	backupName := strings.TrimSuffix(w.Filename, filepath.Ext(w.Filename)) +
		"-" + timestamp + filepath.Ext(w.Filename)

	os.Rename(w.Filename, backupName)

	w.file = nil
	w.size = 0
	return w.openFile()
}

// Logger methods that wrap zerolog

// Debug logs a debug message
func (l *Logger) Debug() *zerolog.Event {
	return l.logger.Debug()
}

// Info logs an info message
func (l *Logger) Info() *zerolog.Event {
	return l.logger.Info()
}

// Warn logs a warning message
func (l *Logger) Warn() *zerolog.Event {
	return l.logger.Warn()
}

// Error logs an error message
func (l *Logger) Error() *zerolog.Event {
	return l.logger.Error()
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal() *zerolog.Event {
	return l.logger.Fatal()
}

// Panic logs a panic message and panics
func (l *Logger) Panic() *zerolog.Event {
	return l.logger.Panic()
}

// With creates a child logger with additional context
func (l *Logger) With() zerolog.Context {
	return l.logger.With()
}

// WithContext returns a copy of ctx with the logger attached
func (l *Logger) WithContext(ctx context.Context) context.Context {
	return l.logger.WithContext(ctx)
}

// GetLogger returns the underlying zerolog.Logger
func (l *Logger) GetLogger() zerolog.Logger {
	return l.logger
}

// Global logger functions for convenience

// Debug logs a debug message using the global logger
func Debug() *zerolog.Event {
	return log.Debug()
}

// Info logs an info message using the global logger
func Info() *zerolog.Event {
	return log.Info()
}

// Warn logs a warning message using the global logger
func Warn() *zerolog.Event {
	return log.Warn()
}

// Error logs an error message using the global logger
func Error() *zerolog.Event {
	return log.Error()
}

// Fatal logs a fatal message and exits using the global logger
func Fatal() *zerolog.Event {
	return log.Fatal()
}

// Panic logs a panic message and panics using the global logger
func Panic() *zerolog.Event {
	return log.Panic()
}

// WithContext returns a copy of ctx with the global logger attached
func WithContext(ctx context.Context) context.Context {
	return log.WithContext(ctx)
}

// FromContext returns the logger from the given context
func FromContext(ctx context.Context) *zerolog.Logger {
	return zerolog.Ctx(ctx)
}

// Structured logging helpers

// HTTPRequest logs HTTP request information
func HTTPRequest(method, path string, statusCode int, duration time.Duration, userAgent string) *zerolog.Event {
	return log.Info().
		Str("method", method).
		Str("path", path).
		Int("status_code", statusCode).
		Dur("duration", duration).
		Str("user_agent", userAgent)
}

// DBQuery logs database query information
func DBQuery(query string, duration time.Duration, rowsAffected int64) *zerolog.Event {
	return log.Debug().
		Str("query", query).
		Dur("duration", duration).
		Int64("rows_affected", rowsAffected)
}

// ServiceCall logs external service call information
func ServiceCall(service, method string, statusCode int, duration time.Duration) *zerolog.Event {
	return log.Info().
		Str("service", service).
		Str("method", method).
		Int("status_code", statusCode).
		Dur("duration", duration)
}

// UserAction logs user action information
func UserAction(userID string, action string, resource string) *zerolog.Event {
	return log.Info().
		Str("user_id", userID).
		Str("action", action).
		Str("resource", resource)
}

// BusinessEvent logs business-related events
func BusinessEvent(event string, details map[string]interface{}) *zerolog.Event {
	evt := log.Info().Str("event_type", "business").Str("event", event)
	for k, v := range details {
		evt = evt.Interface(k, v)
	}
	return evt
}

// SecurityEvent logs security-related events
func SecurityEvent(event string, userID string, ip string, details map[string]interface{}) *zerolog.Event {
	evt := log.Warn().
		Str("event_type", "security").
		Str("event", event).
		Str("user_id", userID).
		Str("ip", ip)

	for k, v := range details {
		evt = evt.Interface(k, v)
	}
	return evt
}

// Performance logs performance metrics
func Performance(operation string, duration time.Duration, details map[string]interface{}) *zerolog.Event {
	evt := log.Info().
		Str("event_type", "performance").
		Str("operation", operation).
		Dur("duration", duration)

	for k, v := range details {
		evt = evt.Interface(k, v)
	}
	return evt
}

// Middleware functions

{{if eq .HTTP.ID "gin"}}// GinMiddleware returns a Gin middleware for request logging
func GinMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		HTTPRequest(
			param.Method,
			param.Path,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
		).
		Str("client_ip", param.ClientIP).
		Int("body_size", param.BodySize).
		Msg("HTTP request completed")
		return ""
	})
}
{{else if eq .HTTP.ID "echo"}}// EchoMiddleware returns an Echo middleware for request logging
func EchoMiddleware() echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogMethod:    true,
		LogURI:       true,
		LogStatus:    true,
		LogLatency:   true,
		LogUserAgent: true,
		LogValuesFunc: func(c echo.Context, values middleware.RequestLoggerValues) error {
			HTTPRequest(
				values.Method,
				values.URI,
				values.Status,
				values.Latency,
				values.UserAgent,
			).
			Str("client_ip", c.RealIP()).
			Msg("HTTP request completed")
			return nil
		},
	})
}
{{else if eq .HTTP.ID "fiber"}}// FiberMiddleware returns a Fiber middleware for request logging
func FiberMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Process request
		err := c.Next()

		// Log request
		HTTPRequest(
			c.Method(),
			c.Path(),
			c.Response().StatusCode(),
			time.Since(start),
			string(c.Request().Header.UserAgent()),
		).
		Str("client_ip", c.IP()).
		Int("body_size", len(c.Response().Body())).
		Msg("HTTP request completed")

		return err
	}
}
{{else}}// HTTPMiddleware returns an HTTP middleware for request logging
func HTTPMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap response writer to capture status code
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Process request
		next.ServeHTTP(wrapped, r)

		// Log request
		HTTPRequest(
			r.Method,
			r.URL.Path,
			wrapped.statusCode,
			time.Since(start),
			r.UserAgent(),
		).
		Str("client_ip", getClientIP(r)).
		Msg("HTTP request completed")
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func getClientIP(r *http.Request) string {
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		return strings.Split(xff, ",")[0]
	}

	xri := r.Header.Get("X-Real-IP")
	if xri != "" {
		return xri
	}

	return r.RemoteAddr
}
{{end}}

// Context helpers

// AddFields adds structured fields to the logger context
func AddFields(ctx context.Context, fields map[string]interface{}) context.Context {
	logger := FromContext(ctx)
	logCtx := logger.With()

	for k, v := range fields {
		logCtx = logCtx.Interface(k, v)
	}

	return logCtx.Logger().WithContext(ctx)
}

// AddField adds a single field to the logger context
func AddField(ctx context.Context, key string, value interface{}) context.Context {
	logger := FromContext(ctx)
	return logger.With().Interface(key, value).Logger().WithContext(ctx)
}

// AddUser adds user information to the logger context
func AddUser(ctx context.Context, userID string, username string) context.Context {
	logger := FromContext(ctx)
	return logger.With().
		Str("user_id", userID).
		Str("username", username).
		Logger().WithContext(ctx)
}

// AddRequest adds request information to the logger context
func AddRequest(ctx context.Context, requestID string, method string, path string) context.Context {
	logger := FromContext(ctx)
	return logger.With().
		Str("request_id", requestID).
		Str("method", method).
		Str("path", path).
		Logger().WithContext(ctx)
}

// Shutdown gracefully shuts down the logger
func Shutdown() error {
	// Close any open log files
	// In a real implementation, you'd properly close file handles
	return nil
}
