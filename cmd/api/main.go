package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type serverConfig struct {
	port        int
	environment string
}

type applicationDependencies struct {
	Username string
	Password string
	Realm    string
	Dir      string
	config   serverConfig
	logger   *slog.Logger
}

func main() {
	var setting serverConfig

	flag.IntVar(&setting.port, "port", 4000, "Server port")
	flag.StringVar(&setting.environment, "env", "development", "Environment (development|staging|production)")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	appInstance := &applicationDependencies{
		Username: "admin",
		Password: "secret",
		Realm:    "Authorized",
		config:   setting,
		logger:   logger,
	}

	apiServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", setting.port),
		Handler:      appInstance.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	logger.Info("Starting server", "address", apiServer.Addr, "environment", setting.environment)
	err := apiServer.ListenAndServe()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
