// Package config loads all runtime configuration from environment variables.
//
// Load panics when a required variable is missing: a misconfigured deployment
// must fail immediately at boot instead of erroring on the first request.
package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config is the fully resolved configuration of the API process.
type Config struct {
	App      App
	Database Database
	CORS     CORS
	Log      Log
}

// App holds the HTTP server settings.
type App struct {
	Host string
	Port string
	// AutoMigrate runs pending SQL migrations during startup.
	AutoMigrate bool
}

// Addr returns the listen address for the HTTP server.
func (a App) Addr() string {
	return fmt.Sprintf("%s:%s", a.Host, a.Port)
}

// Database holds the PostgreSQL connection settings.
type Database struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

// DSN returns a PostgreSQL connection string understood by pgx.
func (d Database) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		d.User, d.Password, d.Host, d.Port, d.Name,
	)
}

// CORS holds the browser origin allowlist.
type CORS struct {
	Origins []string
}

// Log holds logging settings.
type Log struct {
	Level  string
	Format string
}

// Load reads the configuration from the process environment.
func Load() Config {
	return Config{
		App: App{
			Host:        required("BACKEND_HOST"),
			Port:        required("BACKEND_PORT"),
			AutoMigrate: parseBool("AUTO_MIGRATE"),
		},
		Database: Database{
			Host:     required("DB_HOST"),
			Port:     required("DB_PORT"),
			User:     required("DB_USER"),
			Password: required("DB_PASSWORD"),
			Name:     required("DB_NAME"),
		},
		CORS: CORS{
			Origins: parseList(required("CORS_ORIGINS")),
		},
		Log: Log{
			Level:  required("LOG_LEVEL"),
			Format: required("LOG_FORMAT"),
		},
	}
}

func required(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok || strings.TrimSpace(value) == "" {
		panic(fmt.Sprintf("config: required environment variable %q is not set", key))
	}
	return value
}

func parseBool(key string) bool {
	value, err := strconv.ParseBool(os.Getenv(key))
	if err != nil {
		return false
	}
	return value
}

func parseList(value string) []string {
	parts := strings.Split(value, ",")
	list := make([]string, 0, len(parts))
	for _, part := range parts {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			list = append(list, trimmed)
		}
	}
	return list
}
