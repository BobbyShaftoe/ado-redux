package config

import (
	"fmt"
	"log/slog"
	"os"
)

type EnvVars struct {
	DBPass       string
	PAT          string
	ORGANIZATION string
	PROJECT      string
	USER         string
}

type localLogger struct {
	json *slog.Logger
}

var logger = &localLogger{
	json: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
}

func fatal(v ...any) {
	logger.json.Error("main", "err", fmt.Sprint(v...))
	os.Exit(1)
}

func New() *EnvVars {
	envVars := &EnvVars{}
	envVars.getEnv()
	return envVars
}

func (EnvVars *EnvVars) getEnv() {
	EnvVars.DBPass = os.Getenv("DB_PASS")
	EnvVars.PAT = os.Getenv("AZURE_TOKEN")
	EnvVars.ORGANIZATION = os.Getenv("ADO_ORG")
	EnvVars.PROJECT = os.Getenv("ADO_DEFAULT_PROJECT")
	EnvVars.USER = os.Getenv("ADO_DEFAULT_USER")

	if EnvVars.DBPass == "" {
		fatal("DB_PASS environment variable not set")
	}
	if EnvVars.PAT == "" {
		fatal("AZURE_TOKEN environment variable not set")
	}
	if EnvVars.ORGANIZATION == "" {
		fatal("ADO_ORG environment variable not set")
	}
	if EnvVars.PROJECT == "" {
		fatal("ADO_DEFAULT_PROJECT environment variable not set")
	}
	if EnvVars.USER == "" {
		fatal("ADO_DEFAULT_USER environment variable not set")
	}
}
