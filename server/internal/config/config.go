package config

import (
	"os"
	"strings"
)

type AppConfig struct {
	Port      string
	APIKey    string
	TTSAPIKey string
	DBPath    string

	CORSOrigins []string
	CORSMethods []string
	CORSHeaders []string
}

func GetString(envName string, defaultValue string) string {
	if value := os.Getenv(envName); value != "" {
		return value
	}

	return defaultValue
}

func GetStrings(envName string, defaultValue []string) []string {
	if value := GetString(envName, ""); value != "" {
		return strings.Split(value, ",")
	}

	return defaultValue
}
