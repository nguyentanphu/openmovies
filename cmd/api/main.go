package main

import (
	"context"
	"database/sql"
	"flag"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "openmovies/docs"
	"openmovies/internal/data"
	"openmovies/internal/jsonlog"
	"openmovies/internal/mailer"
	"os"
	"sync"
	"time"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
	smtp struct {
		host     string
		port     int
		username string
		password string
		sender   string
	}
}

type application struct {
	config        config
	logger        *jsonlog.Logger
	validator     *validator.Validate
	schemaDecoder *schema.Decoder
	models        data.Models
	mailer        mailer.Mailer
	wg            sync.WaitGroup
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("OPENMOVIES_DB_DSN"), "POSTGRES DSN")

	flag.StringVar(&cfg.smtp.host, "smtp-host", "sandbox.smtp.mailtrap.io", "Smtp host")
	flag.IntVar(&cfg.smtp.port, "smtp-port", 25, "Smtp port")
	flag.StringVar(&cfg.smtp.username, "smtp-username", "10f7d074fb1524", "Smtp username")
	flag.StringVar(&cfg.smtp.password, "smtp-password", "43703899194531", "Smtp password")
	flag.StringVar(&cfg.smtp.sender, "smtp-sender", "nguyentanphu@flirtingapp.com", "Smtp sender")
	flag.Parse()
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	validate := validator.New(validator.WithRequiredStructEnabled())
	db, err := openDB(cfg)
	defer db.Close()
	if err != nil {
		logger.LogFatal(err, nil)
	}
	app := application{
		logger:        jsonlog.New(os.Stdout, jsonlog.LevelInfo),
		config:        cfg,
		validator:     validate,
		models:        data.NewModels(db),
		schemaDecoder: schema.NewDecoder(),
		mailer:        mailer.New(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username, cfg.smtp.password, cfg.smtp.sender),
		wg:            sync.WaitGroup{},
	}
	err = app.serve()
	if err != nil {
		app.logger.LogFatal(err, nil)
	}
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("pgx", cfg.db.dsn)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
