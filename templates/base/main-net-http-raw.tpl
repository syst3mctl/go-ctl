package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"{{.ProjectName}}/cmd/config"
	"{{.ProjectName}}/internal/db"
	"{{.ProjectName}}/internal/handlers"
	"{{.ProjectName}}/internal/store"
	"{{.ProjectName}}/internal/validate"
)

var app *config.Application

var (
	readTimeout  = 5 * time.Second
	writeTimeout = 10 * time.Second
)

func main() {
	run()
}

func run() {
	flag.Parse()

	infoLog := log.New(os.Stderr, "INFO: ", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Llongfile)
	debug := flag.Bool("debug", false, "Enable debug mode")

	cfConn := config.NewConfig()

{{if eq .DbDriver.ID "gorm"}}	dbConn, err := db.NewConn(cfConn.DBConfig.Addr, cfConn.DBConfig.MaxOpenConns, cfConn.DBConfig.MaxIdleConns, cfConn.DBConfig.MaxIdleTime)
	if err != nil {
		errorLog.Fatalln(err)
	}
	sqlDB, err := dbConn.DB()
	if err != nil {
		errorLog.Fatalln(err)
	}
	defer sqlDB.Close()

	st := store.NewStorage(dbConn)
{{else if eq .DbDriver.ID "sqlx"}}	dbConn, err := db.NewConn(cfConn.DBConfig.Addr, cfConn.DBConfig.MaxOpenConns, cfConn.DBConfig.MaxIdleConns, cfConn.DBConfig.MaxIdleTime)
	if err != nil {
		errorLog.Fatalln(err)
	}
	defer dbConn.Close()

	st := store.NewStorage(dbConn)
{{else}}	dbConn, err := db.NewConn(cfConn.DBConfig.Addr, cfConn.DBConfig.MaxOpenConns, cfConn.DBConfig.MaxIdleConns, cfConn.DBConfig.MaxIdleTime)
	if err != nil {
		errorLog.Fatalln(err)
	}
	defer dbConn.Close()

	st := store.NewStorage(dbConn)
{{end}}

{{if or (eq .Database.ID "redis") (hasRedis .Databases)}}	redisConn, err := db.NewRedisClient(cfConn.RedisConfig.Addr, cfConn.RedisConfig.Password, cfConn.RedisConfig.DB)
	if err != nil {
		errorLog.Fatalln(err)
	}
	defer redisConn.Close()

{{end}}	app = &config.Application{
		Debug:    *debug,
		InfoLog:  infoLog,
		ErrorLog: errorLog,
		Store:    st,
		Config:   cfConn,
{{if or (eq .Database.ID "redis") (hasRedis .Databases)}}		Redis:    redisConn,
{{end}}	}

	handlers.NewHandler(app)
	validate.NewValidate(app)

	srv := &http.Server{
		Addr:        cfConn.Addr,
		Handler:     handlers.CorsMiddleware(handlers.Routes()),
		ErrorLog:    errorLog,
		IdleTimeout: readTimeout,
		ReadTimeout: writeTimeout,
	}

	infoLog.Printf("Server starting on %s", cfConn.Addr)

	err = srv.ListenAndServe()
	if err != nil {
		errorLog.Fatalln(err)
	}
}

