package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/madeindra/mock-interview/server/internal/config"
	"github.com/madeindra/mock-interview/server/internal/handler"
)

const (
	envPort   = "PORT"
	envAPIKey = "OPENAI_API_KEY"
	envDBURI  = "DB_URI"

	envCORSOrigins = "CORS_ALLOWED_ORIGINS"
	envCORSMethods = "CORS_ALLOWED_METHODS"
	envCORSHeaders = "CORS_ALLOWED_HEADERS"

	defaultPort = "8080"
)

var (
	defaultCORSOrigin  = []string{"*"}
	defaultCORSMethods = []string{"GET", "POST"}
	defaultCORSHeaders = []string{"Accept", "Authorization", "Content-Type"}
)

func main() {
	cfg, err := initConfig()
	if err != nil {
		log.Fatal(err)
	}

	router := handler.NewHandler(cfg)

	server := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%s", cfg.Port),
		Handler: router,
	}

	log.Printf("Server listening on port %s", cfg.Port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func initConfig() (config.AppConfig, error) {
	cfg := config.AppConfig{
		Port:        config.GetString(envPort, defaultPort),
		APIKey:      config.GetString(envAPIKey, ""),
		DBURI:       config.GetString(envDBURI, ""),
		CORSOrigins: config.GetStrings(envCORSOrigins, defaultCORSOrigin),
		CORSMethods: config.GetStrings(envCORSMethods, defaultCORSMethods),
		CORSHeaders: config.GetStrings(envCORSHeaders, defaultCORSHeaders),
	}

	if cfg.APIKey == "" && cfg.DBURI == "" {
		return config.AppConfig{}, fmt.Errorf("API Key and DB URI is needed")
	}

	return cfg, nil
}
