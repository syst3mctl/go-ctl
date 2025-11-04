package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"{{.ProjectName}}/internal/store"
{{if or (eq .Database.ID "redis") (hasRedis .Databases)}}	"github.com/redis/go-redis/v9"
{{end}})

var (
	maxOpenCons = 30
	maxIdleCons = 30
)

type Application struct {
	Debug    bool
	ErrorLog *log.Logger
	InfoLog  *log.Logger
	Config   Config
	Store    store.Store
{{if or (eq .Database.ID "redis") (hasRedis .Databases)}}	Redis    *redis.Client
{{end}}}

type Config struct {
	Addr        string
	ApiURL      string
	DBConfig    DBConfig
	Env         string
{{if .HasFeature "nats"}}	NatsCluster string
{{end}}{{if or (eq .Database.ID "redis") (hasRedis .Databases)}}	RedisConfig RedisConfig
{{end}}}

type DBConfig struct {
	Addr         string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  time.Duration
}

{{if or (eq .Database.ID "redis") (hasRedis .Databases)}}type RedisConfig struct {
	Addr     string
	Password string
	DB       string
}
{{end}}

func NewConfig() Config {
	idleTime := time.Duration(GetInt("DB_MAX_IDLE_TIME", 30))

	return Config{
{{if .HasFeature "nats"}}		NatsCluster: GetString("NATS_CLUSTER", "frx.session"),
{{end}}		Addr:        GetString("APP_ADDR", ":8080"),
		DBConfig: DBConfig{
			Addr:         GetString("DB_ADDR", getDefaultDBAddr()),
			MaxOpenConns: maxOpenCons,
			MaxIdleConns: maxIdleCons,
			MaxIdleTime:  idleTime * time.Second,
		},
		Env:    GetString("APP_ENV", "development"),
		ApiURL: GetString("APP_API_URL", "http://localhost:8080"),
{{if or (eq .Database.ID "redis") (hasRedis .Databases)}}		RedisConfig: RedisConfig{
			Addr:     GetString("REDIS_ADDR", "localhost:6379"),
			Password: GetString("REDIS_PASSWORD", ""),
			DB:       strconv.Itoa(GetInt("REDIS_DB", 0)),
		},
{{end}}	}
}

func getDefaultDBAddr() string {
{{if eq .Database.ID "postgres"}}	return "postgres://user:psw@localhost/db?sslmode=disable"
{{else if eq .Database.ID "mysql"}}	return "user:psw@tcp(localhost:3306)/db"
{{else if eq .Database.ID "sqlite"}}	return "file:./db.sqlite?cache=shared&mode=rwc"
{{else}}	return "postgres://user:psw@localhost/db?sslmode=disable"
{{end}}}

func GetString(key, fallback string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	return val
}

func GetInt(key string, fallback int) int {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	i, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}
	return i
}

