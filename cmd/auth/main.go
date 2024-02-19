package main

import (
	"flag"
	"os"
	"sync"

	"medods/internal/jsonlog"
)

type serverConfig struct {
	port int
	env  string
}

type tokenConfig struct {
	secret  string
	expires int // Hours until token expires
}

type application struct {
	serverConfig serverConfig
	logger       *jsonlog.Logger
	wg           sync.WaitGroup
	tokenConfig  tokenConfig
}

func main() {
	var serverConfig serverConfig
	var tokenConfig tokenConfig

	flag.IntVar(&serverConfig.port, "port", 4000, "API server port")
	flag.StringVar(&serverConfig.env, "env", "development", "Environment (development|staging|production)")

	flag.StringVar(&tokenConfig.secret, "token-secret", "secret", "JWT secret")
	flag.IntVar(&tokenConfig.expires, "token-expires", 72, "JWT expiration in hours")

	flag.Parse()

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	app := &application{
		serverConfig: serverConfig,
		logger:       logger,
		tokenConfig:  tokenConfig,
	}

	err := app.serve()
	if err != nil {
		logger.PrintFatal(err, nil)
	}
}
