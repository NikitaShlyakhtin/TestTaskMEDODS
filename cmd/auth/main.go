package main

import (
	"context"
	"flag"
	"os"
	"sync"
	"time"

	"medods/internal/jsonlog"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type config struct {
	port  int
	env   string
	token struct {
		secret  string
		expires int // Hours until token expires
	}
	db struct {
		connectionString string
	}
}

type application struct {
	config config
	logger *jsonlog.Logger
	wg     sync.WaitGroup
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

	flag.StringVar(&cfg.token.secret, "token-secret", "secret", "JWT secret")
	flag.IntVar(&cfg.token.expires, "token-expires", 72, "JWT expiration in hours")

	flag.StringVar(&cfg.db.connectionString, "db-connection-string", "", "MongoDB connection string")

	flag.Parse()

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	db, err := openDB(cfg)
	if err != nil {
		logger.PrintFatal(err, nil)
	}

	defer db.Disconnect(context.Background())

	logger.PrintInfo("database connection pool established", nil)

	app := &application{
		config: cfg,
		logger: logger,
	}

	err = app.serve()
	if err != nil {
		logger.PrintFatal(err, nil)
	}
}

func openDB(cfg config) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(cfg.db.connectionString)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}
