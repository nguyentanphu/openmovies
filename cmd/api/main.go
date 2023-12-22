package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"github.com/go-playground/validator/v10"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
	"net/http"
	_ "openmovies/docs"
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
	config    config
	logger    *log.Logger
	validator *validator.Validate
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("OPENMOVIES_DB_DSN"), "POSTGRES DSN")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	validate := validator.New(validator.WithRequiredStructEnabled())
	db, err := openDB(cfg)
	defer db.Close()
	if err != nil {
		logger.Fatal(err)
	}
	app := application{
		logger:    logger,
		config:    cfg,
		validator: validate,
	}
	srv := http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	app.logger.Printf("Starting server %s with port %d", app.config.env, app.config.port)
	err = srv.ListenAndServe()
	if err != nil {
		app.logger.Fatal(err)
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
