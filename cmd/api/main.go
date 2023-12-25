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
	"os"
	"time"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
}

type application struct {
	config        config
	logger        *jsonlog.Logger
	validator     *validator.Validate
	schemaDecoder *schema.Decoder
	models        data.Models
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("OPENMOVIES_DB_DSN"), "POSTGRES DSN")
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
